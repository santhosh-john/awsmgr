---
layout: default
title: Commands
---

# Commands

## `awsmgr list`

Lists AWS profiles discovered from AWS CLI configuration.

```bash
awsmgr list
```

## `awsmgr current`

Prints the active AWS profile from the local environment.

```bash
awsmgr current
```

## `awsmgr validate`

Validates AWS profile sessions with:

```bash
aws sts get-caller-identity
```

Run:

```bash
awsmgr validate
```

## `awsmgr doctor`

Runs workstation diagnostics for AWS, Kubernetes, Terraform, and network access.

```bash
awsmgr doctor
```

## `awsmgr version`

Prints version metadata embedded at build or release time.

```bash
awsmgr version
```

Example:

```text
awsmgr 1.1.0
commit: ad78c11
built: 2026-05-18T13:09:35Z
```

## `awsmgr help`

Shows command help.

```bash
awsmgr help
awsmgr --help
```
