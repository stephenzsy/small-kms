using Azure.Core;
using Microsoft.Extensions.Configuration;
using Microsoft.Identity.Client;
using Microsoft.Identity.Client.Broker;
using Microsoft.Identity.Client.Extensions.Msal;
using Microsoft.Identity.Client.NativeInterop;
using Microsoft.Kiota.Http.HttpClientLibrary;
using Microsoft.VisualBasic;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading.Tasks;

namespace SmallKMSCertClient
{
	class MsalTokenCredential : TokenCredential
	{
		private readonly IPublicClientApplication app;
		private readonly string? accountIdentifier;
		private readonly IEnumerable<string> loginScopes;
		private IAccount? account;
		private readonly IConfiguration configProvider;
		public MsalTokenCredential(IConfiguration config)
		{
			this.configProvider = config;
			var clientId = ConfigUtils.MustGet(config, "AZURE_CLIENT_ID");
			this.loginScopes = new string[] { ConfigUtils.MustGet(config, "AZURE_LOGIN_SCOPE") };


			this.app = PublicClientApplicationBuilder
				.Create(clientId)
				.WithTenantId(ConfigUtils.MustGet(config, "AZURE_TENANT_ID"))
				.WithDefaultRedirectUri()
				.Build();

			this.accountIdentifier = config.GetValue<string>("AZURE_ACCOUNT_IDENTIFIER");
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

		public async Task<bool> IsLoggedIn()
		{
			try
			{
				return (await acquireTokenSilently()) != null;
			}
			catch
			{ // swallow error
			}
			return false;
		}

		private async Task<AuthenticationResult> acquireTokenSilently()
		{
			if (account == null)
			{
				if (string.IsNullOrEmpty(accountIdentifier))
				{
					throw new Exception("No account identifier provided");
				}
				account = await app.GetAccountAsync(accountIdentifier);
			}
			if (account != null)
			{
				var result = await app.AcquireTokenSilent(loginScopes, account).ExecuteAsync();
				Console.WriteLine(result.AccessToken);
				return result;
			}
			throw new Exception("No account loaded");
		}

		public override AccessToken GetToken(TokenRequestContext requestContext, CancellationToken cancellationToken)
		{
			return GetTokenAsync(requestContext, cancellationToken).GetAwaiter().GetResult();
		}

		public override async ValueTask<AccessToken> GetTokenAsync(TokenRequestContext requestContext, CancellationToken cancellationToken)
		{
			var authResult = await acquireTokenSilently();

			return new AccessToken(authResult.AccessToken, authResult.ExpiresOn);
		}

		internal async Task Login(bool useDeviceCode = false)
		{
			AuthenticationResult result = await (useDeviceCode
				? app.AcquireTokenWithDeviceCode(loginScopes, deviceCodeResult =>
				{
					Console.WriteLine(deviceCodeResult.Message);
					return Task.FromResult(deviceCodeResult);
				}).ExecuteAsync()
				: app.AcquireTokenInteractive(loginScopes).ExecuteAsync());
			account = result.Account;
			ConfigUtils.StoreConfiguration("AZURE_ACCOUNT_IDENTIFIER", account.HomeAccountId.Identifier);
		}
	}
}
