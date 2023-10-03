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
			Option groupIdOption = new Option<Guid>("--group-id", "Group ID");
			Option templateIdOption = new Option<Guid>("--template-id", "Template ID");
			Option appIdOption = new Option<Guid>("--app-id", "App client ID");
			Option linkIdOption = new Option<Guid>("--link-id", "Device link ID");
			Option deviceNamespaceIdOption = new Option<Guid>("--device-namespace-id", "Device object ID");
			Option servicePrincipalIdOption = new Option<Guid>("--service-principal-id", "Service principal object ID");

			var cmd = new Command("enroll-device", "Enroll this device to Small KMS provisioning a certificate to be used for authenticate as Microsoft Entra ID service principal");
			cmd.AddOption(groupIdOption);
			cmd.AddOption(templateIdOption);
			cmd.AddOption(appIdOption);
			cmd.AddOption(linkIdOption);
			cmd.AddOption(deviceNamespaceIdOption);
			cmd.AddOption(servicePrincipalIdOption);

			cmd.SetHandler((Guid groupId, Guid templateId, Guid appId, Guid linkId, Guid deviceNamespaceId, Guid servicePrincipalId) =>
			serviceProvider.GetService<EnrollDeviceService>()?.StartEnrollment(groupId, templateId, new SmallKms.Client.Models.CertificateEnrollmentRequestDeviceLinkedServicePrincipal
			{
				AppId = appId,
				LinkId = linkId,
				Type = SmallKms.Client.Models.CertificateEnrollmentTargetType.DeviceLinkedServicePrincipal,
				DeviceNamespaceId = deviceNamespaceId,
				ServicePrincipalId= servicePrincipalId,
			}) ?? Task.CompletedTask,
				(System.CommandLine.Binding.IValueDescriptor<Guid>)groupIdOption,
				(System.CommandLine.Binding.IValueDescriptor<Guid>)templateIdOption,
				(System.CommandLine.Binding.IValueDescriptor<Guid>)appIdOption,
				(System.CommandLine.Binding.IValueDescriptor<Guid>)linkIdOption,
				(System.CommandLine.Binding.IValueDescriptor<Guid>)deviceNamespaceIdOption,
				(System.CommandLine.Binding.IValueDescriptor<Guid>)servicePrincipalIdOption);
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
