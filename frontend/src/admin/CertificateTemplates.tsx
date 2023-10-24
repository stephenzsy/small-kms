import { useRequest } from "ahooks";
import type { Result } from "ahooks/lib/useRequest/src/types";
import { Button, Form, Input, Select, Table, type TableColumnType } from "antd";
import { useForm } from "antd/es/form/Form";
import { DefaultOptionType } from "antd/es/select";
import React, { useContext } from "react";
import { Link } from "../components/Link";
import {
  AdminApi,
  CertificateTemplateRef,
  LinkedCertificateTemplateUsage,
  NamespaceKind1 as NamespaceKind,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";

export interface CertificateTemplatesContextData
  extends Pick<
    Result<CertificateTemplateRef[], []>,
    "data" | "run" | "loading"
  > {
  namespaceKind: NamespaceKind;
  namespaceId: string;
}

const CertificateTemplatesContext =
  React.createContext<CertificateTemplatesContextData>({
    data: undefined,
    run: () => {},
    loading: false,
    namespaceKind: "" as never,
    namespaceId: "",
  });

function useColumns(
  nsKind: NamespaceKind,
  nsId: string
): TableColumnType<CertificateTemplateRef>[] {
  return [
    {
      key: "id",
      title: "ID",
      render: (item) => {
        return <pre className="font-mono">{item.id}</pre>;
      },
    },
    {
      key: "name",
      title: "Name",
      render: (item: CertificateTemplateRef) => {
        return item.linkTo ? (
          <span>Link to: {item.linkTo}</span>
        ) : (
          item.subjectCommonName
        );
      },
    },
    {
      key: "enabled",
      title: "Enabled",
      render: (item: CertificateTemplateRef) => {
        return !item.deleted && item.updated ? "Yes" : "No";
      },
    },
    {
      key: "actions",
      title: "Actions",
      render: (item: CertificateTemplateRef) => {
        return (
          <Link
            to={`/admin/${nsKind}/${nsId}/certificate-templates/${item.id}`}
          >
            View
          </Link>
        );
      },
    },
  ];
}

export function CertificateTemplateTable() {
  const { data, loading, namespaceKind, namespaceId } = React.useContext(
    CertificateTemplatesContext
  );
  const columns = useColumns(namespaceKind, namespaceId);
  return (
    <Table<CertificateTemplateRef>
      loading={loading}
      dataSource={data}
      columns={columns}
      rowKey={(item) => item.id}
    />
  );
}

type LinkTemplateCardFormState = {
  linkTarget?: string;
  linkUsage?: LinkedCertificateTemplateUsage;
};

const linkTemplateUsageSelectOptions: DefaultOptionType[] = [
  {
    value:
      LinkedCertificateTemplateUsage.LinkedCertificateTemplateUsageClientAuthorization,
    label: "Client Authorization",
  },
  {
    value:
      LinkedCertificateTemplateUsage.LinkedCertificateTemplateUsageMemberDelegatedEnrollment,
    label: "Member Delegated Enrollment",
  },
];

export function LinkTemplateForm() {
  const {
    namespaceId,
    namespaceKind,
    run: refresh,
  } = useContext(CertificateTemplatesContext);
  const [form] = useForm<LinkTemplateCardFormState>();
  const adminApi = useAuthedClient(AdminApi);
  const { run: createLink } = useRequest(
    async (
      targetLocator: string,
      selectedUsage: LinkedCertificateTemplateUsage
    ) => {
      targetLocator = targetLocator.trim();
      if (!targetLocator) {
        return;
      }
      await adminApi.createLinkedCertificateTemplate({
        namespaceId,
        namespaceKind,
        createLinkedCertificateTemplateParameters: {
          targetTemplate: targetLocator,
          usage: selectedUsage,
        },
      });

      refresh();
    },
    { manual: true }
  );
  return (
    <Form<LinkTemplateCardFormState>
      form={form}
      labelCol={{ span: 8 }}
      wrapperCol={{ span: 16 }}
      style={{ maxWidth: 600 }}
      onFinish={(values) => {
        const { linkTarget, linkUsage } = values;
        if (linkTarget && linkUsage) {
          createLink(linkTarget, linkUsage);
        }
      }}
      initialValues={{
        linkUsage: linkTemplateUsageSelectOptions[0].value,
      }}
    >
      <Form.Item<LinkTemplateCardFormState>
        label="Link to"
        name="linkTarget"
        required
      >
        <Input />
      </Form.Item>
      <Form.Item<LinkTemplateCardFormState> label="Usage" name="linkUsage">
        <Select options={linkTemplateUsageSelectOptions} />
      </Form.Item>

      <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
        <Button htmlType="submit" type="primary">
          Add link
        </Button>
      </Form.Item>
    </Form>
  );
}
