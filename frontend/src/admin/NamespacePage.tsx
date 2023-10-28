import { Card } from "antd";
import { useContext } from "react";
import { Link } from "../components/Link";
import { CertPolicyRefTable } from "./CertPolicyRefTable";
import { NamespaceContext } from "./contexts/NamespaceContext";

export default function NamespacePage() {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);

  //  const adminApi = useAuthedClient(AdminApi);

  return (
    <>
      <h1>{namespaceIdentifier}</h1>
      <div>{namespaceKind}</div>
      <Card
        title="Certificate Policies"
        extra={
          <Link to="./cert-policy/_create">Create certificate policy</Link>
        }
      >
        <CertPolicyRefTable routePrefix="./cert-policy/" />
      </Card>
    </>
  );
}
