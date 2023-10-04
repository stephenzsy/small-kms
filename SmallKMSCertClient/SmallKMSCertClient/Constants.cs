using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace SmallKMSCertClient
{
	internal static class Constants
	{
		public const string UserSecretsId = @"6383a9b0-1c9a-489e-8918-0c84756fc377";
	}

	internal static class ConvertUtils
	{
		public static string Base64UrlToStd(string src)
		{
			var s = src;
			s = s.Replace('-', '+');
			s = s.Replace('_', '/');
			switch (s.Length % 4)
			{
				case 0:
					break;
				case 2:
					s += "==";
					break;
				case 3:
					s += "=";
					break;
				default:
					throw new Exception("Illegal base64url string!");
			}
			return s;
		}

		internal static string Base64StdToUrl(string v)
		{
			var s = v;
			s = s.Replace('+', '-');
			s = s.Replace('/', '_');
			s.Replace("=", "");
			return s;
		}
	}
}
