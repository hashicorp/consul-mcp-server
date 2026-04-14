# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

schema = "2"

project "consul-mcp-server" {
  team = "namespaces/doormat/groups/github-hashicorp-team-rnd-consul-india"

  # slack channel : feed-consul-mcp-server-releases
  slack {
    notification_channel = "C09F6S18DMY"
  }

  github {
    organization     = "hashicorp"
    repository       = "consul-mcp-server"
    release_branches = ["main", "release/**"]
  }
}

event "merge" {
}

event "build" {
  action "build" {
    organization = "hashicorp"
    repository   = "consul-mcp-server"
    workflow     = "build"
    depends      = null
    config       = ""
  }

  depends = ["merge"]
}

event "prepare" {
  action "prepare" {
    organization = "hashicorp"
    repository   = "crt-workflows-common"
    workflow     = "prepare"
    depends      = ["build"]
    config       = ""
  }

  depends = ["build"]

  notification {
    on = "fail"
  }
}

event "trigger-staging" {
}

event "promote-staging" {
  action "promote-staging" {
    organization = "hashicorp"
    repository   = "crt-workflows-common"
    workflow     = "promote-staging"
    depends      = null
    config       = "oss-release-metadata.hcl"
  }

  depends = ["trigger-staging"]

  notification {
    on = "always"
  }

}

event "trigger-production" {
}

event "promote-production" {
  action "promote-production" {
    organization = "hashicorp"
    repository   = "crt-workflows-common"
    workflow     = "promote-production"
    depends      = null
    config       = ""
  }

  depends = ["trigger-production"]

  notification {
    on = "always"
  }

}
