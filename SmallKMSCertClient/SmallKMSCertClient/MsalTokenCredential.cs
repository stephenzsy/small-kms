using Azure.Core;
using Azure.Identity;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Options;
using Microsoft.Identity.Client;
using Microsoft.Identity.Client.Extensions.Msal;
using Spectre.Console;
using System.CommandLine;

namespace SmallKMSCertClient
{
	class MsalTokenCredential : Azure.Core.TokenCredential
	{
		private readonly IPublicClientApplication app;
		private string? accountIdentifier;

		public override ValueTask<AccessToken> GetTokenAsync(TokenRequestContext requestContext, CancellationToken cancellationToken)
		{
			return this.GetTokenImplAsync(requestContext, cancellationToken);
		}

		public override AccessToken GetToken(TokenRequestContext requestContext, CancellationToken cancellationToken)
		{
			return GetTokenImplAsync(requestContext, cancellationToken).GetAwaiter().GetResult();
		}


		public MsalTokenCredential(DeviceCodeCredentialOptions opts)
		{

			this.app = PublicClientApplicationBuilder
				.Create(opts.ClientId)
				.WithTenantId(opts.TenantId)
				.WithDefaultRedirectUri()
				.Build();

			try
			{
				this.accountIdentifier = File.ReadAllText(Path.Join(AppContext.BaseDirectory, "accountIdentifier.txt")).Trim();
			}
			catch
			{

			}
			this.registerCache();
		}
		private async void registerCache()
		{

			var storageProperties =
				new StorageCreationPropertiesBuilder("tokens.bin", AppContext.BaseDirectory)
				.Build();
			var cacheHelper = await MsalCacheHelper.CreateAsync(storageProperties);
			cacheHelper.RegisterCache(app.UserTokenCache);
		}

		private async ValueTask<AccessToken> GetTokenImplAsync(TokenRequestContext requestContext, CancellationToken cancellationToken)
		{
			var useDeviceCode = string.IsNullOrWhiteSpace(this.accountIdentifier);
			if (!useDeviceCode)
			{
				try
				{
					var account = await this.app.GetAccountAsync(this.accountIdentifier);
					var tokenResultSilent = await this.app.AcquireTokenSilent(requestContext.Scopes, account).ExecuteAsync(cancellationToken);
					return new AccessToken(tokenResultSilent.AccessToken, tokenResultSilent.ExpiresOn);
				}
				catch (MsalUiRequiredException)
				{
					useDeviceCode = true;
				}
			}
			var tokenResult = await this.app.AcquireTokenWithDeviceCode(requestContext.Scopes, deviceCodeResult =>
			{
				Console.WriteLine(deviceCodeResult.Message);
				return Task.FromResult(deviceCodeResult);
			}).ExecuteAsync(cancellationToken);
			this.accountIdentifier = tokenResult.Account.HomeAccountId.Identifier;
			File.WriteAllText(Path.Join(AppContext.BaseDirectory, "accountIdentifier.txt"), this.accountIdentifier);
			return new AccessToken(tokenResult.AccessToken, tokenResult.ExpiresOn);
		}

	}
}