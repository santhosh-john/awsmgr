# Release Guide

This project uses SemVer tags, GoReleaser, and Homebrew formula publishing.

## Versioning

Use Semantic Versioning:

- Patch: fixes, docs, small CLI improvements, packaging fixes
- Minor: new commands or meaningful new user workflows
- Major: breaking command behavior, output contracts, or installation changes

Do not reuse or overwrite a published release tag. If a release is bad, roll
forward with a new patch release.

## Local Checks

Run these before committing substantial changes:

```bash
make fmt
make vet
make test
make build
./bin/awsmgr version
```

`make build` derives local version metadata from Git tags. If `HEAD` is not on a
clean stable release tag, the binary reports a development version such as:

```text
awsmgr 1.0.1-dev.cd8b901
commit: cd8b901
built: 2026-05-18T11:25:55Z
```

## Commit And Push

Review the changes:

```bash
git status
git diff
```

Commit and push:

```bash
git add .
git commit -m "Describe the change"
git push origin main
```

## Choose The Next Version

Inspect the next SemVer options:

```bash
make next-patch
make next-minor
make next-major
```

Most routine releases should be patch releases.

## Create A Release

Create the tag only from a clean worktree:

```bash
git status
make tag-patch
```

Use `make tag-minor` or `make tag-major` when appropriate.

Push the tag:

```bash
git push origin vX.Y.Z
```

The tag push should trigger the release workflow. GoReleaser builds binaries,
archives, checksums, the GitHub release, and the Homebrew formula update.

## Test Release Artifacts Locally

Build release artifacts without publishing:

```bash
make release-snapshot
```

Inspect the output:

```bash
ls dist
./dist/awsmgr_darwin_arm64_v8.0/awsmgr version
```

Or test the downloadable archive:

```bash
mkdir -p /tmp/awsmgr-test
tar -xzf dist/awsmgr_Darwin_arm64.tar.gz -C /tmp/awsmgr-test
/tmp/awsmgr-test/awsmgr version
```

The `dist/awsmgr_<os>_<arch>_<variant>` folders are GoReleaser build work
directories. The user-facing downloads are the archives, for example:

```text
awsmgr_Darwin_arm64.tar.gz
awsmgr_Darwin_x86_64.tar.gz
awsmgr_Linux_arm64.tar.gz
awsmgr_Linux_x86_64.tar.gz
awsmgr_Windows_x86_64.zip
```

## Rollback Before A Release

If a bad commit was pushed to `main` but no release tag was pushed, revert it:

```bash
git revert HEAD
git push origin main
```

For an older commit:

```bash
git revert <bad_commit_sha>
git push origin main
```

## Rollback After A Release

If a bad release was already published, do not reuse the same tag. Revert the
bad change and publish a new patch release:

```bash
git revert <bad_commit_sha>
make fmt
make vet
make test
make build
git add .
git commit -m "Revert broken change"
git push origin main

make tag-patch
git push origin vX.Y.Z
```

## Fix A Mistaken Local Tag

If the tag exists only locally:

```bash
git tag -d vX.Y.Z
```

If the mistaken tag was pushed but the release has not been consumed, delete it
with care:

```bash
git push origin :refs/tags/vX.Y.Z
git tag -d vX.Y.Z
```

Prefer a new patch release for normal rollback.
