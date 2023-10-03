using Azure.Core;
using Azure.Identity;
using Microsoft.Extensions.Configuration;
using Microsoft.Kiota.Abstractions;
using Microsoft.Kiota.Abstractions.Authentication;
using Microsoft.Kiota.Authentication.Azure;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace SmallKMSCertClient
{
	internal enum AuthPoviderCredentialMode
	{
		interactive,
		deviceCode
	}


	internal class AdminAuthProvider : IAuthenticationProvider
	{
		private readonly MsalTokenCredential tokenCredential;
		private readonly AzureIdentityAuthenticationProvider innerProvider;
		public AdminAuthProvider(IConfiguration config)
		{
			this.tokenCredential = new MsalTokenCredential(config);
			this.innerProvider = new AzureIdentityAuthenticationProvider(this.tokenCredential);
		}

		public async Task EnsureLoggedIn()
		{
			if (!await this.tokenCredential.IsLoggedIn())
			{
				throw new Exception("Not logged in");
			}
		}

		public async Task AuthenticateRequestAsync(RequestInformation request, Dictionary<string, object>? additionalAuthenticationContext = null, CancellationToken cancellationToken = default)
		{
			var token = await this.tokenCredential.GetTokenAsync(new TokenRequestContext { }, cancellationToken);
			request.Headers["Authorization"] = new string[]{ $"Bearer {token.Token}" };
		}

		internal Task Login(bool useDeviceCode = false)
		{
			return this.tokenCredential.Login(useDeviceCode);
		}
	}
}
