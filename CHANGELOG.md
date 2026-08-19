## 0.1.0 (Unreleased)

FEATURES:

* Initial release of the BiznetGIO provider: 20 resources (`Baremetal`, `BaremetalKeypair`, `BaremetalAdditionalIp`, `BaremetalAdditionalIpAssignment`, `BaremetalElasticStorage`, `GpuInstance`, `GpuKeypair`, `NeoliteVm`, `NeoliteKeypair`, `NeoliteSnapshot`, `NeoliteVmFromSnapshot`, `NeoliteDisk`, `NeoliteProVm`, `NeoliteProKeypair`, `NeoliteProSnapshot`, `NeoliteProDisk`, `ObjectStorage`, `ObjectStorageBucket`, `ObjectStorageCredential`, `ObjectStorageObject`) and 19 functions (`baremetalProducts`, `baremetalRebuildOsList`, `baremetalOpenvpn`, `gpuProducts`, `gpuConsole`, `gpuGraph`, `neoliteProducts`, `neoliteOsList`, `neoliteChangePackageOptions`, `neoliteStorageUpgradeOptions`, `neoliteIPAvailability`, `neoliteProProducts`, `neoliteProOsList`, `neoliteProChangePackageOptions`, `neoliteProStorageUpgradeOptions`, `neoliteProIPAvailability`, `objectStorageInstances`, `objectStorageBuckets`, `objectStorageCredentials`).
* SDKs for Node.js (`@shirasakaren/biznetgio`), Python (`pulumi_biznetgio`), Go, .NET (`Shirasakaren.Biznetgio`), and Java (`com.pulumi:biznetgio`).
* API token auth via `apiToken` config secret or `BIZNETGIO_API_KEY`; `baseUrl` override via config or `BIZNETGIO_BASE_URL`.
