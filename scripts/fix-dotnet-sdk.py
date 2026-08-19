#!/usr/bin/env python3
"""post-codegen patch: dotnet codegen ignores packageName, always emits
Biznetgio.Biznetgio. rename to Shirasakaren.Biznetgio everywhere."""
import sys
from pathlib import Path

root = Path(sys.argv[1] if len(sys.argv) > 1 else "sdk/dotnet")
version = sys.argv[2] if len(sys.argv) > 2 else "0.0.0"

old = root / "Biznetgio.Biznetgio.csproj"
new = root / "Shirasakaren.Biznetgio.csproj"
if old.exists():
    old.rename(new)

csproj = new
t0 = csproj.read_text()
if "<Version>" not in t0:
    t0 = t0.replace("    <TargetFramework>", f"    <Version>{version}</Version>\n    <TargetFramework>")
    csproj.write_text(t0)

for p in list(root.rglob("*.cs")) + list(root.rglob("*.csproj")) + list(root.rglob("*.md")):
    t = p.read_text()
    if "Biznetgio.Biznetgio" in t or "github.com/shirasakaren" in t:
        t = t.replace("Biznetgio.Biznetgio", "Shirasakaren.Biznetgio")
        t = t.replace(">github.com/shirasakaren/pulumi-biznetgio<",
                      ">https://github.com/shirasakaren/pulumi-biznetgio<")
        p.write_text(t)

# wire a real package readme into the nupkg so nuget.org renders docs on the
# package page instead of the one-line codegen placeholder
t = csproj.read_text()
if "<PackageReadmeFile>" not in t:
    t = t.replace(
        "    <PackageIcon>logo.png</PackageIcon>",
        "    <PackageIcon>logo.png</PackageIcon>\n    <PackageReadmeFile>README.md</PackageReadmeFile>",
    )
    t = t.replace(
        "  <ItemGroup>\n    <None Include=\"logo.png\">",
        "  <ItemGroup>\n    <None Include=\"README.md\">\n      <Pack>True</Pack>\n      <PackagePath>\\</PackagePath>\n    </None>\n  </ItemGroup>\n\n  <ItemGroup>\n    <None Include=\"logo.png\">",
    )
    csproj.write_text(t)

README = """# Shirasakaren.Biznetgio

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
"""

(root / "README.md").write_text(README)

print("patched dotnet sdk")
