---
layout: default
title: Install
---

# Install

## Homebrew

```bash
brew install santhosh-john/tap/awsmgr
```

Upgrade:

```bash
brew update
brew upgrade awsmgr
```

Check the installed version:

```bash
awsmgr version
```

## Requirements

- macOS ARM64 first
- AWS CLI installed and configured
- Go 1.26+ for local development
- `kubectl` and `terraform` for full `doctor` diagnostics

## Direct Binaries

Direct binary downloads are published as GitHub release assets.

Release archive names follow this pattern:

```text
awsmgr_Darwin_arm64.tar.gz
awsmgr_Darwin_x86_64.tar.gz
awsmgr_Linux_arm64.tar.gz
awsmgr_Linux_x86_64.tar.gz
awsmgr_Windows_x86_64.zip
```

Download from:

```text
https://github.com/santhosh-john/awsmgr/releases
```
