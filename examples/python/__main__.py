import pulumi
import pulumi_biznetgio

# token: `pulumi config set biznetgio:apiToken <token> --secret` atau env BIZNETGIO_API_KEY
api_token = pulumi_biznetgio.config.api_token
if not api_token:
    raise RuntimeError("set biznetgio:apiToken di config atau BIZNETGIO_API_KEY di env")

products = pulumi_biznetgio.neolite_products()
os_list = pulumi_biznetgio.neolite_os_list_output(
    product_id=products.products[0].product_id,
)

keypair = pulumi_biznetgio.NeoliteKeypair(
    "demo-keypair",
    name="neo-lite-key",
)

vm = pulumi_biznetgio.NeoliteVm(
    "demo-vm",
    vm_name="neo-lite-1",
    product_id=products.products[0].product_id,
    select_os=os_list.oss[0].name,
    keypair_id=keypair.keypair_id,
    cycle="m",
    ssh_and_console_user="admin",
    console_password="s3cretP4ssw0rd",
)

pulumi.export("vmStatus", vm.status)
