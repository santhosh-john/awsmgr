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

## Commands

```bash
awsmgr list
awsmgr current
awsmgr validate
awsmgr doctor
```

## Local Testing

```bash
make test
make build
./bin/awsmgr list
./bin/awsmgr current
./bin/awsmgr validate
./bin/awsmgr doctor
```

## Current Version: 1.0.2

`validate` and `doctor` call local tools such as `aws`, `kubectl`, and `terraform`, so their results depend on your workstation configuration and active credentials.

## Author

**Santhosh John** - [GitHub](https://github.com/santhosh-john)

## License

MIT License - see [LICENSE](LICENSE) file for details.

