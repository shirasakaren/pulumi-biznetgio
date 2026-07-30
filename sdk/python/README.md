# Pulumi BiznetGIO Provider

A native [Pulumi](https://www.pulumi.com) provider for managing
[BiznetGIO](https://www.biznetgio.com) cloud infrastructure via the
[BiznetGIO Portal API](https://api.portal.biznetgio.com/v1/docs).

Covers NEO Metal (bare metal), NEO Lite / NEO Lite Pro (VMs), NEO GPU, and
NEO Object Storage (S3-compatible), with matching provider functions for
catalog lookups (products, OS lists, upgrade options, IP availability).

## Install

The provider is published as `biznetgio` in the Pulumi Registry. SDK packages:

| Language | Package |
|---|---|
| Node.js | `@biznetgio/biznetgio` |
| Python | `pulumi_biznetgio` |
| Go | `github.com/biznetgio/pulumi-biznetgio/sdk/go/pulumi-biznetgio` |
| .NET | `Biznetgio.Biznetgio` |
| Java | `com.pulumi:biznetgio` |

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
import * as biznetgio from "@biznetgio/biznetgio";

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

`make ci-mgmt` regenerates release workflows from `.ci-mgmt.yaml`. Before
publishing, replace the placeholder `biznetgio` GitHub org with the real one
and set `publishRegistry: true`. Release secrets come from a Pulumi ESC
environment (`imports/github-secrets` in `.ci-mgmt.yaml`).

## License

Apache-2.0
