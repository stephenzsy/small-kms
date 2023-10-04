using Azure.Core;
using Microsoft.Extensions.Configuration;
using Microsoft.Kiota.Serialization.Json;
using SmallKms.Client.Models;
using System;
using System.Collections.Generic;
using System.CommandLine;
using System.Configuration;
using System.Linq;
using System.Runtime.Versioning;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace SmallKMSCertClient
{
	public record JwtHeader
	{
		[JsonPropertyName("typ")]
		public string Type { get; init; } = "JWT";

		[JsonPropertyName("alg")]
		public required string Algorithm { get; init; }
	}

	[SupportedOSPlatform("windows")]
	internal class EnrollCertificateCommand : Command
	{


		public EnrollCertificateCommand() : base("sign-enroll-receipt", "sign enrollment receipt")
		{
			var receiptArg = new Argument<FileInfo>("receipt", "file path of enrollment receipt");
			var signedArg = new Argument<FileInfo>("signed-out", "file path of enrollment receipt");

			AddArgument(receiptArg);
			AddArgument(signedArg);
			this.SetHandler((receipt, outfile) =>
			{
				if (receipt != null)
				{
					var jsonContent = File.ReadAllText(receipt.FullName);
					using var jsonDocument = JsonDocument.Parse(File.ReadAllText(receipt.FullName));
					var parserNode = new JsonParseNode(jsonDocument.RootElement);
					var receiptObj = parserNode.GetObjectValue(CertificateEnrollmentReceipt.CreateFromDiscriminatorValue);

					var header = new JwtHeader()
					{
						Algorithm = "RS256"
					};
					var headerEncoded = ConvertUtils.Base64StdToUrl(Convert.ToBase64String(Encoding.UTF8.GetBytes(JsonSerializer.Serialize(header))));

					var keyId = receiptObj.Ref?.Id.ToString();
					if (string.IsNullOrWhiteSpace(keyId))
					{
						Console.Error.Write("Key ID is missing");
						return;
					}
					var cngKey = CngKey.Exists(keyId)
					? CngKey.Open(keyId)
					: CngKey.Create(CngAlgorithm.Rsa, receiptObj.Ref?.Id.ToString(), new CngKeyCreationParameters
					{
						ExportPolicy = CngExportPolicies.None,
						Parameters = { new CngProperty("Length", BitConverter.GetBytes(2048), CngPropertyOptions.None) }
					});
					var rsaKey = new RSACng(cngKey);
					var sBytes = rsaKey.SignData(
						Encoding.UTF8.GetBytes(headerEncoded + "." + receiptObj.JwtClaims),
						HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1);
					var signatureEncoded = ConvertUtils.Base64StdToUrl(Convert.ToBase64String(sBytes));
					var rsaPem = rsaKey.ExportRSAPublicKeyPem();

					var finalize = new CertificateEnrollmentReplyFinalize
					{
						JwtHeader = headerEncoded,
						JwtSignature = signatureEncoded,
						PublicKeyPem = rsaPem
					};

					var writer = new JsonSerializationWriter();
					writer.WriteObjectValue(string.Empty, finalize);
					var readStream = new StreamReader(writer.GetSerializedContent());
					var serializedOutput = readStream.ReadToEnd();
					File.WriteAllText(outfile.FullName, serializedOutput);
				}
			}, receiptArg, signedArg);
		}
	}
}
