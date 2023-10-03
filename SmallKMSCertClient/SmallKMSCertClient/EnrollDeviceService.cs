using Microsoft.Kiota.Http.HttpClientLibrary;
using SmallKms.Client.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace SmallKMSCertClient
{

	internal class EnrollDeviceService
	{
		private readonly AdminAuthProvider authProvider;

		public EnrollDeviceService(AdminAuthProvider authProvider)
		{
			this.authProvider = authProvider;
		}

		public async Task StartEnrollment(Guid groupId, Guid templateId, CertificateEnrollmentRequestDeviceLinkedServicePrincipal req)
		{
			await authProvider.EnsureLoggedIn();
			var adapter = new HttpClientRequestAdapter(authProvider);
			adapter.BaseUrl = "http://localhost:9001";
			var client = new SmallKms.Client.SmallKmsClient(adapter);
			try
			{
				await client.V2.Group[groupId.ToString()].CertificateTemplates[templateId.ToString()].Enroll.PostAsync(req);
			}
			catch (CertificateEnrollmentReceipt400Error e)
			{
				Console.WriteLine(e.AdditionalData["error"]);
			}
		}
	}
}
