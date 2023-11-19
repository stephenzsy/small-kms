import { Button, Card, Form, Input, Select, Typography } from "antd";
import { AdminApi, NamespaceProvider } from "../generated/apiv2";
import { useAuthedClientV2, useGraphClient } from "../utils/useCertsApi";
import { useRequest } from "ahooks";
import { RefsTable } from "./RefsTable";
import { ResourceRefsTable } from "./tables/ResourceRefsTable";
import { GraphRequest } from "@microsoft/microsoft-graph-client";
import { useForm } from "antd/es/form/Form";
import { useMemo } from "react";
import { DefaultOptionType } from "antd/es/select";
import { ArrowPathIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";

export default function ProfilesPage({
  namespaceProvider,
  title,
}: {
  namespaceProvider: NamespaceProvider;
  title: React.ReactNode;
}) {
  const api = useAuthedClientV2(AdminApi);
  const { data, loading, refresh } = useRequest(
    () => {
      return api.listProfiles({
        namespaceProvider: namespaceProvider,
      });
    },
    { refreshDeps: [namespaceProvider] }
  );

  const graphClient = useGraphClient();
  const {
    data: dirObj,
    run: getDirectoryObjects,
    loading: dirObjLoading,
  } = useRequest(
    async (): Promise<[NamespaceProvider, any]> => {
      let builder: GraphRequest;
      switch (namespaceProvider) {
        case NamespaceProvider.NamespaceProviderServicePrincipal:
          builder = graphClient.api("/servicePrincipals");
          break;
        case NamespaceProvider.NamespaceProviderUser:
          builder = graphClient.api("/users");
          break;
        case NamespaceProvider.NamespaceProviderGroup:
          builder = graphClient.api("/groups");
          break;
        default:
          return ["" as NamespaceProvider, []];
      }
      return [
        namespaceProvider,
        await builder.select(["id", "displayName"]).get(),
      ];
    },
    {
      manual: true,
      refreshDeps: [namespaceProvider],
    }
  );

  const [importForm] = useForm<{ id: string }>();

  const dirObjOptions = useMemo((): DefaultOptionType[] | undefined => {
    if (!dirObj) {
      return undefined;
    }
    const [reqProvider, directoryObjects] = dirObj;
    if (namespaceProvider !== reqProvider) {
      return undefined;
    }
    return directoryObjects.value.map((v: any) => {
      return {
        value: v.id,
        label: `${v.displayName} (${v.id})`,
      };
    });
  }, [dirObj, namespaceProvider]);

  const { run: importProfile } = useRequest(
    async (id: string) => {
      await api.syncProfile({
        namespaceProvider: namespaceProvider,
        namespaceId: id,
      });
      refresh();
    },
    {
      manual: true,
    }
  );

  return (
    <>
      <Typography.Title>{title}</Typography.Title>

      <Card title="Profiles">
        <ResourceRefsTable dataSource={data} loading={loading} />
      </Card>

      <Card title="Import Profile">
        <Form
          form={importForm}
          layout="vertical"
          onFinish={(values) => {
            importProfile(values.id);
          }}
        >
          <Form.Item
            name="id"
            label={
              <div className="inline-flex items-center gap-4">
                <span>Select graph object</span>
                <Button
                  type="link"
                  size="small"
                  onClick={getDirectoryObjects}
                  className="inline-flex items-center gap-2"
                >
                  <ArrowPathIcon
                    className={classNames(
                      "h-em w-em",
                      dirObjLoading && "animate-spin"
                    )}
                  />
                  <span>Get List from Microsoft Graph API</span>
                </Button>
              </div>
            }
          >
            <Select options={dirObjOptions} loading={dirObjLoading} />
          </Form.Item>
          <Form.Item name="id" label={"Enter Object ID"}>
            <Input />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              Import
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </>
  );
}
