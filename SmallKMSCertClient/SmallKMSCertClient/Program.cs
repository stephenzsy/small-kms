// See https://aka.ms/new-console-template for more information
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using SmallKMSCertClient;
using System.CommandLine;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Configuration.Json;

HostApplicationBuilder builder = Host.CreateApplicationBuilder();
builder.Configuration
	.AddJsonFile(Path.Join(AppContext.BaseDirectory, "appsettings.json"), optional: true, reloadOnChange: true);
builder.Services
	.AddSingleton<AdminAuthProvider>()
	.AddSingleton<EnrollDeviceService>()
	.AddSingleton<CommandService>();
var host = builder.Build();


var rootCommand = new RootCommand("Small KMS certificate utils");
host.Services.GetService<CommandService>()?.ConfigureCommand(rootCommand);

return await rootCommand.InvokeAsync(args);
