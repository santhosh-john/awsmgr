# Development

## Local Checks

Run these before committing substantial Go changes:

```bash
make fmt
make vet
make test
make build
```

Check the local binary:

```bash
./bin/awsmgr version
./bin/awsmgr help
```

## Versioning

Local builds derive version metadata from Git tags.

Inspect the next versions:

```bash
make version
make next-patch
make next-minor
make next-major
```

Create release tags from a clean worktree:

```bash
make tag-patch
make tag-minor
make tag-major
```

## Release Artifacts

Build local release artifacts without publishing:

```bash
make release-snapshot
```

Inspect:

```bash
ls dist
```

Run the local macOS ARM64 release binary:

```bash
./dist/awsmgr_darwin_arm64_v8.0/awsmgr version
```

## More

See:

- Release guide: https://github.com/santhosh-john/awsmgr/blob/main/RELEASE.md
- Website and documentation guide: https://github.com/santhosh-john/awsmgr/blob/main/DOCS_WEBSITE.md
