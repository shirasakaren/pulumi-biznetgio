package provider

import (
	"context"
	"fmt"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
)

var Version string

const Name string = "biznetgio"

func Provider() p.Provider {
	p, err := infer.NewProviderBuilder().
		WithDisplayName("Pulumi BiznetGIO Provider").
		WithDescription("Unofficial Pulumi provider for Biznet GIO cloud by Shirasaka Ren — NEO Metal, NEO Lite/Lite Pro, NEO GPU, and Object Storage.").
		WithHomepage("https://biznetgio.creations.ren").
		WithRepository("github.com/shirasakaren/pulumi-biznetgio").
		WithPublisher("shirasakaren").
		WithKeywords("biznetgio", "cloud", "indonesia", "neo").
		WithLicense("Apache-2.0").
		WithNamespace("biznetgio").
		WithPluginDownloadURL("github://api.github.com/shirasakaren/pulumi-biznetgio").
		WithLogoURL("https://raw.githubusercontent.com/shirasakaren/pulumi-biznetgio/main/assets/logo.svg").
		WithLanguageMap(map[string]any{
			"nodejs": map[string]any{"packageName": "@shirasakaren/biznetgio"},
			"python": map[string]any{"packageName": "pulumi_biznetgio"},
			"dotnet": map[string]any{"packageName": "Shirasakaren.Biznetgio"},
			"java":   map[string]any{"packages": map[string]any{"biznetgio": "ren.shirasaka:biznetgio"}, "basePackage": "ren.shirasaka.biznetgio"},
		}).
		WithGoImportPath("github.com/shirasakaren/pulumi-biznetgio/sdk/go/pulumi-biznetgio").
		WithResources(
			// metal dulu
			infer.Resource(Baremetal{}),
			infer.Resource(BaremetalKeypair{}),
			infer.Resource(BaremetalAdditionalIp{}),
			infer.Resource(BaremetalAdditionalIpAssignment{}),
			infer.Resource(BaremetalElasticStorage{}),
			// gpu cekidot
			infer.Resource(GpuInstance{}),
			infer.Resource(GpuKeypair{}),
			// neolite gaskeun
			infer.Resource(NeoliteVm{}),
			infer.Resource(NeoliteKeypair{}),
			infer.Resource(NeoliteSnapshot{}),
			infer.Resource(NeoliteVmFromSnapshot{}),
			infer.Resource(NeoliteDisk{}),
			// pro gacor
			infer.Resource(NeoliteProVm{}),
			infer.Resource(NeoliteProKeypair{}),
			infer.Resource(NeoliteProSnapshot{}),
			infer.Resource(NeoliteProDisk{}),
			// object storage jos
			infer.Resource(ObjectStorage{}),
			infer.Resource(ObjectStorageBucket{}),
			infer.Resource(ObjectStorageCredential{}),
			infer.Resource(ObjectStorageObject{}),
		).
		WithFunctions(
			// metal dulu
			infer.Function(BaremetalProducts{}),
			infer.Function(BaremetalRebuildOsList{}),
			infer.Function(BaremetalOpenvpn{}),
			// gpu cekidot
			infer.Function(GpuProducts{}),
			infer.Function(GpuConsole{}),
			infer.Function(GpuGraph{}),
			// neolite gaskeun
			infer.Function(NeoliteProducts{}),
			infer.Function(NeoliteOsList{}),
			infer.Function(NeoliteChangePackageOptions{}),
			infer.Function(NeoliteStorageUpgradeOptions{}),
			infer.Function(NeoliteIPAvailability{}),
			// pro gacor
			infer.Function(NeoliteProProducts{}),
			infer.Function(NeoliteProOsList{}),
			infer.Function(NeoliteProChangePackageOptions{}),
			infer.Function(NeoliteProStorageUpgradeOptions{}),
			infer.Function(NeoliteProIPAvailability{}),
			// object storage jos
			infer.Function(ObjectStorageInstances{}),
			infer.Function(ObjectStorageBuckets{}),
			infer.Function(ObjectStorageCredentials{}),
		).
		WithConfig(infer.Config(&Config{})).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		}).Build()
	if err != nil {
		panic(fmt.Errorf("unable to build provider: %w", err))
	}
	return p
}

type Config struct {
	ApiToken *string `pulumi:"apiToken,optional" provider:"secret"`
	BaseURL  *string `pulumi:"baseUrl,optional"`

	client *client.Client
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.ApiToken, "BiznetGIO API token (x-token header). Falls back to BIZNETGIO_API_KEY.")
	a.SetDefault(&c.ApiToken, nil, "BIZNETGIO_API_KEY")
	a.Describe(&c.BaseURL, "BiznetGIO API base URL. Falls back to BIZNETGIO_BASE_URL.")
	a.SetDefault(&c.BaseURL, "https://api.portal.biznetgio.com/v1", "BIZNETGIO_BASE_URL")
}

func (c *Config) Configure(_ context.Context) error {
	if c.ApiToken == nil || *c.ApiToken == "" {
		return fmt.Errorf("apiToken is required (set via `pulumi config set --secret biznetgio:apiToken <token>` " +
			"or BIZNETGIO_API_KEY)")
	}
	c.client = client.New(*c.BaseURL, *c.ApiToken, 30*time.Second)
	return nil
}

func GetClient(ctx context.Context) *client.Client {
	return infer.GetConfig[Config](ctx).client
}
