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

		public async Task StartEnrollment()
		{
			await authProvider.EnsureLoggedIn();
		}
	}
}
