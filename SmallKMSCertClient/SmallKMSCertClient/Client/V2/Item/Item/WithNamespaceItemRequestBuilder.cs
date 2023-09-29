// <auto-generated/>
using Microsoft.Kiota.Abstractions;
using SmallKms.Client.V2.Item.Item.CertificateTemplates;
using SmallKms.Client.V2.Item.Item.GraphSync;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System;
namespace SmallKms.Client.V2.Item.Item {
    /// <summary>
    /// Builds and executes requests for operations under \v2\{namespaceType}\{namespaceId}
    /// </summary>
    public class WithNamespaceItemRequestBuilder : BaseRequestBuilder {
        /// <summary>The certificateTemplates property</summary>
        public CertificateTemplatesRequestBuilder CertificateTemplates { get =>
            new CertificateTemplatesRequestBuilder(PathParameters, RequestAdapter);
        }
        /// <summary>The graphSync property</summary>
        public GraphSyncRequestBuilder GraphSync { get =>
            new GraphSyncRequestBuilder(PathParameters, RequestAdapter);
        }
        /// <summary>
        /// Instantiates a new WithNamespaceItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="pathParameters">Path parameters for the request</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public WithNamespaceItemRequestBuilder(Dictionary<string, object> pathParameters, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/v2/{namespaceType}/{namespaceId}", pathParameters) {
        }
        /// <summary>
        /// Instantiates a new WithNamespaceItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public WithNamespaceItemRequestBuilder(string rawUrl, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/v2/{namespaceType}/{namespaceId}", rawUrl) {
        }
    }
}
