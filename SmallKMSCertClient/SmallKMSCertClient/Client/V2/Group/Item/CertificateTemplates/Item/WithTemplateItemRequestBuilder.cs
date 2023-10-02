// <auto-generated/>
using Microsoft.Kiota.Abstractions;
using SmallKms.Client.V2.Group.Item.CertificateTemplates.Item.Enroll;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System;
namespace SmallKms.Client.V2.Group.Item.CertificateTemplates.Item {
    /// <summary>
    /// Builds and executes requests for operations under \v2\group\{namespaceId}\certificate-templates\{templateId}
    /// </summary>
    public class WithTemplateItemRequestBuilder : BaseRequestBuilder {
        /// <summary>The enroll property</summary>
        public EnrollRequestBuilder Enroll { get =>
            new EnrollRequestBuilder(PathParameters, RequestAdapter);
        }
        /// <summary>
        /// Instantiates a new WithTemplateItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="pathParameters">Path parameters for the request</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public WithTemplateItemRequestBuilder(Dictionary<string, object> pathParameters, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/v2/group/{namespaceId}/certificate-templates/{templateId}", pathParameters) {
        }
        /// <summary>
        /// Instantiates a new WithTemplateItemRequestBuilder and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public WithTemplateItemRequestBuilder(string rawUrl, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/v2/group/{namespaceId}/certificate-templates/{templateId}", rawUrl) {
        }
    }
}
