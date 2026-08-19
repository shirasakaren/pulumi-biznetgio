#!/usr/bin/env python3
"""post-codegen patch: java codegen derives group from basePackage, we want
group ren.shirasaka + artifact biznetgio + basePackage ren.shirasaka.biznetgio."""
import re
import sys
from pathlib import Path

root = Path(sys.argv[1] if len(sys.argv) > 1 else "sdk/java")
bg = root / "build.gradle"
sg = root / "settings.gradle"

t = bg.read_text()
t = t.replace('group = "ren.shirasaka.biznetgio"', 'group = "ren.shirasaka"')
t = t.replace('groupId = "ren.shirasaka.biznetgio"', 'groupId = "ren.shirasaka"')
t = t.replace('url = "github.com/shirasakaren/pulumi-biznetgio"',
              'url = "https://github.com/shirasakaren/pulumi-biznetgio"')
t = t.replace('connection = "github.com/shirasakaren/pulumi-biznetgio"',
              'connection = "scm:git:https://github.com/shirasakaren/pulumi-biznetgio.git"')
t = t.replace('developerConnection = "github.com/shirasakaren/pulumi-biznetgio"',
              'developerConnection = "scm:git:https://github.com/shirasakaren/pulumi-biznetgio.git"')
t = t.replace('id = ""\n                        name = ""\n                        email = ""',
              'id = "shirasakaren"\n                        name = "Shirasaka Ren"\n                        email = "ren@shirasaka.ren"')
t = t.replace('name = ""\n                packaging = "jar"',
              'name = "pulumi-biznetgio"\n                packaging = "jar"')
bg.write_text(t)

s = sg.read_text()
s = s.replace('rootProject.name = "ren.shirasaka.biznetgio.biznetgio"',
              'rootProject.name = "ren.shirasaka.biznetgio"')
sg.write_text(s)
print("patched java gradle files")
