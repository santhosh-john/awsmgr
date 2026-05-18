---
layout: default
title: Home
---
# awsmgr

<p class="lead">
An operational workflow and environment diagnostics CLI for
engineers working with AWS SSO, multiple AWS profiles, Kubernetes/EKS,
Terraform, and enterprise network environments.

It does not replace AWS CLI. AWS CLI remains the source of truth for
authentication, SSO, credentials, and profile resolution.
</p>

## Install

```bash
brew install santhosh-john/tap/awsmgr
```

## Commands

```bash
awsmgr list
awsmgr current
awsmgr validate
awsmgr doctor
awsmgr version
awsmgr help
```

## Documentation

- [Install](install.html)
- [Commands](commands.html)
- [Troubleshooting](troubleshooting.html)
- [Development](development.html)
- [Release notes](release-notes.html)
