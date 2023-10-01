using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using System.CommandLine;

namespace SmallKMSCertClient
{
	internal class CommandService
	{
		private readonly IServiceProvider serviceProvider;
		private readonly IConfiguration configs;
		public CommandService(IServiceProvider serviceProvider, IConfiguration configs)
		{
			this.serviceProvider = serviceProvider;
			this.configs = configs;
		}

		internal void ConfigureCommand(RootCommand rootCommand)
		{
			rootCommand.AddCommand(buildConfigCmd());
			rootCommand.AddCommand(buildLoginCmd());
			rootCommand.AddCommand(buildEnrollDeviceCmd());
		}

		private Command buildEnrollDeviceCmd()
		{
			var cmd = new Command("enroll-device", "Enroll this device to Small KMS provisioning a certificate to be used for authenticate as Microsoft Entra ID service principal");

			cmd.SetHandler(() =>
			serviceProvider.GetService<EnrollDeviceService>()?.StartEnrollment() ?? Task.CompletedTask);
			return cmd;
		}

		private Command buildConfigCmd()
		{
			Option keyOption = new Option<string>("--key", "Configuration key");
			Option valueOption = new Option<string>("--value", "Configuration value");

			var cmd = new Command("config", "View or edit configurations");

			var setCmd = new Command("set", "Set configuration");
			setCmd.AddOption(keyOption);
			setCmd.AddOption(valueOption);
			setCmd.SetHandler(ConfigUtils.StoreConfiguration, (System.CommandLine.Binding.IValueDescriptor<string>)keyOption, (System.CommandLine.Binding.IValueDescriptor<string>)valueOption);
			cmd.AddCommand(setCmd);

			return cmd;

		}

		private Command buildLoginCmd()
		{
			var cmd = new Command("login", "Login to Azure AD");
			var deviceCodeOption = new Option<bool>("--device-code", "Use device code to login");
			cmd.AddOption(deviceCodeOption);
			cmd.SetHandler((bool useDeviceCode) => serviceProvider.GetService<AdminAuthProvider>()?.Login(useDeviceCode) ?? Task.CompletedTask, deviceCodeOption);

			return cmd;
		}
	}
}
