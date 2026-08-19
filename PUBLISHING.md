# Publishing Runbook — pulumi-biznetgio

Complete guide to get this community provider published everywhere users can
consume it. Maintained by Shirasaka Ren under the `shirasakaren` namespace.

## Current status

| Target | Status |
| --- | --- |
| npm (`@shirasakaren/biznetgio`) | ✅ **LIVE** — latest `0.1.2` (`npm install @shirasakaren/biznetgio`) |
| GitHub Releases (plugin binaries) | ✅ **LIVE** — `v0.1.2` (plugin auto-downloads) |
| PyPI (`pulumi-biznetgio`) | ⬜ not yet |
| NuGet (`Shirasakaren.Biznetgio`) | ⬜ not yet |
| Maven Central (`ren.shirasaka:biznetgio`) | ⬜ publishing via the Central Portal plugin |
| Go module | ✅ resolved from GitHub automatically |
| Pulumi Registry listing | ⬜ optional discovery layer — submit after everything else |

## Architecture

```
                    Pulumi Registry (discovery + docs, optional)
                                 │
                                 ▼
        github.com/shirasakaren/pulumi-biznetgio
                                 │
          ┌──────────────────────┼──────────────────────┐
          ▼                      ▼                      ▼
   GitHub Releases         SDK packages            schema.json
   pulumi-resource-        (npm/PyPI/NuGet/        (fed to the
   biznetgio binaries      Maven/Go module)        Pulumi Registry)
```

Users never need the Pulumi Registry: Pulumi pulls the plugin binary from
**GitHub Releases** (`pluginDownloadURL` in the schema), and each language
installs its SDK from its own registry.

## Do you need PyPI / NuGet / Maven?

Yes — each is a separate ecosystem:

- **npm** covers TypeScript/JavaScript only (done).
- **PyPI** covers Python users (`import pulumi_biznetgio`).
- **NuGet** covers C#/.NET users.
- **Maven Central** covers Java users.
- **Go** needs no upload — Go resolves the module straight from this GitHub repo via tags.

Skipping one registry simply means that language can't install the SDK.

---

## 1. GitHub environment and secrets (one-time)

Create a GitHub environment named `release` (repo → Settings → Environments):

- Required reviewers: your account
- Deployment branches/tags: tags only (`v*.*.*`)

Add these **repository secrets** as they become available:

| Secret | Used for | Required now? |
| --- | --- | --- |
| `NPM_TOKEN` | npm publish (CI) | yes |
| `PYPI_API_TOKEN` | PyPI publish | when doing Python |
| `NUGET_PUBLISH_KEY` | NuGet publish | when doing .NET |
| `OSSRH_USERNAME` / `OSSRH_PASSWORD` | Maven Central upload | when doing Java |
| `JAVA_SIGNING_KEY_ID` / `JAVA_SIGNING_KEY` / `JAVA_SIGNING_PASSWORD` | Maven artifact signing | when doing Java |
| `CODECOV_TOKEN` | coverage upload | optional |
| `SLACK_WEBHOOK_URL` | failure notifications | optional |

Add the token you used for the manual npm publish as `NPM_TOKEN` so future
releases automate npm automatically.

<Note>
  Prefer OIDC / Trusted Publishing where supported (npm and PyPI both support
  it) so no long-lived token sits in the repo. Instructions below.
</Note>

## 2. npm — automation (publish is live)

The `release.yml` workflow publishes npm automatically on every `v*.*.*` tag
using the `NPM_TOKEN` secret. To replace the token with **npm Trusted
Publishing**:

1. npm: package page → Settings → Add GitHub Repo → owner `shirasakaren`,
   repo `pulumi-biznetgio`, workflow `release.yml`, environment `release`.
2. The publish job then runs with `id-token: write` and no `NPM_TOKEN`.

Manual publish (for reference):

