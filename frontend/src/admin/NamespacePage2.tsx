import { useRequest } from "ahooks";
import { useContext, useState } from "react";
import { Button } from "../components/Button";
import Select, { SelectItem } from "../components/Select";
import {
  AdminApi,
  CertificateTemplateRef,
  LinkedCertificateTemplateUsage,
  NamespaceKind1 as NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { InputField } from "./InputField";
import { NamespaceContext } from "./NamespaceContext2";
import { RefTableColumn, RefsTable } from "./RefsTable";
import { Card } from "antd";
import { Link } from "../components/Link";

const subjectCnColumn: RefTableColumn<CertificateTemplateRef> = {
  columnKey: "subjectCommonName",
  header: "Subject Common Name",
  render: (item) => {
    return item.linkTo ? (
      <span>Link to: {item.linkTo}</span>
    ) : (
      item.subjectCommonName
    );
  },
};
const enabledColumn: RefTableColumn<CertificateTemplateRef> = {
  columnKey: "enabled",
  header: "Enabled",
  render: (item) => (!item.deleted && item.updated ? "Yes" : "No"),
};

type UsageSelectItem = SelectItem<LinkedCertificateTemplateUsage>;

const selectItems: UsageSelectItem[] = [
  {
    id: LinkedCertificateTemplateUsage.LinkedCertificateTemplateUsageClientAuthorization,
    name: "Client Authorization",
  },
  {
    id: LinkedCertificateTemplateUsage.LinkedCertificateTemplateUsageMemberDelegatedEnrollment,
    name: "Member Delegated Enrollment",
  },
];

export default function NamespacePage() {
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);
  return (
    <>
      <h1>{namespaceId}</h1>
      <div>{namespaceKind}</div>
      <Card
        title="Certificate Policies"
        extra={<Link to="./cert-policy/_create">Create certificate policy</Link>}
      ></Card>
    </>
  );
}
