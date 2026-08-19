using System;
using System.Collections.Generic;
using Pulumi;
using Shirasakaren.Biznetgio;

return await Deployment.RunAsync(() =>
{
    // token: `pulumi config set biznetgio:apiToken <token> --secret` atau env BIZNETGIO_API_KEY
    var apiToken = Config.ApiToken;
    if (string.IsNullOrEmpty(apiToken))
    {
        throw new Exception("set biznetgio:apiToken di config atau BIZNETGIO_API_KEY di env");
    }

    var products = NeoliteProducts.Invoke();
    var osList = NeoliteOsList.Invoke(new NeoliteOsListInvokeArgs
    {
        ProductId = products.Apply(p => p.Products[0].ProductId),
    });

    var keypair = new NeoliteKeypair("demo-keypair", new NeoliteKeypairArgs
    {
        Name = "neo-lite-key",
    });

    var vm = new NeoliteVm("demo-vm", new NeoliteVmArgs
    {
        VmName = "neo-lite-1",
        ProductId = products.Apply(p => p.Products[0].ProductId),
        SelectOs = osList.Apply(o => o.Oss[0].Name),
        KeypairId = keypair.KeypairId,
        Cycle = "m",
        SshAndConsoleUser = "admin",
        ConsolePassword = "s3cretP4ssw0rd",
    });

    return new Dictionary<string, object?>
    {
        ["vmStatus"] = vm.Status,
    };
});
