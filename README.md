# Documentation Site Workflow

This directory is the source for the public `awsmgr` GitHub Pages website.

The website is published from this repository to the `gh-pages` branch by
`.github/workflows/docs-pages.yml`.

## First-Time Setup

In GitHub repository settings:

1. Open `Settings`.
2. Open `Pages`.
3. Set `Build and deployment` source to `Deploy from a branch`.
4. Select branch `gh-pages`.
5. Select folder `/`.
6. Save.

The expected website URL is:

```text
https://santhosh-john.github.io/awsmgr/
```

## Edit Docs

Change files under `docs/`:

```text
docs/index.md
docs/install.md
docs/commands.md
docs/troubleshooting.md
docs/development.md
docs/release-notes.md
```

Then commit to `main`:

```bash
git status
git add docs README.md DOCS_WEBSITE.md
git commit -m "docs: update website"
git push origin main
```

The docs workflow publishes the site to `gh-pages`.

## Documentation Release Tags

Documentation tags are separate from CLI release tags.

CLI release tags:

```text
v1.1.0
v1.1.1
v1.2.0
```

Documentation release tags:

```text
docs-v2026.05.18
docs-v2026.05.18.1
```

Create a documentation tag only after the docs have been pushed to `main`:

```bash
git tag -a docs-vYYYY.MM.DD -m "Documentation release docs-vYYYY.MM.DD"
git push origin docs-vYYYY.MM.DD
```

If publishing more than once on the same day:

```bash
git tag -a docs-vYYYY.MM.DD.1 -m "Documentation release docs-vYYYY.MM.DD.1"
git push origin docs-vYYYY.MM.DD.1
```

## Local Review

The docs are plain Markdown. Review them directly in your editor or through
GitHub's Markdown preview.

If you want to preview them with a local static server:

```bash
cd docs
python3 -m http.server 8000
```

Then open:

```text
http://localhost:8000/
```

## Rollback

If a docs change is wrong and no docs tag was pushed:

```bash
git revert HEAD
git push origin main
```

If a docs tag was already pushed, revert and publish a new docs tag:

```bash
git revert <bad_docs_commit_sha>
git push origin main
git tag -a docs-vYYYY.MM.DD.N -m "Documentation rollback docs-vYYYY.MM.DD.N"
git push origin docs-vYYYY.MM.DD.N
```

Do not create CLI release tags for documentation-only changes.
