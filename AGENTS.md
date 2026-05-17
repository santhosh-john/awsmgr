# AGENTS

## Project Focus

`awsmgr` is an operational workflow and environment diagnostics CLI for engineers working with AWS SSO, multiple AWS profiles, Kubernetes/EKS, Terraform, and enterprise network environments.

- Do not turn this into an AWS CLI replacement.
- AWS CLI remains the source of truth for authentication, SSO, credentials, and profile resolution.
- Prefer workflow orchestration, visibility, validation, and safety checks over deeper AWS API coverage.

See [README.md](README.md) for product intent and local usage.

## Build And Validation

Run these before finishing substantial Go changes:

- `make fmt`
- `make vet`
- `make test`
- `make build`

CI currently enforces format, vet, and tests on macOS with Go 1.26.x. See [.github/workflows/ci.yml](.github/workflows/ci.yml).

## Architecture

- [main.go](main.go) owns CLI argument routing and the top-level command timeout.
- [internal/aws](internal/aws) owns AWS profile discovery and AWS CLI based validation.
- [internal/doctor](internal/doctor) orchestrates environment diagnostics across packages.
- [internal/kube](internal/kube) owns `kubectl` detection and current-context lookup.
- [internal/terraform](internal/terraform) owns Terraform availability checks.
- [internal/network](internal/network) owns basic reachability checks.
- [internal/ui](internal/ui) owns terminal output formatting.

Keep packages small and focused. Add new packages only when a responsibility is clearly separate.

## Go Conventions For This Repo

- Prefer the standard library first.
- Avoid unnecessary interfaces and dependency injection patterns.
- Keep data flow explicit with typed structs and direct function calls.
- Use `context.Context` for command execution and time-bounded checks.
- When validating multiple profiles, preserve the existing goroutine, channel, and `sync.WaitGroup` style instead of introducing heavier abstractions.
- Wrap errors with context using `%w`.
- Add comments only where the intent is not obvious from the code.

## External Tooling Assumptions

- Primary target is macOS ARM64 for now.
- `validate` depends on `aws sts get-caller-identity` succeeding for each profile.
- `doctor` depends on local workstation tools and environment state, especially `aws`, `kubectl`, and `terraform`.
- Missing external tools should be reported clearly, not hidden.

When changing diagnostics, prefer graceful failure with actionable messages over fallback behavior that obscures the real issue.

## Change Guardrails

- Do not add AWS SDK usage unless there is a strong reason that AWS CLI orchestration cannot satisfy the requirement.
- Do not implement future roadmap commands unless explicitly asked.
- Preserve human-readable CLI output unless the change explicitly introduces a machine-readable mode.
- Keep command behavior composable and easy to extend.

## Release And Distribution Notes

- The repo includes a starter GoReleaser config in [.goreleaser.yaml](.goreleaser.yaml).
- Current distribution targets are Homebrew and direct binaries.
- If release-related changes are made, ensure they still build for `darwin`, `linux`, and `windows` with `arm64` and `amd64` as configured.

## Useful Files To Read First

- [README.md](README.md)
- [main.go](main.go)
- [internal/aws/profiles.go](internal/aws/profiles.go)
- [internal/aws/validate.go](internal/aws/validate.go)
- [internal/doctor/doctor.go](internal/doctor/doctor.go)
- [Makefile](Makefile)