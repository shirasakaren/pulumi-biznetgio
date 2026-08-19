## 0.1.7 (Unreleased)

BUG FIXES:

* Fix the NuGet package: `assets/logo.png` was actually SVG content saved with a `.png` extension (from an
  earlier attempt to work around a 404), so NuGet.org rejected the whole package with "Unsupported icon image
  format." It's a real rasterized PNG now, referenced via the schema's `logoUrl`.
* Publish the release GPG key to public keyservers - Central Portal validates signatures against them, and
  without this every Maven Central deployment failed validation even with a correct signature.

## 0.1.6

BUG FIXES:

* Fix the Java SDK's package naming (was double-nested as
  `ren.shirasaka.biznetgio.ren.shirasaka_biznetgio`, now the correct `ren.shirasaka.biznetgio`).
* Fix the NuGet publish step: the shared release tooling only looks for `Pulumi.*.nupkg`, so it silently
  never pushed ours (`Shirasakaren.Biznetgio.*.nupkg`). It's pushed directly now.
* Fix the Maven Central publish bundle: it was missing per-artifact checksums and didn't preserve the Maven
  repo directory layout, so Central Portal accepted the upload but never synced it to Maven Central.
* Fix `schema.json` version drift and two golangci-lint line-length violations.

IMPROVEMENTS:

* Add `docs/_index.md` and `docs/installation-configuration.md` for the Pulumi Registry.
* Translate all resource/property descriptions and code comments to English.
* Remove an accidentally committed 40MB example binary.

## 0.1.1 - 0.1.5

* SDK publishing fixes across npm, PyPI, and Go (no schema or resource changes).

## 0.1.0 (Unreleased)

FEATURES:

* Initial release of the BiznetGIO provider: 20 resources (`Baremetal`, `BaremetalKeypair`, `BaremetalAdditionalIp`, `BaremetalAdditionalIpAssignment`, `BaremetalElasticStorage`, `GpuInstance`, `GpuKeypair`, `NeoliteVm`, `NeoliteKeypair`, `NeoliteSnapshot`, `NeoliteVmFromSnapshot`, `NeoliteDisk`, `NeoliteProVm`, `NeoliteProKeypair`, `NeoliteProSnapshot`, `NeoliteProDisk`, `ObjectStorage`, `ObjectStorageBucket`, `ObjectStorageCredential`, `ObjectStorageObject`) and 19 functions (`baremetalProducts`, `baremetalRebuildOsList`, `baremetalOpenvpn`, `gpuProducts`, `gpuConsole`, `gpuGraph`, `neoliteProducts`, `neoliteOsList`, `neoliteChangePackageOptions`, `neoliteStorageUpgradeOptions`, `neoliteIPAvailability`, `neoliteProProducts`, `neoliteProOsList`, `neoliteProChangePackageOptions`, `neoliteProStorageUpgradeOptions`, `neoliteProIPAvailability`, `objectStorageInstances`, `objectStorageBuckets`, `objectStorageCredentials`).
* SDKs for Node.js (`@shirasakaren/biznetgio`), Python (`pulumi_biznetgio`), Go, .NET (`Shirasakaren.Biznetgio`), and Java (`com.pulumi:biznetgio`).
* API token auth via `apiToken` config secret or `BIZNETGIO_API_KEY`; `baseUrl` override via config or `BIZNETGIO_BASE_URL`.
