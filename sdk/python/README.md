# Pulumi BiznetGIO Provider

> **Unofficial** community provider by [Shirasaka Ren](https://shirasaka.ren),
> not affiliated with or endorsed by PT Biznet Gio Nusantara.

A native [Pulumi](https://www.pulumi.com) provider for managing
[BiznetGIO](https://www.biznetgio.com) cloud infrastructure via the
[BiznetGIO Portal API](https://api.portal.biznetgio.com/v1/docs).

Covers NEO Metal (bare metal), NEO Lite / NEO Lite Pro (VMs), NEO GPU, and
NEO Object Storage (S3-compatible), with matching provider functions for
catalog lookups (products, OS lists, upgrade options, IP availability).

Documentation: https://biznetgio.creations.ren

## Install

SDK packages (the `pulumi-resource-biznetgio` plugin is downloaded
automatically from GitHub Releases via `pluginDownloadURL`):

| Language | Package |
|---|---|
| Node.js | `@shirasakaren/biznetgio` |
| Python | `shirasakaren-biznetgio` (import `pulumi_biznetgio`) |
| Go | `github.com/shirasakaren/pulumi-biznetgio/sdk/go/pulumi-biznetgio` |
| .NET | `Shirasakaren.Biznetgio` |
| Java | `com.shirasakaren.biznetgio` |

## Authentication

Set the API token via config (secret) or environment:

```bash
pulumi config set --secret biznetgio:apiToken <token>
# or: export BIZNETGIO_API_KEY=<token>
```

The token is sent as the `x-token` header. `baseUrl` defaults to
`https://api.portal.biznetgio.com/v1` (override with `BIZNETGIO_BASE_URL` or
`pulumi config set biznetgio:baseUrl ...`).

## Example

```typescript
import * as biznetgio from "@shirasakaren/biznetgio";

const plans = biznetgio.getNeoliteProducts();

const keypair = new biznetgio.NeoliteKeypair("deploy", { name: "deploy-key" });

const vm = new biznetgio.NeoliteVm("web", {
  vmName: "web-1",
  productId: plans.products[0].productId,
  selectOs: "Ubuntu 22.04",
  keypairId: keypair.id,
  cycle: "m",
  // defaults to true: the invoice is paid automatically with the stored card.
  // set false to keep the order pending until paid manually in the portal.
  payWithCreditCard: true,
});

export const vmStatus = vm.status;
```

> **Billing note**: every create/upgrade call places a real order and may
> charge the credit card on file. Resources created with
> `payWithCreditCard = false` stay `Pending` until the invoice is paid in the
> portal.

## Resources

`Baremetal`, `BaremetalKeypair`, `BaremetalAdditionalIp`,
`BaremetalAdditionalIpAssignment`, `BaremetalElasticStorage`, `GpuInstance`,
`GpuKeypair`, `NeoliteVm`, `NeoliteKeypair`, `NeoliteSnapshot`,
`NeoliteVmFromSnapshot`, `NeoliteDisk`, `NeoliteProVm`, `NeoliteProKeypair`,
`NeoliteProSnapshot`, `NeoliteProDisk`, `ObjectStorage`,
`ObjectStorageBucket`, `ObjectStorageCredential`, `ObjectStorageObject`.

## Functions

`baremetalProducts`, `baremetalRebuildOsList`, `baremetalOpenvpn`,
`gpuProducts`, `gpuConsole`, `gpuGraph`, `neoliteProducts`,
`neoliteOsList`, `neoliteChangePackageOptions`,
`neoliteStorageUpgradeOptions`, `neoliteIPAvailability`,
`neoliteProProducts`, `neoliteProOsList`,
`neoliteProChangePackageOptions`, `neoliteProStorageUpgradeOptions`,
`neoliteProIPAvailability`, `objectStorageInstances`,
`objectStorageBuckets`, `objectStorageCredentials`.

## Notes on the BiznetGIO API

- The Portal API does not publish response schemas. Response handling is
  defensive and was cross-checked against BiznetGIO's own SDKs and CLI; report
  any field mismatch as an issue. Every resource exposes a secret-marked `raw`
  output (secrets redacted) with the full last-read payload.
- Power actions are declarative (`powerState`); the API is only called when
  the value changes. One-shot actions (reset, rebuild, reserve GPU hours,
  migrate-to-pro) are trigger attributes: change the value to re-fire.
- Long-running provisioning is polled until active; on timeout the partial
  state is kept and the next `pulumi up` resumes the update.
- Products with no public API (NEO Virtual Compute, NEO Kubernetes, NEO DNS,
  domains, web hosting, gio-private, gio-enterprise-cloud, gio-backup) are out
  of scope.

## Development

```sh
make build install   # provider binary + all-language SDK codegen
make test_provider   # Go unit tests
make lint            # golangci-lint
```

## Publishing

Push a `v*.*.*` tag (or trigger the `release` workflow manually) and GitHub
Actions builds the SDKs, runs the language test matrix, and publishes:

- plugin binaries to **GitHub Releases** (Pulumi downloads them via
  `pluginDownloadURL: github://api.github.com/shirasakaren/pulumi-biznetgio`)
- `@shirasakaren/biznetgio` to npm, `shirasakaren-biznetgio` to PyPI,
  `Shirasakaren.Biznetgio` to NuGet, `com.shirasakaren.biznetgio` to Maven
  Central, and the Go SDK module to GitHub

**npm is live:** `npm install @shirasakaren/biznetgio` (latest `0.1.1`). See
[PUBLISHING.md](PUBLISHING.md) for the full release runbook covering every
registry.

Secrets required (GitHub repository secrets):

| Secret | Used for |
|---|---|
| `NPM_TOKEN` | npm publish |
| `PYPI_API_TOKEN` | PyPI publish |
| `NUGET_PUBLISH_KEY` | NuGet publish |
| `OSSRH_USERNAME` / `OSSRH_PASSWORD` | Maven Central (Sonatype) |
| `JAVA_SIGNING_KEY_ID` / `JAVA_SIGNING_KEY` / `JAVA_SIGNING_PASSWORD` | Maven artifact signing |
| `CODECOV_TOKEN` | coverage upload (optional) |
| `SLACK_WEBHOOK_URL` | failure notifications (optional) |

## License

Apache-2.0
