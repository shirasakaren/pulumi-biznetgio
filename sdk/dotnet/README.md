# Shirasakaren.Biznetgio

Unofficial [Pulumi](https://www.pulumi.com) provider for [BiznetGIO](https://www.biznetgio.com), an Indonesian
cloud platform. Manage NEO Metal (bare metal), NEO Lite / NEO Lite Pro (VMs), NEO GPU, and NEO Object Storage
(S3-compatible) as code.

This package is community maintained by [Shirasaka Ren](https://shirasaka.ren) and is not affiliated with or
endorsed by PT Biznet Gio Nusantara.

## Install

```bash
dotnet add package Shirasakaren.Biznetgio
```

Requires the [Pulumi](https://www.nuget.org/packages/Pulumi) package. The `pulumi-resource-biznetgio` plugin
binary downloads automatically from GitHub Releases on first use.

## Quickstart

```csharp
using Pulumi;
using Shirasakaren.Biznetgio;

return await Deployment.RunAsync(() =>
{
    var vm = new NeoliteVm("web", new NeoliteVmArgs
    {
        VmName = "web-1",
        ProductId = 123,
        SelectOs = "Ubuntu 22.04",
        Cycle = "m",
        ConsolePassword = "change-this-now!",
        SshAndConsoleUser = "root",
        // defaults to true: bills the card on file immediately.
        PayWithCreditCard = true,
    });

    return new Dictionary<string, object?> { ["status"] = vm.Status };
});
```

## Documentation

- Full docs and guides: [biznetgio.creations.ren](https://biznetgio.creations.ren)
- Source code: [github.com/shirasakaren/pulumi-biznetgio](https://github.com/shirasakaren/pulumi-biznetgio)

> **Billing note**: `payWithCreditCard` defaults to `true`, so the first `pulumi up` places a real order and may
> charge the card on file.

## License

Apache-2.0
