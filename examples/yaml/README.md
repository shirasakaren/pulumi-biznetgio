# YAML Example

Kitchen-sink Pulumi program in YAML: invokes the `neoliteProducts` and
`neoliteOsList` functions, then creates a `NeoliteKeypair` and a `NeoliteVm`
that references it, and exports the VM status.

Set the API token before running:

```bash
pulumi config set biznetgio:apiToken <token> --secret
# atau: export BIZNETGIO_API_KEY=<token>
pulumi up
```
