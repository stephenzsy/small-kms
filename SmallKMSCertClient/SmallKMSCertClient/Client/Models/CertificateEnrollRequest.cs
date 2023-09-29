// <auto-generated/>
using Microsoft.Kiota.Abstractions.Serialization;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System;
namespace SmallKms.Client.Models {
    public class CertificateEnrollRequest : IAdditionalDataHolder, IParsable {
        /// <summary>Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.</summary>
        public IDictionary<string, object> AdditionalData { get; set; }
        /// <summary>The issuer property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public CertificateIssuerParameters? Issuer { get; set; }
#nullable restore
#else
        public CertificateIssuerParameters Issuer { get; set; }
#endif
        /// <summary>The issueToUser property</summary>
        public bool? IssueToUser { get; set; }
        /// <summary>ID of the policy to use for certificate enrollment</summary>
        public Guid? PolicyId { get; set; }
        /// <summary>Property bag of JSON Web Key (RFC 7517) with additional fields</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public JwkProperties? PublicKey { get; set; }
#nullable restore
#else
        public JwkProperties PublicKey { get; set; }
#endif
        /// <summary>The renew property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public CertificateRenewalParameters? Renew { get; set; }
#nullable restore
#else
        public CertificateRenewalParameters Renew { get; set; }
#endif
        /// <summary>The targetFqdn property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? TargetFqdn { get; set; }
#nullable restore
#else
        public string TargetFqdn { get; set; }
#endif
        /// <summary>The usage property</summary>
        public CertificateUsage? Usage { get; set; }
        /// <summary>The validity_months property</summary>
        public int? ValidityMonths { get; set; }
        /// <summary>
        /// Instantiates a new CertificateEnrollRequest and sets the default values.
        /// </summary>
        public CertificateEnrollRequest() {
            AdditionalData = new Dictionary<string, object>();
        }
        /// <summary>
        /// Creates a new instance of the appropriate class based on discriminator value
        /// </summary>
        /// <param name="parseNode">The parse node to use to read the discriminator value and create the object</param>
        public static CertificateEnrollRequest CreateFromDiscriminatorValue(IParseNode parseNode) {
            _ = parseNode ?? throw new ArgumentNullException(nameof(parseNode));
            return new CertificateEnrollRequest();
        }
        /// <summary>
        /// The deserialization information for the current model
        /// </summary>
        public IDictionary<string, Action<IParseNode>> GetFieldDeserializers() {
            return new Dictionary<string, Action<IParseNode>> {
                {"issuer", n => { Issuer = n.GetObjectValue<CertificateIssuerParameters>(CertificateIssuerParameters.CreateFromDiscriminatorValue); } },
                {"issueToUser", n => { IssueToUser = n.GetBoolValue(); } },
                {"policyId", n => { PolicyId = n.GetGuidValue(); } },
                {"publicKey", n => { PublicKey = n.GetObjectValue<JwkProperties>(JwkProperties.CreateFromDiscriminatorValue); } },
                {"renew", n => { Renew = n.GetObjectValue<CertificateRenewalParameters>(CertificateRenewalParameters.CreateFromDiscriminatorValue); } },
                {"targetFqdn", n => { TargetFqdn = n.GetStringValue(); } },
                {"usage", n => { Usage = n.GetEnumValue<CertificateUsage>(); } },
                {"validity_months", n => { ValidityMonths = n.GetIntValue(); } },
            };
        }
        /// <summary>
        /// Serializes information the current object
        /// </summary>
        /// <param name="writer">Serialization writer to use to serialize this model</param>
        public void Serialize(ISerializationWriter writer) {
            _ = writer ?? throw new ArgumentNullException(nameof(writer));
            writer.WriteObjectValue<CertificateIssuerParameters>("issuer", Issuer);
            writer.WriteBoolValue("issueToUser", IssueToUser);
            writer.WriteGuidValue("policyId", PolicyId);
            writer.WriteObjectValue<JwkProperties>("publicKey", PublicKey);
            writer.WriteObjectValue<CertificateRenewalParameters>("renew", Renew);
            writer.WriteStringValue("targetFqdn", TargetFqdn);
            writer.WriteEnumValue<CertificateUsage>("usage", Usage);
            writer.WriteIntValue("validity_months", ValidityMonths);
            writer.WriteAdditionalData(AdditionalData);
        }
    }
}
