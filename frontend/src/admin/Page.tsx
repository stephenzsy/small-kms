import { Card, CardSection } from "../components/Card";
import { NamespaceKind1 as NamespaceKind } from "../generated";

export default function AdminPage() {
  return (
    <>
      {(
        [
          [NamespaceKind.NamespaceKindCaRoot, "Root Certificate Authorities"],
          [
            NamespaceKind.NamespaceKindCaInt,
            "Intermediate Certificate Authorities",
          ],
          [NamespaceKind.NamespaceKindServicePrincipal, "Service Principals"],
          [NamespaceKind.NamespaceKindGroup, "Groups"],
          [NamespaceKind.NamespaceKindDevice, "Devices"],
          [NamespaceKind.NamespaceKindUser, "Users"],
          [NamespaceKind.NamespaceKindApplication, "Applications"],
        ] as Array<[NamespaceKind, string]>
      ).map(([t, title]: [NamespaceKind, string]) => (
        <Card key={t}>
          <CardSection></CardSection>
        </Card>
      ))}
    </>
  );
}
