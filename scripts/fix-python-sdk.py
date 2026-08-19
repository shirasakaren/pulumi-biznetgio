#!/usr/bin/env python3
"""post-codegen patch: python codegen hardcodes VERSION = "0.0.0" in setup.py.
substitute the real version before building."""
import sys
from pathlib import Path

root = Path(sys.argv[1] if len(sys.argv) > 1 else "sdk/python")
version = sys.argv[2] if len(sys.argv) > 2 else "0.0.0"

p = root / "setup.py"
t = p.read_text()
t = t.replace('VERSION = "0.0.0"', f'VERSION = "{version}"', 1)
t = t.replace("'Repository': 'github.com/shirasakaren/pulumi-biznetgio'",
              "'Repository': 'https://github.com/shirasakaren/pulumi-biznetgio'")
p.write_text(t)
print(f"patched python setup.py -> {version}")
