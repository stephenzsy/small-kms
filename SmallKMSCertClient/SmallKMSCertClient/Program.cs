using System.CommandLine.Builder;
using System.CommandLine.Parsing;
using SmallKms.Client;
using Microsoft.Kiota.Abstractions;
using Microsoft.Kiota.Abstractions.Authentication;
using Microsoft.Kiota.Cli.Commons.Extensions;
using Microsoft.Kiota.Http.HttpClientLibrary;
using Microsoft.Kiota.Serialization.Form;
using Microsoft.Kiota.Serialization.Json;
using Microsoft.Kiota.Serialization.Text;
using Microsoft.Kiota.Authentication.Azure;
using Azure.Identity;

var rootCommand = new SmallKmsClient().BuildRootCommand();
rootCommand.Description = "Small KMS CLI";

var builder = new CommandLineBuilder(rootCommand)
	.UseDefaults()
	.UseRequestAdapter(context =>
	{
		var authProvider = new AzureIdentityAuthenticationProvider(new DefaultAzureCredential(new DefaultAzureCredentialOptions()
		{
			ExcludeInteractiveBrowserCredential = false
		}));
		var adapter = new HttpClientRequestAdapter(authProvider);
		adapter.BaseUrl = "http://localhost:9001";

		// Register default serializers
		ApiClientBuilder.RegisterDefaultSerializer<JsonSerializationWriterFactory>();
		ApiClientBuilder.RegisterDefaultSerializer<TextSerializationWriterFactory>();
		ApiClientBuilder.RegisterDefaultSerializer<FormSerializationWriterFactory>();

		// Register default deserializers
		ApiClientBuilder.RegisterDefaultDeserializer<JsonParseNodeFactory>();
		ApiClientBuilder.RegisterDefaultDeserializer<TextParseNodeFactory>();
		ApiClientBuilder.RegisterDefaultDeserializer<FormParseNodeFactory>();

		return adapter;
	}).RegisterCommonServices();

return await builder.Build().InvokeAsync(args);