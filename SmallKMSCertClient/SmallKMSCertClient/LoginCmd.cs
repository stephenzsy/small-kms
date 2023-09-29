using System.CommandLine;
using System.Configuration;

namespace SmallKMSCertClient
{
	internal static class LoginCmd
	{
		public static Command BuildCommand()
		{
			var configCmd = new Command("login", "Login with Microsoft Entra ID");

			var showCmd = new Command("show", "Show app configurations");
			showCmd.SetHandler(showAppConfigurations);
			configCmd.AddCommand(showCmd);


			return configCmd;
		}

		private static void showAppConfigurations()
		{
			Console.WriteLine("Small KMS certificate utilities configurations");
			Console.WriteLine("==============================================");
			ConfigurationManager.AppSettings.AllKeys.ToList()
				.ForEach(key => Console.WriteLine($"{key}: {ConfigurationManager.AppSettings[key]}"));
		}
	}
}
