using Microsoft.Extensions.DependencyInjection;
using System.CommandLine;

namespace SmallKMSCertClient
{
	internal class CommandService
	{
		private readonly IServiceProvider serviceProvider;
		public CommandService(IServiceProvider serviceProvider)
		{
			this.serviceProvider = serviceProvider;
		}

		internal void ConfigureCommand(RootCommand rootCommand)
		{
			rootCommand.AddCommand(buildEnrollDeviceCmd());
		}

		private Command buildEnrollDeviceCmd()
		{
			var cmd = new Command("enroll-device", "Enroll this device to Small KMS provisioning a certificate to be used for authenticate as Microsoft Entra ID service principal");
			cmd.SetHandler(() =>
			serviceProvider.GetService<EnrollDeviceService>()?.StartEnrollment() ?? Task.CompletedTask);
			return cmd;
		}
	}
}
