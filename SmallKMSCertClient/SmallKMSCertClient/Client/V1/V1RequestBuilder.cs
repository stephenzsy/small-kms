// <auto-generated/>
using Microsoft.Kiota.Abstractions;
using SmallKms.Client.V1.Diagnostics;
using SmallKms.Client.V1.Item;
using SmallKms.Client.V1.My;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System;
namespace SmallKms.Client.V1 {
    /// <summary>
    /// Builds and executes requests for operations under \v1
    /// </summary>
    public class V1RequestBuilder : BaseRequestBuilder {
        /// <summary>The diagnostics property</summary>
        public DiagnosticsRequestBuilder Diagnostics { get =>
            new DiagnosticsRequestBuilder(PathParameters, RequestAdapter);
        }
        /// <summary>The my property</summary>
        public MyRequestBuilder My { get =>
            new MyRequestBuilder(PathParameters, RequestAdapter);
        }
        /// <summary>Gets an item from the SmallKms.Client.v1.item collection</summary>
        /// <param name="position">Unique identifier of the item</param>
        public WithNamespaceItemRequestBuilder this[string position] { get {
            var urlTplParams = new Dictionary<string, object>(PathParameters);
            urlTplParams.Add("namespaceId", position);
            return new WithNamespaceItemRequestBuilder(urlTplParams, RequestAdapter);
        } }
        /// <summary>
        /// Instantiates a new V1RequestBuilder and sets the default values.
        /// </summary>
        /// <param name="pathParameters">Path parameters for the request</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public V1RequestBuilder(Dictionary<string, object> pathParameters, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/v1", pathParameters) {
        }
        /// <summary>
        /// Instantiates a new V1RequestBuilder and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public V1RequestBuilder(string rawUrl, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/v1", rawUrl) {
        }
    }
}
