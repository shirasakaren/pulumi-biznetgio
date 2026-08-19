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
print("patched dotnet sdk")
