import { useRequest } from "ahooks";
import { Button, Card, Form, Input, SelectProps, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import Select, { DefaultOptionType } from "antd/es/select";
import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { useNamespace } from "../../admin/contexts/useNamespace";
import { CertificateExternalIssuerAcme } from "../../generated/apiv2";
import { useAdminApi } from "../../utils/useCertsApi";
import { XMarkIcon } from "@heroicons/react/24/outline";

export default function CertificateIssuerPage() {
  const { namespaceProvider, namespaceId } = useNamespace();
  const { id } = useParams<{ id: string }>();

  return (
    <>
      <Typography.Title>Certificate Issuer</Typography.Title>
      <div className="font-mono">
        {namespaceProvider}:{namespaceId}:cert-policy/{id}
      </div>
      <Card title="Create or update issuer">
        {id && <AcmeIssuerForm issuerId={id} />}
      </Card>
    </>
  );
}

function KeySelect({
  value,
  onChange,
}: {
  value?: string;
  onChange?: SelectProps<string>["onChange"];
}) {
  const api = useAdminApi();
  const { namespaceProvider, namespaceId } = useNamespace();

  const { data: keys } = useRequest(
    async () => {
      return api?.listKeys({
        namespaceId,
        namespaceProvider,
      });
    },
    {
      refreshDeps: [namespaceId, namespaceProvider],
    }
  );

  const keyOptions = useMemo((): DefaultOptionType[] | undefined => {
    return keys?.map((key) => ({
      label: (
        <span>
          {key.id} ({key.policyIdentifier})
        </span>
      ),
      value: key.id,
    }));
  }, [keys]);

  return (
    <>
      <Form.Item label="Select key" required>
        <Select options={keyOptions} value={value} onChange={onChange} />
      </Form.Item>
    </>
  );
}

function AcmeIssuerForm({ issuerId }: { issuerId: string }) {
  const [formInstance] = useForm<CertificateExternalIssuerAcme>();
  const { namespaceId } = useNamespace();
  const api = useAdminApi();

  const { run } = useRequest(
    async (acme: CertificateExternalIssuerAcme) => {
      return api?.putExternalCertificateIssuer({
        id: issuerId,
        namespaceId: namespaceId,
        certificateExternalIssuerFields: {
          acme,
        },
      });
    },
    {
      manual: true,
    }
  );

  return (
    <Form
      layout="vertical"
      form={formInstance}
      onFinish={(values) => {
        run(values);
      }}
    >
      <Form.Item<CertificateExternalIssuerAcme>
        name="directoryUrl"
        label="Directory URL"
        required
      >
        <Input type="url" inputMode="url" />
      </Form.Item>
      <Form.Item<CertificateExternalIssuerAcme> name="accountKeyId">
        <KeySelect />
      </Form.Item>
      <Form.Item<CertificateExternalIssuerAcme> label="Contacts">
        <Form.List name="contacts">
          {(subFields, subOpt) => {
            return (
              <div className="flex flex-col gap-4">
                {subFields.map((subField) => (
                  <div key={subField.key} className="flex items-center gap-4">
                    <Form.Item
                      noStyle
                      name={subField.name}
                      className="flex-auto"
                    >
                      <Input
                        placeholder="email@example.com"
                        type="email"
                        inputMode="email"
                      />
                    </Form.Item>
                    <Button
                      type="text"
                      onClick={() => {
                        subOpt.remove(subField.name);
                      }}
                    >
                      <XMarkIcon className="h-em w-em" />
                    </Button>
                  </div>
                ))}
                <Button type="dashed" onClick={() => subOpt.add()} block>
                  Add contact email
                </Button>
              </div>
            );
          }}
        </Form.List>
      </Form.Item>
      <Form.Item<CertificateExternalIssuerAcme>
        label="Azure DNS Zone Resource ID"
        name="azureDnsZoneResourceId"
        required
      >
        <Input />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Register or import
        </Button>
      </Form.Item>
    </Form>
  );
}
