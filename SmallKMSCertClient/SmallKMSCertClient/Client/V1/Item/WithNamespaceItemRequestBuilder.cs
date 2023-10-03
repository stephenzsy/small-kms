// <auto-generated/>
using Microsoft.Kiota.Abstractions;
using Microsoft.Kiota.Cli.Commons.IO;
using Microsoft.Kiota.Cli.Commons;
using SmallKms.Client.V1.Item.Policies;
using SmallKms.Client.V1.Item.Profile;
using System.Collections.Generic;
using System.CommandLine;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System;
namespace SmallKms.Client.V1.Item {
    /// <summary>
    /// Builds and executes requests for operations under \v1\{namespaceId}
    /// </summary>
    public class WithNamespaceItemRequestBuilder : BaseCliRequestBuilder {
        /// <summary>
        /// The policies property
        /// </summary>
        public Command BuildPoliciesNavCommand() {
            var command = new Command("policies");
            command.Description = "The policies property";
            var builder = new PoliciesRequestBuilder(PathParameters);
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
        /// The profile property
        /// </summary>
        public Command BuildProfileNavCommand() {
            var command = new Command("profile");
            command.Description = "The profile property";
            var builder = new ProfileRequestBuilder(PathParameters);
            var execCommands = new List<Command>();
            execCommands.Add(builder.BuildGetCommand());
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
        public WithNamespaceItemRequestBuilder(Dictionary<string, object> pathParameters) : base("{+baseurl}/v1/{namespaceId}", pathParameters) {
        }
        /// <summary>
        /// Instantiates a new WithNamespaceItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        public WithNamespaceItemRequestBuilder(string rawUrl) : base("{+baseurl}/v1/{namespaceId}", rawUrl) {
        }
    }
}
