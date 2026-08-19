#!/usr/bin/env python3
"""post-codegen patch: nodejs codegen writes ${VERSION} placeholder and no
entrypoints. substitute the real version and add main/types/files."""
import json
import sys
from pathlib import Path

root = Path(sys.argv[1] if len(sys.argv) > 1 else "sdk/nodejs")
version = sys.argv[2] if len(sys.argv) > 2 else "0.0.0"

p = root / "package.json"
d = json.loads(p.read_text())
d["version"] = version
d["main"] = "./index.js"
d["types"] = "./index.d.ts"
d["files"] = ["**/*"]
d.setdefault("homepage", "https://biznetgio.creations.ren")
d["repository"] = "https://github.com/shirasakaren/pulumi-biznetgio"
p.write_text(json.dumps(d, indent=4) + "\n")
print(f"patched nodejs package.json -> {version}")
