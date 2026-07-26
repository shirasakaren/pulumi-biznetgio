package provider

import (
	"context"
	"fmt"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

var Version string

const Name string = "biznetgio"

func Provider() p.Provider {
	p, err := infer.NewProviderBuilder().
		WithDisplayName("Pulumi BiznetGIO Provider").
		WithDescription("Pulumi provider for BiznetGIO cloud — NEO Metal, NEO Lite/Lite Pro, NEO GPU, and Object Storage.").
		WithHomepage("https://www.biznetgio.com").
		WithRepository("github.com/biznetgio/pulumi-biznetgio").
		WithPublisher("biznetgio").
		WithKeywords("biznetgio", "cloud", "indonesia", "neo").
		WithLicense("Apache-2.0").
		WithNamespace("biznetgio").
		WithLanguageMap(map[string]any{
			"nodejs": map[string]any{"packageName": "@biznetgio/biznetgio"},
			"python": map[string]any{"packageName": "pulumi_biznetgio"},
			"dotnet": map[string]any{"packageName": "Biznetgio.Biznetgio"},
			"java":   map[string]any{},
		}).
		WithGoImportPath("github.com/biznetgio/pulumi-biznetgio/sdk/go/pulumi-biznetgio").
		WithResources(
			// metal dulu
			infer.Resource(Baremetal{}),
			infer.Resource(BaremetalKeypair{}),
			infer.Resource(BaremetalAdditionalIp{}),
			infer.Resource(BaremetalAdditionalIpAssignment{}),
			infer.Resource(BaremetalElasticStorage{}),
