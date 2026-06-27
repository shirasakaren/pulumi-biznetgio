package main

import (
	"fmt"

	biznetgio "github.com/biznetgio/pulumi-biznetgio/sdk/go/pulumi-biznetgio"
	bngcfg "github.com/biznetgio/pulumi-biznetgio/sdk/go/pulumi-biznetgio/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// token: `pulumi config set biznetgio:apiToken <token> --secret` atau env BIZNETGIO_API_KEY
		if bngcfg.GetApiToken(ctx) == "" {
			return fmt.Errorf("set biznetgio:apiToken di config atau BIZNETGIO_API_KEY di env")
		}

		products, err := biznetgio.NeoliteProducts(ctx, &biznetgio.NeoliteProductsArgs{})
		if err != nil {
			return err
		}
		osList, err := biznetgio.NeoliteOsList(ctx, &biznetgio.NeoliteOsListArgs{
			ProductId: products.Products[0].ProductId,
		})
		if err != nil {
			return err
		}

		keypair, err := biznetgio.NewNeoliteKeypair(ctx, "demo-keypair", &biznetgio.NeoliteKeypairArgs{
			Name: pulumi.String("neo-lite-key"),
		})
		if err != nil {
			return err
		}

		vm, err := biznetgio.NewNeoliteVm(ctx, "demo-vm", &biznetgio.NeoliteVmArgs{
			VmName:            pulumi.String("neo-lite-1"),
			ProductId:         pulumi.Int(products.Products[0].ProductId),
			SelectOs:          pulumi.String(osList.Oss[0].Name),
			KeypairId:         keypair.KeypairId,
			Cycle:             pulumi.String("m"),
			SshAndConsoleUser: pulumi.String("admin"),
			ConsolePassword:   pulumi.String("s3cretP4ssw0rd"),
		})
		if err != nil {
			return err
		}

		ctx.Export("vmStatus", vm.Status)
		return nil
	})
}
