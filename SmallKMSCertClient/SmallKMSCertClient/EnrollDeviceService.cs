using Microsoft.Kiota.Http.HttpClientLibrary;
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

		public async Task StartEnrollment(Guid groupId, Guid templateId)
		{
			await authProvider.EnsureLoggedIn();
			var adapter = new HttpClientRequestAdapter(authProvider);
			adapter.BaseUrl = "http://localhost:9001";
			var client = new SmallKms.Client.SmallKmsClient(adapter);
			var body = new SmallKms.Client.Models.CertificateEnrollmentRequestDeviceLinkedServicePrincipal();
			await client.V2.Group[groupId.ToString()].CertificateTemplates[templateId.ToString()].Enroll.PostAsync(body);
		}
	}
}
