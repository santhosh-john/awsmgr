---
layout: default
title: Troubleshooting
---

# Troubleshooting

## Homebrew Download Returns 404

Check that the repository is public and that the release asset exists with the
exact name used by the formula.

```bash
brew update
brew cat santhosh-john/tap/awsmgr
```

## AWS Profiles Are Missing

`awsmgr` reads profiles from AWS CLI configuration. Check:

```bash
aws configure list-profiles
```

## AWS Validation Fails

Refresh AWS SSO credentials:

```bash
aws sso login --profile <profile>
```

Then retry:

```bash
awsmgr validate
```

## Kubernetes Context Is Missing

Check:

```bash
kubectl config current-context
```

`awsmgr doctor` reports missing or unavailable `kubectl` clearly.

## Terraform Is Missing

Check:

```bash
terraform version
```

Install Terraform before expecting full `doctor` coverage.
