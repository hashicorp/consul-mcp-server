// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// generateTestCertAndKey creates a self-signed certificate and private key for testing.
func generateTestCertAndKey() (certPEM []byte, keyPEM []byte, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		IsCA:         true,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal key: %w", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return certPEM, keyPEM, nil
}

// writeTempFile creates a temporary file with the given content and returns its path.
func writeTempFile(t *testing.T, pattern string, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp("", pattern)
	require.NoError(t, err)
	_, err = f.Write(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestTLSConfig(t *testing.T) {
	// Test with no TLS config
	os.Unsetenv("MCP_TLS_CERT_FILE")
	os.Unsetenv("MCP_TLS_KEY_FILE")

	tlsConfig, err := GetTLSConfigFromEnv()
	require.NoError(t, err)
	require.Nil(t, tlsConfig)
}

func TestHTTPServerWithTLS(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	srv := httptest.NewTLSServer(mux)
	defer srv.Close()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Get(srv.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTLSConfigWithValidCert(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertAndKey()
	require.NoError(t, err)

	certPath := writeTempFile(t, "test_cert_*.pem", certPEM)
	keyPath := writeTempFile(t, "test_key_*.pem", keyPEM)

	os.Setenv("MCP_TLS_CERT_FILE", certPath)
	os.Setenv("MCP_TLS_KEY_FILE", keyPath)
	defer func() {
		os.Unsetenv("MCP_TLS_CERT_FILE")
		os.Unsetenv("MCP_TLS_KEY_FILE")
	}()

	tlsConfig, err := GetTLSConfigFromEnv()
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)
	require.Equal(t, certPath, tlsConfig.CertFile)
	require.Equal(t, keyPath, tlsConfig.KeyFile)
	require.Equal(t, uint16(tls.VersionTLS12), tlsConfig.Config.MinVersion)
}

func TestTLSConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		certFile  string
		keyFile   string
		wantNil   bool
		wantError bool
	}{
		{"both empty", "", "", true, false},
		{"cert only", "cert.pem", "", false, true},
		{"key only", "", "key.pem", false, true},
		{"nonexistent files", "nonexistent.pem", "nonexistent.key", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MCP_TLS_CERT_FILE", tt.certFile)
			os.Setenv("MCP_TLS_KEY_FILE", tt.keyFile)
			defer func() {
				os.Unsetenv("MCP_TLS_CERT_FILE")
				os.Unsetenv("MCP_TLS_KEY_FILE")
			}()

			config, err := GetTLSConfigFromEnv()
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, config)
			} else if tt.wantNil {
				require.NoError(t, err)
				require.Nil(t, config)
			} else {
				require.NoError(t, err)
				require.NotNil(t, config)
			}
		})
	}
}

func TestTLSConfigWithMismatchedKeyPair(t *testing.T) {
	// Generate two separate cert/key pairs
	certPEM1, _, err := generateTestCertAndKey()
	require.NoError(t, err)
	_, keyPEM2, err := generateTestCertAndKey()
	require.NoError(t, err)

	// Use cert from pair 1 with key from pair 2
	certPath := writeTempFile(t, "test_cert_*.pem", certPEM1)
	keyPath := writeTempFile(t, "test_key_*.pem", keyPEM2)

	os.Setenv("MCP_TLS_CERT_FILE", certPath)
	os.Setenv("MCP_TLS_KEY_FILE", keyPath)
	defer func() {
		os.Unsetenv("MCP_TLS_CERT_FILE")
		os.Unsetenv("MCP_TLS_KEY_FILE")
	}()

	config, err := GetTLSConfigFromEnv()
	require.Error(t, err)
	require.Nil(t, config)
	require.Contains(t, err.Error(), "invalid TLS certificate/key pair")
}

func TestTLSConfigWithInvalidCertContent(t *testing.T) {
	certPath := writeTempFile(t, "test_cert_*.pem", []byte("not a real cert"))
	keyPath := writeTempFile(t, "test_key_*.pem", []byte("not a real key"))

	os.Setenv("MCP_TLS_CERT_FILE", certPath)
	os.Setenv("MCP_TLS_KEY_FILE", keyPath)
	defer func() {
		os.Unsetenv("MCP_TLS_CERT_FILE")
		os.Unsetenv("MCP_TLS_KEY_FILE")
	}()

	config, err := GetTLSConfigFromEnv()
	require.Error(t, err)
	require.Nil(t, config)
	require.Contains(t, err.Error(), "invalid TLS certificate/key pair")
}

func TestTLSConfigCipherSuites(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertAndKey()
	require.NoError(t, err)

	certPath := writeTempFile(t, "test_cert_*.pem", certPEM)
	keyPath := writeTempFile(t, "test_key_*.pem", keyPEM)

	os.Setenv("MCP_TLS_CERT_FILE", certPath)
	os.Setenv("MCP_TLS_KEY_FILE", keyPath)
	defer func() {
		os.Unsetenv("MCP_TLS_CERT_FILE")
		os.Unsetenv("MCP_TLS_KEY_FILE")
	}()

	tlsConfig, err := GetTLSConfigFromEnv()
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)

	// Verify cipher suites are configured
	expectedCiphers := []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	}
	require.Equal(t, expectedCiphers, tlsConfig.Config.CipherSuites)

	// Verify curve preferences
	expectedCurves := []tls.CurveID{
		tls.X25519,
		tls.CurveP256,
		tls.CurveP384,
		tls.X25519MLKEM768,
	}
	require.Equal(t, expectedCurves, tlsConfig.Config.CurvePreferences)
}

func TestIsLocalHost(t *testing.T) {
	tests := []struct {
		host     string
		expected bool
	}{
		{"localhost", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"[::1]", true},
		{"0.0.0.0", true},
		{"LOCALHOST", true},
		{"Localhost", true},
		{"example.com", false},
		{"192.168.1.1", false},
		{"10.0.0.1", false},
		{"my-server.internal", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			result := IsLocalHost(tt.host)
			require.Equal(t, tt.expected, result)
		})
	}
}
