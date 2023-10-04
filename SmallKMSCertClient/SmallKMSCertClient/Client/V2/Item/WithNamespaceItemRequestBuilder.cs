// <auto-generated/>
using Microsoft.Kiota.Abstractions;
using Microsoft.Kiota.Cli.Commons.IO;
using Microsoft.Kiota.Cli.Commons;
using SmallKms.Client.V2.Item.CertificateTemplates;
using SmallKms.Client.V2.Item.Certificates;
using SmallKms.Client.V2.Item.LinkServicePrincipal;
using System.Collections.Generic;
using System.CommandLine;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System;
namespace SmallKms.Client.V2.Item {
    /// <summary>
    /// Builds and executes requests for operations under \v2\{namespaceId}
    /// </summary>
    public class WithNamespaceItemRequestBuilder : BaseCliRequestBuilder {
        /// <summary>
        /// The certificates property
        /// </summary>
        public Command BuildCertificatesNavCommand() {
            var command = new Command("certificates");
            command.Description = "The certificates property";
            var builder = new CertificatesRequestBuilder(PathParameters);
            var nonExecCommands = new List<Command>();
            var cmds = builder.BuildCommand();
            nonExecCommands.AddRange(cmds.Item2);
            foreach (var cmd in nonExecCommands.OrderBy(static c => c.Name, StringComparer.Ordinal))
            {
                command.AddCommand(cmd);
            }
            return command;
        }
        /// <summary>
        /// The certificateTemplates property
        /// </summary>
        public Command BuildCertificateTemplatesNavCommand() {
            var command = new Command("certificate-templates");
            command.Description = "The certificateTemplates property";
            var builder = new CertificateTemplatesRequestBuilder(PathParameters);
            var execCommands = new List<Command>();
            var nonExecCommands = new List<Command>();
            execCommands.Add(builder.BuildListCommand());
            var cmds = builder.BuildCommand();
            execCommands.AddRange(cmds.Item1);
            nonExecCommands.AddRange(cmds.Item2);
            foreach (var cmd in execCommands)
            {
                command.AddCommand(cmd);
            }
            foreach (var cmd in nonExecCommands.OrderBy(static c => c.Name, StringComparer.Ordinal))
            {
                command.AddCommand(cmd);
            }
            return command;
        }
        /// <summary>
        /// The linkServicePrincipal property
        /// </summary>
        public Command BuildLinkServicePrincipalNavCommand() {
            var command = new Command("link-service-principal");
            command.Description = "The linkServicePrincipal property";
            var builder = new LinkServicePrincipalRequestBuilder(PathParameters);
            var execCommands = new List<Command>();
            execCommands.Add(builder.BuildGetCommand());
            execCommands.Add(builder.BuildPostCommand());
            foreach (var cmd in execCommands)
            {
                command.AddCommand(cmd);
            }
            return command;
        }
        /// <summary>
        /// Instantiates a new WithNamespaceItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="pathParameters">Path parameters for the request</param>
        public WithNamespaceItemRequestBuilder(Dictionary<string, object> pathParameters) : base("{+baseurl}/v2/{namespaceId}", pathParameters) {
        }
        /// <summary>
        /// Instantiates a new WithNamespaceItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        public WithNamespaceItemRequestBuilder(string rawUrl) : base("{+baseurl}/v2/{namespaceId}", rawUrl) {
        }
    }
}