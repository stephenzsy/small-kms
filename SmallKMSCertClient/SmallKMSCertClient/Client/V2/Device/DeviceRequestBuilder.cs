// <auto-generated/>
using Microsoft.Kiota.Abstractions;
using Microsoft.Kiota.Cli.Commons.IO;
using Microsoft.Kiota.Cli.Commons;
using SmallKms.Client.V2.Device.Item;
using System.Collections.Generic;
using System.CommandLine;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System;
namespace SmallKms.Client.V2.Device {
    /// <summary>
    /// Builds and executes requests for operations under \v2\device
    /// </summary>
    public class DeviceRequestBuilder : BaseCliRequestBuilder {
        /// <summary>
        /// Gets an item from the SmallKms.Client.v2.device.item collection
        /// </summary>
        public Tuple<List<Command>, List<Command>> BuildCommand() {
            var commands = new List<Command>();
            var builder = new WithNamespaceItemRequestBuilder(PathParameters);
            commands.Add(builder.BuildLinkServicePrincipalNavCommand());
            return new(new(0), commands);
        }
        /// <summary>
        /// Instantiates a new DeviceRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="pathParameters">Path parameters for the request</param>
        public DeviceRequestBuilder(Dictionary<string, object> pathParameters) : base("{+baseurl}/v2/device", pathParameters) {
        }
        /// <summary>
        /// Instantiates a new DeviceRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        public DeviceRequestBuilder(string rawUrl) : base("{+baseurl}/v2/device", rawUrl) {
        }
    }
}
