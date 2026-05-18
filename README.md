# awsmgr

`awsmgr` is an operational workflow and environment diagnostics CLI for engineers working with AWS SSO, multiple AWS profiles/accounts, Kubernetes/EKS, Terraform, and enterprise network environments.

It does not replace the AWS CLI. AWS CLI remains the source of truth for authentication, SSO, credentials, and profile resolution. `awsmgr` orchestrates visibility and safety checks around those workflows.

## Requirements

- Go 1.26+
- AWS CLI installed and configured
- macOS ARM64 first; Linux and Windows support can follow

## Build

```bash
make build
```

or:

```bash
go build -o bin/awsmgr .
```

## Installation via Homebrew (macOS)
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

- `list` lists AWS profiles discovered from AWS CLI configuration.
- `current` prints the active AWS profile from the local environment.
- `validate` checks profile sessions with `aws sts get-caller-identity`.
- `doctor` runs workstation diagnostics for AWS, Kubernetes, Terraform, and network access.
- `version` prints release metadata embedded at build time.

## Local Testing

```bash
make test
make build
./bin/awsmgr list
./bin/awsmgr current
./bin/awsmgr validate
./bin/awsmgr doctor
./bin/awsmgr version
```

To test custom local version metadata:

```bash
make build VERSION=1.0.2-local
./bin/awsmgr version
```

To build local release artifacts without publishing:

```bash
make release-snapshot
```

Local builds derive version metadata from Git tags. If `HEAD` is not on a stable
release tag, `make build` uses the next patch version with a development suffix,
for example `1.0.1-dev.cd8b901`.

To inspect or create the next release tag:

```bash
make version
make next-patch
make next-minor
make next-major
make tag-patch
```

Use `tag-minor` or `tag-major` when the release contains a feature or breaking
change. After creating a tag, push it with `git push origin vX.Y.Z` so the
release workflow can publish that SemVer version.

See [RELEASE.md](RELEASE.md) for the full commit, tagging, release, and rollback
workflow.

## Info

`validate` and `doctor` call local tools such as `aws`, `kubectl`, and `terraform`, so their results depend on your workstation configuration and active credentials.

## Author

**Santhosh John** - [GitHub](https://github.com/santhosh-john)

## License

MIT License - see [LICENSE](LICENSE) file for details.
