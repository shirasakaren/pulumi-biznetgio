import * as pulumi from "@pulumi/pulumi";
import * as biznetgio from "@shirasakaren/biznetgio";

// token: `pulumi config set biznetgio:apiToken <token> --secret` atau env BIZNETGIO_API_KEY
const apiToken = biznetgio.config.apiToken;
if (!apiToken) {
    throw new Error("set biznetgio:apiToken di config atau BIZNETGIO_API_KEY di env");
}

const products = biznetgio.neoliteProductsOutput();
const osList = biznetgio.neoliteOsListOutput({ productId: products.products[0].productId });

const keypair = new biznetgio.NeoliteKeypair("demo-keypair", {
    name: "neo-lite-key",
});

const vm = new biznetgio.NeoliteVm("demo-vm", {
    vmName: "neo-lite-1",
    productId: products.products[0].productId,
    selectOs: osList.oss[0].name,
    keypairId: keypair.keypairId,
    cycle: "m",
    sshAndConsoleUser: "admin",
    consolePassword: "s3cretP4ssw0rd",
});

export const vmStatus = vm.status;