```bash
cd sdk/nodejs && npm install && npm run build
cp ../../README.md ../../LICENSE package.json bin/
cd bin && npm pack --dry-run          # inspect the tarball
npm publish --access public
```

## 3. PyPI

Package: `pulumi-biznetgio` (import stays `pulumi_biznetgio` — PyPI normalizes `_` to `-`).

1. Create a [PyPI](https://pypi.org) account (2FA on).
2. Reserve the name by publishing once, or use a Trusted Publisher directly:
   - PyPI → Publishing → Add a new pending publisher: owner `shirasakaren`,
     repository `pulumi-biznetgio`, workflow `release.yml`, environment `release`.
3. Build locally to sanity-check:

   ```bash
   cd sdk/python && python -m build
   ls dist/                    # pulumi_biznetgio-0.1.2-py3-none-any.whl
   ```

4. The release workflow publishes via `pypa/gh-action-pypi-publish` with
   `id-token: write` — no token needed once the Trusted Publisher exists.

## 4. NuGet

Package ID: `Shirasakaren.Biznetgio`.

1. Create a [nuget.org](https://www.nuget.org) account (2FA on).
2. Generate an API key (scoped to push this package ID) and store it as the
   `NUGET_PUBLISH_KEY` secret.
3. Sanity-check the package locally:

   ```bash
   cd sdk/dotnet && dotnet pack
   ls bin/Release/*.nupkg
   ```

   Verify inside the `.nupkg`: Package ID, version, description, repository
   URL, license, authors.
4. Release workflow pushes with `dotnet nuget push`.

## 5. Maven Central

Coordinates: `ren.shirasaka:biznetgio` (group `ren.shirasaka`, artifact
`biznetgio`; generated classes live under base package
`ren.shirasaka.biznetgio`). `ren.shirasaka` is the reverse-DNS of the
`shirasaka.ren` domain you own, so the namespace is verified at Sonatype via
a DNS TXT record on that domain (see §13 for the exact steps).

1. Create an account on the [Sonatype Central Portal](https://central.sonatype.com)
   and **verify the `ren.shirasaka` namespace** with a DNS TXT record on
   `shirasaka.ren` (exact steps in §13).
2. Generate a GPG keypair for signing and store:
   - `OSSRH_USERNAME`, `OSSRH_PASSWORD` (Central Portal credentials)
   - `JAVA_SIGNING_KEY_ID`, `JAVA_SIGNING_KEY`, `JAVA_SIGNING_PASSWORD` (GPG)
3. The Java SDK ships `settings.gradle` + `build.gradle` (gradle-nexus: sources
   + javadoc JARs, in-memory GPG signing from env vars); the workflow runs:

   ```bash
   gradle -p ./sdk/java publishToSonatype closeAndReleaseSonatypeStagingRepository
   ```

   Do this step early — Maven has the most account/signing setup.

## 6. Go SDK

No registry to upload. The module
`github.com/shirasakaren/pulumi-biznetgio/sdk/go/pulumi-biznetgio` resolves
from this repo. Every release must create **two tags**:

```bash
git tag v0.1.0
git tag sdk/go/pulumi-biznetgio/v0.1.0
```

The release workflow's `publish_go_sdk` job creates the SDK tag
automatically. Users install:

```bash
go get github.com/shirasakaren/pulumi-biznetgio/sdk/go/pulumi-biznetgio@v0.1.0
```

## 7. GitHub Releases + plugin binaries

Tagging `v0.1.0` triggers `release.yml`: tests → SDK builds → GoReleaser →
GitHub Release with `pulumi-resource-biznetgio-v0.1.0-{linux,darwin,windows}-{amd64,arm64}`
archives. Pulumi downloads these automatically because the schema contains:

```json
"pluginDownloadURL": "github://api.github.com/shirasakaren/pulumi-biznetgio"
```

Verify from a clean machine after the first release:

```bash
pulumi plugin install resource biznetgio 0.1.0
pulumi plugin ls
```

## 8. Pulumi Registry listing (optional, do last)

The Registry is a discovery/docs layer — installs work without it.

1. Fork [pulumi/registry](https://github.com/pulumi/registry).
2. Add to `community-packages/package-list.json`:

   ```json
   {
     "repoSlug": "shirasakaren/pulumi-biznetgio",
     "schemaFile": "provider/cmd/pulumi-resource-biznetgio/schema.json"
   }
   ```

3. Open a PR. Publisher identity shows as `shirasakaren` (from provider
   metadata) — correct for a community provider.
4. Update this README with the Registry URL after acceptance.

## 9. First release — the exact order

1. All identity checks (done): no `github.com/biznetgio` / `@biznetgio`
   references anywhere in the repo.
2. Commit and push everything to `main`.
3. Create the `release` environment + secrets above.
4. Local rehearsal:

   ```bash
   make build
   make test_provider
   make codegen && make build_sdks
   git status           # only intentional generated changes
   jq . provider/cmd/pulumi-resource-biznetgio/schema.json >/dev/null
   ```

5. Tag and push:

   ```bash
   git tag -a v0.1.0 -m "Pulumi BiznetGIO Provider v0.1.0"
   git push origin v0.1.0
   ```

6. Watch the workflow; fix failures and re-tag `v0.1.1` if needed.
7. Run the verification matrix below from clean environments.
8. Submit the Pulumi Registry PR.

## 10. Verification matrix (clean environment per ecosystem)

| Ecosystem | Test |
| --- | --- |
| npm | `npm install @shirasakaren/biznetgio` + `import * as biznetgio from "@shirasakaren/biznetgio"` |
| PyPI | `pip install pulumi-biznetgio` + `import pulumi_biznetgio` |
| Go | `go get github.com/shirasakaren/pulumi-biznetgio/sdk/go/pulumi-biznetgio@v0.1.0` |
| NuGet | `dotnet add package Shirasakaren.Biznetgio` |
| Maven | `implementation("ren.shirasaka:biznetgio:0.1.0")` |
| Plugin | `pulumi plugin install resource biznetgio 0.1.0` then `pulumi up` in an example |

## 11. Security checklist

- 2FA on: GitHub, npm, PyPI, NuGet, Sonatype.
- Branch protection on `main` and tags `v*` + `sdk/go/pulumi-biznetgio/v*`.
- Prefer Trusted Publishing (npm, PyPI) over static tokens.
- Normal CI keeps `permissions: contents: read`; only the release job gets
  `contents: write` / `id-token: write`.
- Never commit tokens; the npm token used for the first publish is a
  publish-only token — store it as `NPM_TOKEN` or delete it after enabling
  Trusted Publishing.

## 12. Versioning

Start at `v0.1.0`, follow semver, and do not rush `v1.0.0` — the provider
schema is the public API contract. Each release bumps the SDK versions in
lockstep via the release workflow.

## 13. Sonatype namespace verification — user steps (blocking)

This is the only step that requires the Sonatype web UI (no API exists), so
do it yourself:

1. Sign in at https://central.sonatype.com with your account.
2. Namespaces → **Add Namespace** → `ren.shirasaka`.
3. Pick verification method **DNS record** — the namespace is the
   reverse-DNS of `shirasaka.ren`, the domain you own.
4. The portal shows a TXT record to publish. Add it in your DNS for
   `shirasaka.ren` (use the exact host/value it displays) and click verify.
5. Once verified, generate a **publish token** (User → Publishing tokens →
   Access tokens). Store:
   - token username → `OSSRH_USERNAME` secret
   - token password → `OSSRH_PASSWORD` secret
6. The Java publish job already has the GPG signing secrets configured, so
   the next `v*` tag publishes `ren.shirasaka:biznetgio` to Maven Central.

Until verification completes, the Java publish step fails on namespace
permission — everything else in the release pipeline is unaffected.
