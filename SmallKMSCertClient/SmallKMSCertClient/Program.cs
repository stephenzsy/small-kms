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
using SmallKMSCertClient;
using System.Net;
using System.Runtime.Versioning;

internal class Program
{
	[SupportedOSPlatform("windows")]
	private static async Task<int> Main(string[] args)
	{
		var rootCommand = new SmallKmsClient().BuildRootCommand();
		rootCommand.Description = "Small KMS CLI";

		var builder = new CommandLineBuilder(rootCommand)
			.UseDefaults()
			.UseRequestAdapter(context =>
			{
				var options = new DeviceCodeCredentialOptions
				{
					TenantId = Environment.GetEnvironmentVariable("AZURE_TENANT_ID") ?? "",
					ClientId = Environment.GetEnvironmentVariable("AZURE_CLIENT_ID") ?? "",

					DeviceCodeCallback = (code, cancellation) =>
					{
						Console.WriteLine(code.Message);
						return Task.FromResult(0);
					},
					TokenCachePersistenceOptions = new TokenCachePersistenceOptions
					{
						Name = "tokens.bin",
						UnsafeAllowUnencryptedStorage = true,
					}
				};
				var deviceCodeCredentials = new MsalTokenCredential(options);
				var authProvider = new OverrideHttpsAuthenticationPovider(deviceCodeCredentials,
					scopes: new string[] { Environment.GetEnvironmentVariable("SMALLKMS_LOGIN_SCOPE") ?? "" });
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
		builder.Command.AddCommand(new ViewRecieptCommand());
		builder.Command.AddCommand(new EnrollCertificateCommand());

		return await builder.Build().InvokeAsync(args);
	}
}