using Microsoft.Kiota.Serialization.Form;
using Microsoft.Kiota.Serialization.Json;
using Microsoft.Kiota.Serialization.Text;

using SmallKms.Client.Models;
using System.Buffers.Text;
using System.CommandLine;
using System.Configuration;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace SmallKMSCertClient
{
	internal class ViewRecieptCommand : Command
	{

		public ViewRecieptCommand() : base("view-receipt", "view reciept of request json")
		{
			var fileArg = new Argument<FileInfo>("file", "file path of request json");
			AddArgument(fileArg);
			this.SetHandler((file) =>
			{
				if (file != null)
				{
					var jsonContent = File.ReadAllText(file.FullName);
					using var jsonDocument = JsonDocument.Parse(File.ReadAllText(file.FullName));
					var parserNode = new JsonParseNode(jsonDocument.RootElement);
					var receipt = parserNode.GetObjectValue(CertificateEnrollmentReceipt.CreateFromDiscriminatorValue);

					var claimsJsonStr = System.Text.Encoding.UTF8.GetString(Convert.FromBase64String(ConvertUtils.Base64UrlToStd(receipt.JwtClaims ?? "")));

					var claimsJsonDoc = JsonDocument.Parse(claimsJsonStr);

					Console.WriteLine("JWT Claims:");
					Console.WriteLine("=======================");
					var serializedClaims = JsonSerializer.Serialize(claimsJsonDoc, new JsonSerializerOptions { WriteIndented = true }) ;
					Console.WriteLine(serializedClaims);
					Console.WriteLine("=======================");
					Console.WriteLine("Please verify the names");
				}
			}, fileArg);
		}
	}
}
