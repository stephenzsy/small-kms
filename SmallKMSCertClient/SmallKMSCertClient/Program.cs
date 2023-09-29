// See https://aka.ms/new-console-template for more information
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using SmallKMSCertClient;
using System.CommandLine;

HostApplicationBuilder builder = Host.CreateApplicationBuilder();
builder.Services
	.AddSingleton<EnrollDeviceService>()
	.AddSingleton<CommandService>();
var host = builder.Build();


var rootCommand = new RootCommand("Small KMS certificate utils");
host.Services.GetService<CommandService>()?.ConfigureCommand(rootCommand);

return await rootCommand.InvokeAsync(args);
