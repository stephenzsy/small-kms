// <auto-generated/>
using Microsoft.Kiota.Abstractions.Serialization;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System;
namespace SmallKms.Client.Models {
    public class CertificateInfo : IAdditionalDataHolder, IParsable {
        /// <summary>Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.</summary>
        public IDictionary<string, object> AdditionalData { get; set; }
        /// <summary>Common name</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? CommonName { get; set; }
#nullable restore
#else
        public string CommonName { get; set; }
#endif
        /// <summary>The issuerCertificate property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public SmallKms.Client.Models.Ref? IssuerCertificate { get; set; }
#nullable restore
#else
        public SmallKms.Client.Models.Ref IssuerCertificate { get; set; }
#endif
        /// <summary>Property bag of JSON Web Key (RFC 7517) with additional fields</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public JwkProperties? Jwk { get; set; }
#nullable restore
#else
        public JwkProperties Jwk { get; set; }
#endif
        /// <summary>Expiration date of the certificate</summary>
        public DateTimeOffset? NotAfter { get; set; }
        /// <summary>Expiration date of the certificate</summary>
        public DateTimeOffset? NotBefore { get; set; }
        /// <summary>The pem property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? Pem { get; set; }
#nullable restore
#else
        public string Pem { get; set; }
#endif
        /// <summary>The ref property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public RefWithMetadata? Ref { get; set; }
#nullable restore
#else
        public RefWithMetadata Ref { get; set; }
#endif
        /// <summary>The subject property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? Subject { get; set; }
#nullable restore
#else
        public string Subject { get; set; }
#endif
        /// <summary>The subjectAlternativeNames property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public CertificateSubjectAlternativeNames? SubjectAlternativeNames { get; set; }
#nullable restore
#else
        public CertificateSubjectAlternativeNames SubjectAlternativeNames { get; set; }
#endif
        /// <summary>The template property</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public SmallKms.Client.Models.Ref? Template { get; set; }
#nullable restore
#else
        public SmallKms.Client.Models.Ref Template { get; set; }
#endif
        /// <summary>The usage property</summary>
        public CertificateUsage? Usage { get; set; }
        /// <summary>
        /// Instantiates a new CertificateInfo and sets the default values.
        /// </summary>
        public CertificateInfo() {
            AdditionalData = new Dictionary<string, object>();
        }
        /// <summary>
        /// Creates a new instance of the appropriate class based on discriminator value
        /// </summary>
        /// <param name="parseNode">The parse node to use to read the discriminator value and create the object</param>
        public static CertificateInfo CreateFromDiscriminatorValue(IParseNode parseNode) {
            _ = parseNode ?? throw new ArgumentNullException(nameof(parseNode));
            return new CertificateInfo();
        }
        /// <summary>
        /// The deserialization information for the current model
        /// </summary>
        public IDictionary<string, Action<IParseNode>> GetFieldDeserializers() {
            return new Dictionary<string, Action<IParseNode>> {
                {"commonName", n => { CommonName = n.GetStringValue(); } },
                {"issuerCertificate", n => { IssuerCertificate = n.GetObjectValue<SmallKms.Client.Models.Ref>(SmallKms.Client.Models.Ref.CreateFromDiscriminatorValue); } },
                {"jwk", n => { Jwk = n.GetObjectValue<JwkProperties>(JwkProperties.CreateFromDiscriminatorValue); } },
                {"notAfter", n => { NotAfter = n.GetDateTimeOffsetValue(); } },
                {"notBefore", n => { NotBefore = n.GetDateTimeOffsetValue(); } },
                {"pem", n => { Pem = n.GetStringValue(); } },
                {"ref", n => { Ref = n.GetObjectValue<RefWithMetadata>(RefWithMetadata.CreateFromDiscriminatorValue); } },
                {"subject", n => { Subject = n.GetStringValue(); } },
                {"subjectAlternativeNames", n => { SubjectAlternativeNames = n.GetObjectValue<CertificateSubjectAlternativeNames>(CertificateSubjectAlternativeNames.CreateFromDiscriminatorValue); } },
                {"template", n => { Template = n.GetObjectValue<SmallKms.Client.Models.Ref>(SmallKms.Client.Models.Ref.CreateFromDiscriminatorValue); } },
                {"usage", n => { Usage = n.GetEnumValue<CertificateUsage>(); } },
            };
        }
        /// <summary>
        /// Serializes information the current object
        /// </summary>
        /// <param name="writer">Serialization writer to use to serialize this model</param>
        public void Serialize(ISerializationWriter writer) {
            _ = writer ?? throw new ArgumentNullException(nameof(writer));
            writer.WriteStringValue("commonName", CommonName);
            writer.WriteObjectValue<SmallKms.Client.Models.Ref>("issuerCertificate", IssuerCertificate);
            writer.WriteObjectValue<JwkProperties>("jwk", Jwk);
            writer.WriteDateTimeOffsetValue("notAfter", NotAfter);
            writer.WriteDateTimeOffsetValue("notBefore", NotBefore);
            writer.WriteStringValue("pem", Pem);
            writer.WriteObjectValue<RefWithMetadata>("ref", Ref);
            writer.WriteStringValue("subject", Subject);
            writer.WriteObjectValue<CertificateSubjectAlternativeNames>("subjectAlternativeNames", SubjectAlternativeNames);
            writer.WriteObjectValue<SmallKms.Client.Models.Ref>("template", Template);
            writer.WriteEnumValue<CertificateUsage>("usage", Usage);
            writer.WriteAdditionalData(AdditionalData);
        }
    }
}
