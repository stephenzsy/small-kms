using Azure.Core;
using Microsoft.Extensions.Configuration;
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace SmallKMSCertClient
{
	internal static class ConfigUtils
	{
		public static string MustGet(IConfiguration configProvider, string key)
		{
			string? v = configProvider.GetValue<string>(key);
			if (string.IsNullOrEmpty(v))
			{
				throw new Exception($"Missing configuration value for {key}");
			}
			return v;
		}

		public static void StoreConfiguration(string key, string value)
		{
			var filePath = Path.Combine(AppContext.BaseDirectory, "appsettings.json");
			Console.WriteLine($"Configuration file path: {filePath}");
			Dictionary<string, string>? jsonObj;
			if (File.Exists(filePath))
			{
				string json = File.ReadAllText(filePath);
				jsonObj = JsonSerializer.Deserialize<Dictionary<string, string>>(json);
			}
			else
			{
				jsonObj = new Dictionary<string, string>();
			}
			if (jsonObj != null)
			{
				jsonObj[key] = value;
			}
			{
				string json = JsonSerializer.Serialize(jsonObj);
				File.WriteAllText(filePath, json);
			}
		}

	}
}
