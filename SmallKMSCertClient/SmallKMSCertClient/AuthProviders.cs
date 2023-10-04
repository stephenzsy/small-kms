using Azure.Core;
using Azure.Identity;
using Microsoft.Extensions.Configuration;
using Microsoft.Identity.Client.Extensions.Msal;
using Microsoft.Kiota.Abstractions;
using Microsoft.Kiota.Abstractions.Authentication;
using Microsoft.Kiota.Authentication.Azure;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static System.Formats.Asn1.AsnWriter;

namespace SmallKMSCertClient
{
	internal enum AuthPoviderCredentialMode
	{
		interactive,
		deviceCode
	}


	internal class OverrideHttpsTokenPovider : AzureIdentityAccessTokenProvider, IAccessTokenProvider, IDisposable
	{

		public OverrideHttpsTokenPovider(TokenCredential credential, string[]? allowedHosts = null, ObservabilityOptions? observabilityOptions = null, params string[] scopes) : base(credential, allowedHosts, observabilityOptions, scopes)
		{
		}

		public new Task<string> GetAuthorizationTokenAsync(Uri uri, Dictionary<string, object>? additionalAuthenticationContext = default, CancellationToken cancellationToken = default)
		{
			if (uri.Scheme == "http")
			{
				var uribuider = new UriBuilder(uri) { Scheme = "https" };
				uri = uribuider.Uri;
			}
			return base.GetAuthorizationTokenAsync(uri, additionalAuthenticationContext, cancellationToken);
		}


	}

	internal class OverrideHttpsAuthenticationPovider : BaseBearerTokenAuthenticationProvider
	{

		public OverrideHttpsAuthenticationPovider(TokenCredential credential, string[]? allowedHosts = null, ObservabilityOptions? observabilityOptions = null, params string[] scopes) :
			base(new OverrideHttpsTokenPovider(credential, allowedHosts, observabilityOptions, scopes))
		{
		}

	}
}
