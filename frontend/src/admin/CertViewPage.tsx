import { useRequest } from "ahooks";
import { useMemo } from "react";
import { useParams } from "react-router-dom";

import { WellknownId } from "../constants";
import { useCertsApi } from "../utils/useCertsApi";
import { AdminBreadcrumb, BreadcrumbPageMetadata } from "./AdminBreadcrumb";

export const caBreadcrumPages: BreadcrumbPageMetadata[] = [
  { name: "CA", to: "/admin/ca" },
];
/*
function CertDownloadLink({
  download,
  className,
  accept,
  cert,
  children,
  getLinkText,
}: PropsWithChildren<{
  cert: CertificateRef;
  className?: string;
  download: string;
  getLinkText: React.ReactNode;
  accept: Exclude<GetCertificateV1AcceptEnum, "application/json">;
}>) {
  const downloadClient = useAuthedClient(CertsDownloadApi);

  const [blobUrl, setBlobUrl] = useState<string>();
  const { data: certBlob, run: getDownloadUrl } = useRequest(
    () => {
      return downloadClient.getCertificateV1({
        namespaceId: cert.namespaceId,
        id: cert.id,
        accept,
      });
    },
    { manual: true }
  );

  useEffect(() => {
    if (certBlob) {
      const blobUrl = URL.createObjectURL(certBlob);
      setBlobUrl(blobUrl);
      return () => {
        URL.revokeObjectURL(blobUrl);
        setBlobUrl(undefined);
      };
    }
  }, [certBlob]);

  if (!blobUrl) {
    return (
      <button type="button" className={className} onClick={getDownloadUrl}>
        {getLinkText}
      </button>
    );
  }

  return (
    <a href={blobUrl} className={className} download={download}>
      {children}
    </a>
  );
}
*/
export default function CertViewPage() {
  const client = useCertsApi();
  const { certId, namespaceId } = useParams();
  const { data: cert } = useRequest(async () => {
    return client.getCertificateV1({
      namespaceId: namespaceId!,
      id: certId!,
    });
  }, {});

  const breadcrumPages: BreadcrumbPageMetadata[] = useMemo(() => {
    switch (namespaceId) {
      case WellknownId.nsRootCa:
        return [...caBreadcrumPages, { name: "View Certificate", to: "#" }];
    }
    return [{ name: "View Certificate", to: "#" }];
  }, [namespaceId]);

  return (
    <>
      <AdminBreadcrumb pages={breadcrumPages} />
      <section className="divide-y divide-neutral-200 overflow-hidden rounded-lg bg-white shadow">
        <h2 className="px-4 py-5 sm:px-6 text-lg font-semibold">
          Certificate details
        </h2>
        {cert && (
          <>
            <dl className="px-4 py-5 sm:p-6" key={cert.id}>
              <div className="flex gap-x-2">
                <dt className="text-sm font-medium">Common Name</dt>
                <dd>{cert.name}</dd>
              </div>
              <div className="flex gap-x-2">
                <dt className="text-sm font-medium">Expiry</dt>
                <dd>{cert.notAfter.toLocaleString()}</dd>
              </div>
            </dl>
            <div className="px-4 py-5 sm:p-6 flex gap-x-6">
              {/*<CertDownloadLink
                cert={cert}
                download={`${cert.id}.pem`}
                accept={GetCertificateV1AcceptEnum.Accept_Pem}
                className="text-sm font-medium text-indigo-600 hover:text-indigo-500"
                getLinkText="Get download .PEM file link"
              >
                Download .PEM file
              </CertDownloadLink>
              <CertDownloadLink
                cert={cert}
                download={`${cert.id}.crt`}
                accept={GetCertificateV1AcceptEnum.Accept_X509CaCert}
                className="text-sm font-medium text-indigo-600 hover:text-indigo-500"
                getLinkText="Get download .CRT file link"
              >
                Download .CRT file
        </CertDownloadLink>*/}
            </div>
          </>
        )}
      </section>
    </>
  );
}
