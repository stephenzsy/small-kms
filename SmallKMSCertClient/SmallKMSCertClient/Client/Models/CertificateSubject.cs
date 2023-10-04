// <auto-generated/>
using Microsoft.Kiota.Abstractions.Serialization;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System;
namespace SmallKms.Client.Models {
    public class CertificateSubject : IAdditionalDataHolder, IParsable {
        /// <summary>Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.</summary>
        public IDictionary<string, object> AdditionalData { get; set; }
        /// <summary>Country or region</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? C { get; set; }
#nullable restore
#else
        public string C { get; set; }
#endif
        /// <summary>Common name</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? Cn { get; set; }
#nullable restore
#else
        public string Cn { get; set; }
#endif
        /// <summary>Organization</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? O { get; set; }
#nullable restore
#else
        public string O { get; set; }
#endif
        /// <summary>Organizational unit</summary>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public string? Ou { get; set; }
#nullable restore
#else
        public string Ou { get; set; }
#endif
        /// <summary>
        /// Instantiates a new CertificateSubject and sets the default values.
        /// </summary>
        public CertificateSubject() {
            AdditionalData = new Dictionary<string, object>();
        }
        /// <summary>
        /// Creates a new instance of the appropriate class based on discriminator value
        /// </summary>
        /// <param name="parseNode">The parse node to use to read the discriminator value and create the object</param>
        public static CertificateSubject CreateFromDiscriminatorValue(IParseNode parseNode) {
            _ = parseNode ?? throw new ArgumentNullException(nameof(parseNode));
            return new CertificateSubject();
        }
        /// <summary>
        /// The deserialization information for the current model
        /// </summary>
        public IDictionary<string, Action<IParseNode>> GetFieldDeserializers() {
            return new Dictionary<string, Action<IParseNode>> {
                {"c", n => { C = n.GetStringValue(); } },
                {"cn", n => { Cn = n.GetStringValue(); } },
                {"o", n => { O = n.GetStringValue(); } },
                {"ou", n => { Ou = n.GetStringValue(); } },
            };
        }
        /// <summary>
        /// Serializes information the current object
        /// </summary>
        /// <param name="writer">Serialization writer to use to serialize this model</param>
        public void Serialize(ISerializationWriter writer) {
            _ = writer ?? throw new ArgumentNullException(nameof(writer));
            writer.WriteStringValue("c", C);
            writer.WriteStringValue("cn", Cn);
            writer.WriteStringValue("o", O);
            writer.WriteStringValue("ou", Ou);
            writer.WriteAdditionalData(AdditionalData);
        }
    }
}