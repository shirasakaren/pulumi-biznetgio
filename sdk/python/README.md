<p align="center">
  <img src="./assets/logo.svg" width="88" alt="BiznetGIO logo" />
</p>

<h1 align="center">Pulumi BiznetGIO Provider</h1>

<p align="center">
  Manage <a href="https://www.biznetgio.com">BiznetGIO</a> cloud infrastructure - bare metal, VMs, GPUs, object storage - as code.
</p>

<p align="center">
  <a href="https://github.com/shirasakaren/pulumi-biznetgio/actions/workflows/ci.yml"><img src="https://github.com/shirasakaren/pulumi-biznetgio/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/shirasakaren/pulumi-biznetgio/releases/latest"><img src="https://img.shields.io/github/v/release/shirasakaren/pulumi-biznetgio" alt="Latest release"></a>
  <a href="https://www.npmjs.com/package/@shirasakaren/biznetgio"><img src="https://img.shields.io/npm/v/@shirasakaren/biznetgio?label=npm" alt="npm"></a>
  <a href="https://pypi.org/project/pulumi-biznetgio/"><img src="https://img.shields.io/pypi/v/pulumi-biznetgio?label=pypi" alt="PyPI"></a>
  <a href="https://www.nuget.org/packages/Shirasakaren.Biznetgio"><img src="https://img.shields.io/nuget/v/Shirasakaren.Biznetgio?label=nuget" alt="NuGet"></a>
  <a href="https://central.sonatype.com/artifact/ren.shirasaka/biznetgio"><img src="https://img.shields.io/maven-central/v/ren.shirasaka/biznetgio?label=maven" alt="Maven Central"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/shirasakaren/pulumi-biznetgio" alt="License"></a>
</p>

> **Unofficial** community provider maintained by [Shirasaka Ren](https://shirasaka.ren) - not affiliated with or
> endorsed by PT Biznet Gio Nusantara.

Covers **NEO Metal** (bare metal), **NEO Lite / NEO Lite Pro** (VMs), **NEO GPU**, and **NEO Object Storage**
(S3-compatible), backed by the [BiznetGIO Portal API](https://api.portal.biznetgio.com/v1/docs). New to Pulumi?
[pulumi.com/docs](https://www.pulumi.com/docs/) is the best starting point; this README assumes you already have
the Pulumi CLI installed.

📖 **Full docs, every resource, and per-language guides live at [biznetgio.creations.ren](https://biznetgio.creations.ren).**
This README is intentionally short.

## Install

<table>
<tr><td>Node.js</td><td><code>npm install @shirasakaren/biznetgio</code></td></tr>
<tr><td>Python</td><td><code>pip install pulumi-biznetgio</code></td></tr>
<tr><td>Go</td><td><code>go get github.com/shirasakaren/pulumi-biznetgio/sdk/go/pulumi-biznetgio</code></td></tr>
<tr><td>.NET</td><td><code>dotnet add package Shirasakaren.Biznetgio</code></td></tr>
<tr><td>Java</td><td><code>ren.shirasaka:biznetgio</code> on Maven Central</td></tr>
</table>

The plugin binary itself is downloaded automatically from GitHub Releases - no separate install step needed.

## Quickstart

```bash
pulumi config set --secret biznetgio:apiToken <your-token>   # or export BIZNETGIO_API_KEY
```

```typescript
import * as biznetgio from "@shirasakaren/biznetgio";

const keypair = new biznetgio.NeoliteKeypair("deploy", { name: "deploy-key" });

const vm = new biznetgio.NeoliteVm("web", {
    vmName: "web-1",
    productId: 123,
    selectOs: "Ubuntu 22.04",
    keypairId: keypair.keypairId,
    cycle: "m",
    consolePassword: "change-this-now!",
    sshAndConsoleUser: "root",
    payWithCreditCard: true, // ⚠️ bills the card on file immediately, see below
});

export const vmStatus = vm.status;
```

> **💳 Billing note**: `payWithCreditCard` defaults to `true`, so the first `pulumi up` places a real order and may
> charge the card on file. Set it to `false` to leave the order `Pending` until you pay manually in the portal.

Python, Go, C#, and Java examples: [docs site →](https://biznetgio.creations.ren)

<details>
<summary><b>Resources & functions</b> (click to expand)</summary>

**Resources**: `Baremetal`, `BaremetalKeypair`, `BaremetalAdditionalIp`, `BaremetalAdditionalIpAssignment`,
`BaremetalElasticStorage`, `GpuInstance`, `GpuKeypair`, `NeoliteVm`, `NeoliteKeypair`, `NeoliteSnapshot`,
`NeoliteVmFromSnapshot`, `NeoliteDisk`, `NeoliteProVm`, `NeoliteProKeypair`, `NeoliteProSnapshot`, `NeoliteProDisk`,
`ObjectStorage`, `ObjectStorageBucket`, `ObjectStorageCredential`, `ObjectStorageObject`.

**Functions**: `baremetalProducts`, `baremetalRebuildOsList`, `baremetalOpenvpn`, `gpuProducts`, `gpuConsole`,
`gpuGraph`, `neoliteProducts`, `neoliteOsList`, `neoliteChangePackageOptions`, `neoliteStorageUpgradeOptions`,
`neoliteIPAvailability`, `neoliteProProducts`, `neoliteProOsList`, `neoliteProChangePackageOptions`,
`neoliteProStorageUpgradeOptions`, `neoliteProIPAvailability`, `objectStorageInstances`, `objectStorageBuckets`,
`objectStorageCredentials`.

</details>

<details>
<summary><b>Notes on the BiznetGIO API</b> (click to expand)</summary>

- The Portal API doesn't publish response schemas. Response handling is defensive and cross-checked against
  BiznetGIO's own SDKs and CLI - every resource exposes a secret-marked `raw` output (secrets redacted) with the
  full last-read payload, so you're never blocked on an unmodeled field.
- Power actions are declarative (`powerState`); the API is only called when the value changes. One-shot actions
  (reset, rebuild, reserve GPU hours, migrate-to-pro) are trigger attributes - change the value to re-fire.
- Long-running provisioning is polled until active; on timeout the partial state is kept and the next `pulumi up`
  resumes the update.
- Out of scope: products with no public API (NEO Virtual Compute, NEO Kubernetes, NEO DNS, domains, web hosting,
  gio-private, gio-enterprise-cloud, gio-backup).

</details>

## Development

```sh
make build install   # provider binary + all-language SDK codegen
make test_provider    # Go unit tests
make lint             # golangci-lint
```

Releasing and per-registry setup (npm/PyPI/NuGet/Maven/Go) is documented in [PUBLISHING.md](PUBLISHING.md).

## License

[Apache-2.0](LICENSE)
