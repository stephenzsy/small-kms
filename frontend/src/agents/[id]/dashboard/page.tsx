import { useRequest } from "ahooks";
import { Button, Card, Form, Input, Table, Tag, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import { useContext, useMemo } from "react";
import { useParams } from "react-router-dom";
import { AgentContext } from "../../../admin/contexts/AgentContext";
import { DrawerContext } from "../../../admin/contexts/DrawerContext";
import { NamespaceContext } from "../../../admin/contexts/NamespaceContext";
import { useNamespace } from "../../../admin/contexts/useNamespace";
import { JsonDataDisplay } from "../../../components/JsonDataDisplay";
import {
  AgentConfigEndpoint,
  AgentConfigName,
  NamespaceProvider,
  PullImageRequest,
} from "../../../generated/apiv2";
import { useAdminApi } from "../../../utils/useCertsApi";
import Select, { DefaultOptionType } from "antd/es/select";
import { ColumnType } from "antd/es/table";
import { StopIcon, TrashIcon } from "@heroicons/react/24/solid";

function DockerPullImageForm({
  instanceId,
  endpointConfig,
}: {
  instanceId: string;
  endpointConfig: AgentConfigEndpoint;
}) {
  const [form] = useForm<PullImageRequest>();
  const { namespaceId } = useNamespace();
  const api = useAdminApi();

  const { run: pullImage, loading } = useRequest(
    async (req: PullImageRequest) => {
      await api?.agentDockerImagePull({
        namespaceId: namespaceId,
        id: instanceId,
        pullImageRequest: req,
      });
    },
    { manual: true }
  );

  const selectItems = useMemo(
    () =>
      endpointConfig?.allowedImageRepos?.map(
        (r): DefaultOptionType => ({
          label: r,
          value: r,
        })
      ),
    [endpointConfig]
  );

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(s) => {
        pullImage({
          imageRepo: s.imageRepo,
          imageTag: s.imageTag,
        });
      }}
    >
      <Form.Item<PullImageRequest> name="imageRepo" label="Image repository">
        <Select options={selectItems} />
      </Form.Item>
      <Form.Item<PullImageRequest> name="imageTag" label="Image tag" required>
        <Input placeholder="latest" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading}>
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}

// type LaunchContainerFormState = {
//   containerName: string;
//   exposedPortSpecs: string[];
//   hostBinds: string[];
//   imageTag: string;
//   listenerAddress: string;
//   pushEndpoint: string;
//   networkName?: string;
//   secrets: SecretMount[];
//   msEntraIdClientCertSecretName?: string;
//   env?: string[];
// };

// function LaunchContainerForm({ mode }: { mode: AgentMode }) {
//   const { namespaceId: namespaceIdentifier, namespaceKind } =
//     useContext(NamespaceContext);
//   const [form] = useForm<LaunchContainerFormState>();
//   const api = useAuthedClient(AdminApi);
//   const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
//   const { run, loading } = useRequest(
//     async (req: LaunchAgentRequest) => {
//       await api.agentLaunchAgent({
//         namespaceId: namespaceIdentifier,
//         namespaceKind,
//         resourceId: instanceId,
//         launchAgentRequest: req,
//         xCryptocatProxyAuthorization: getAccessToken(),
//       });
//     },
//     { manual: true }
//   );

//   const onFinish = useMemoizedFn((values: LaunchContainerFormState) => {
//     if (
//       !values.imageTag ||
//       !values.containerName ||
//       !values.listenerAddress ||
//       !values.pushEndpoint
//     ) {
//       return;
//     }
//     run({
//       containerName: values.containerName,
//       mode,
//       exposedPortSpecs: values.exposedPortSpecs,
//       hostBinds: values.hostBinds,
//       pushEndpoint: values.pushEndpoint,
//       imageTag: values.imageTag,
//       listenerAddress: values.listenerAddress,
//       networkName: values.networkName,
//       secrets: values.secrets,
//       msEntraIdClientCertSecretName: values.msEntraIdClientCertSecretName,
//       env: values.env,
//     });
//   });

//   return (
//     <>
//       <div className="mb-4">Mode: {mode}</div>
//       <Form<LaunchContainerFormState>
//         form={form}
//         layout="vertical"
//         initialValues={{
//           containerName:
//             mode === AgentMode.AgentModeLauncher
//               ? "cryptocat-agent-launcher"
//               : "cryptocat-agent",
//           exposedPortSpecs: ["11443:11443"],
//           hostBinds: [
//             "/opt/smallkms/config:/opt/smallkms/config:rw",
//             "/var/run/docker.sock:/var/run/docker.sock:rw",
//           ],
//           imageTag: "",
//           listenerAddress: ":10443",
//           networkName: "",
//           secrets: [
//             {
//               source: "/opt/smallkms/sp-client-cert.pem",
//               targetName: "aad-client-creds.pem",
//             },
//           ],
//         }}
//         onFinish={onFinish}
//       >
//         <Form.Item<LaunchContainerFormState>
//           name="containerName"
//           label="Container name"
//           required
//         >
//           <Input placeholder="cryptocat-agent" />
//         </Form.Item>
//         <Form.Item<LaunchContainerFormState>
//           name="imageTag"
//           label="Image tag"
//           required
//         >
//           <Input placeholder="latest" />
//         </Form.Item>
//         <Form.Item<LaunchContainerFormState>
//           name="listenerAddress"
//           label="Listener address"
//           required
//         >
//           <Input placeholder=":11443" />
//         </Form.Item>
//         <Form.Item<LaunchContainerFormState>
//           name="pushEndpoint"
//           label="Push endpoint"
//           required
//         >
//           <Input placeholder="https://localhost:11443" />
//         </Form.Item>
//         <Form.Item<LaunchContainerFormState>
//           name="msEntraIdClientCertSecretName"
//           label="Microsoft Entra ID client certificate secret name"
//         >
//           <Input />
//         </Form.Item>
//         <Form.Item<LaunchContainerFormState>
//           name="networkName"
//           label="Network name"
//         >
//           <Input placeholder="" />
//         </Form.Item>
//         <Form.List name={"exposedPortSpecs"}>
//           {(subFields, subOpt) => {
//             return (
//               <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md">
//                 <div className="text-lg font-semibold">Port bindings</div>
//                 {subFields.map((subField) => (
//                   <div key={subField.key} className="flex items-center gap-4">
//                     <Form.Item
//                       noStyle
//                       name={subField.name}
//                       className="flex-auto"
//                     >
//                       <Input placeholder={"11443:11443"} />
//                     </Form.Item>
//                     <Button
//                       type="text"
//                       onClick={() => {
//                         subOpt.remove(subField.name);
//                       }}
//                     >
//                       <XMarkIcon className="h-em w-em" />
//                     </Button>
//                   </div>
//                 ))}
//                 <Button type="dashed" onClick={() => subOpt.add()} block>
//                   Add exposed port
//                 </Button>
//               </div>
//             );
//           }}
//         </Form.List>
//         <Form.List name={"hostBinds"}>
//           {(subFields, subOpt) => {
//             return (
//               <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
//                 <div className="text-lg font-semibold">Host bindings</div>
//                 {subFields.map((subField) => (
//                   <div key={subField.key} className="flex items-center gap-4">
//                     <Form.Item
//                       noStyle
//                       name={subField.name}
//                       className="flex-auto"
//                     >
//                       <Input placeholder={"source:target"} />
//                     </Form.Item>
//                     <Button
//                       type="text"
//                       onClick={() => {
//                         subOpt.remove(subField.name);
//                       }}
//                     >
//                       <XMarkIcon className="h-em w-em" />
//                     </Button>
//                   </div>
//                 ))}
//                 <Button type="dashed" onClick={() => subOpt.add()} block>
//                   Add host binding
//                 </Button>
//               </div>
//             );
//           }}
//         </Form.List>
//         <Form.List name={"secrets"}>
//           {(subFields, subOpt) => {
//             return (
//               <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
//                 <div className="text-lg font-semibold">Secret bindings</div>
//                 {subFields.map((subField) => (
//                   <div key={subField.key} className="flex items-center gap-4">
//                     <Form.Item
//                       noStyle
//                       name={[subField.name, "targetName"]}
//                       className="flex-auto"
//                       label="Name"
//                     >
//                       <Input
//                         placeholder={"source:target"}
//                         addonBefore={"Name"}
//                       />
//                     </Form.Item>
//                     <Form.Item
//                       noStyle
//                       name={[subField.name, "source"]}
//                       className="flex-auto"
//                       label="Source"
//                     >
//                       <Input
//                         placeholder={"source:target"}
//                         addonBefore={"Source"}
//                       />
//                     </Form.Item>
//                     <Button
//                       type="text"
//                       onClick={() => {
//                         subOpt.remove(subField.name);
//                       }}
//                     >
//                       <XMarkIcon className="h-em w-em" />
//                     </Button>
//                   </div>
//                 ))}
//                 <Button type="dashed" onClick={() => subOpt.add()} block>
//                   Add secret binding
//                 </Button>
//               </div>
//             );
//           }}
//         </Form.List>
//         <Form.List name={"env"}>
//           {(subFields, subOpt) => {
//             return (
//               <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
//                 <div className="text-lg font-semibold">
//                   Enviornment variables
//                 </div>
//                 {subFields.map((subField) => (
//                   <div key={subField.key} className="flex items-center gap-4">
//                     <Form.Item
//                       noStyle
//                       name={[subField.name]}
//                       className="flex-auto"
//                     >
//                       <Input placeholder={"FOO=BAR"} />
//                     </Form.Item>
//                     <Button
//                       type="text"
//                       onClick={() => {
//                         subOpt.remove(subField.name);
//                       }}
//                     >
//                       <XMarkIcon className="h-em w-em" />
//                     </Button>
//                   </div>
//                 ))}
//                 <Button type="dashed" onClick={() => subOpt.add()} block>
//                   Add environment variable
//                 </Button>
//               </div>
//             );
//           }}
//         </Form.List>
//         <Form.Item className="mt-4">
//           <Button type="primary" htmlType="submit" loading={loading}>
//             Launch
//           </Button>
//         </Form.Item>
//       </Form>
//     </>
//   );
// }

type DockerContainer = {
  Id: string;
  Image: string;
  Command: string;
  Created: number;
  Status: string;
  State: string;
  Ports: {
    IP: string;
    PrivatePort: number;
    PublicPort: number;
    Type: string;
  }[];
  Names: string[];
};

function useDockerContainerColumns(
  onInspect: (id: string) => void,
  onStop: (id: string) => void,
  onRemove: (id: string) => void
): ColumnType<DockerContainer>[] {
  return useMemo(
    () => [
      {
        title: "ID",
        dataIndex: "Id",
        key: "Id",
        render: (id: DockerContainer["Id"]) => {
          return <span className="font-mono">{id.substring(0, 12)}</span>;
        },
      },
      {
        title: "Image",
        dataIndex: "Image",
        key: "Image",
        className: "max-w-[200px]",
      },
      {
        title: "Command",
        dataIndex: "Command",
        key: "Command",
        className: "font-mono",
      },
      {
        title: "Created",
        dataIndex: "Created",
        key: "Created",
        render: (created: DockerContainer["Created"]) => {
          const d = new Date(created * 1000);
          return (
            <time dateTime={d.toISOString()} className="font-mono">
              {d.toISOString()}
            </time>
          );
        },
      },
      {
        title: "Status",
        dataIndex: "Status",
        key: "Status",
      },
      {
        title: "Ports",
        dataIndex: "Ports",
        key: "Ports",
        render: (ports: DockerContainer["Ports"]) => {
          return ports.map((port) => (
            <div key={port.PrivatePort}>
              {port.IP}:{port.PrivatePort}
              {" -> "}
              {port.PublicPort}
            </div>
          ));
        },
      },
      {
        title: "Names",
        dataIndex: "Names",
        key: "Names",
        render: (names: DockerContainer["Names"]) => {
          return names.map((name) => <div key={name}>{name}</div>);
        },
      },
      {
        title: "Actions",
        key: "actions",
        render: (c: DockerContainer) => {
          return (
            <div className="flex gap-2">
              <Button
                size="small"
                onClick={() => {
                  onInspect(c.Id);
                }}
              >
                Inspect
              </Button>
              {c.State == "running" ? (
                <Button
                  size="small"
                  className="flex items-center"
                  danger
                  onClick={() => {
                    onStop(c.Id);
                  }}
                  icon={<StopIcon className="h-em w-em" />}
                >
                  Stop
                </Button>
              ) : (
                <Button
                  size="small"
                  className="flex items-center"
                  danger
                  onClick={() => {
                    onRemove(c.Id);
                  }}
                  icon={<TrashIcon className="h-4 w-4" />}
                >
                  Remove
                </Button>
              )}
            </div>
          );
        },
      },
    ],
    [onInspect, onStop, onRemove]
  );
}

function ContainersTableCard({ instanceId }: { instanceId: string }) {
  const { openDrawer } = useContext(DrawerContext);
  const { namespaceId } = useNamespace();
  const api = useAdminApi();
  const { data: containers, refresh: refreshContainersList } = useRequest(
    async (): Promise<DockerContainer[] | undefined> => {
      return (await api?.agentDockerContainerList({
        namespaceId,
        id: instanceId,
      })) as DockerContainer[];
    },
    { refreshDeps: [namespaceId, instanceId] }
  );

  // const onInspect = useMemoizedFn((id: string) => {
  //   openDrawer({
  //     title: `Inspect ${id}`,
  //     onGetData: async () => {
  //       return await api.agentDockerContainerInspect({
  //         namespaceId: namespaceIdentifier,
  //         namespaceKind,
  //         resourceId: instanceId,
  //         xCryptocatProxyAuthorization: getAccessToken(),
  //         containerId: id,
  //       });
  //     },
  //   });
  // });

  // const { run: onStop } = useRequest(
  //   async (containerId: string): Promise<void> => {
  //     return await api.agentDockerContainerStop({
  //       namespaceId: namespaceIdentifier,
  //       namespaceKind,
  //       resourceId: instanceId,
  //       containerId,
  //       xCryptocatProxyAuthorization: getAccessToken(),
  //     });
  //   },
  //   { manual: true }
  // );

  // const { run: onRemove } = useRequest(
  //   async (containerId: string): Promise<void> => {
  //     return await api.agentDockerContainerRemove({
  //       namespaceId: namespaceIdentifier,
  //       namespaceKind,
  //       resourceId: instanceId,
  //       containerId,
  //       xCryptocatProxyAuthorization: getAccessToken(),
  //     });
  //   },
  //   { manual: true }
  // );

  const columns = useDockerContainerColumns(
    () => {},
    () => {},
    () => {}
  );
  return (
    <Card
      title="Containers"
      extra={
        <div className="flex flex-row gap-4 items-center">
          <Button type="link" onClick={refreshContainersList}>
            Refresh
          </Button>
          <Button
            type="link"
            onClick={() =>
              openDrawer(<JsonDataDisplay data={containers} />, {
                title: "Containers",
              })
            }
          >
            View JSON
          </Button>
        </div>
      }
    >
      <Table<DockerContainer>
        columns={columns}
        dataSource={containers}
        rowKey={(c) => c.Id}
      />
    </Card>
  );
}

// function RadiusConfigurationCard({ api }: { api: AdminApi }) {
//   const { namespaceId: namespaceIdentifier, namespaceKind } =
//     useContext(NamespaceContext);
//   const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
//   const { run: pushAgentConfigRadius } = useRequest(
//     async () => {
//       if (instanceId) {
//         const result = await api.pushAgentConfigRadius({
//           namespaceKind,
//           namespaceId: namespaceIdentifier,
//           resourceId: instanceId,
//           xCryptocatProxyAuthorization: getAccessToken(),
//         });

//         return result;
//       }
//     },
//     {
//       manual: true,
//     }
//   );
//   return (
//     <Card title="RADIUS configuration">
//       <div>
//         <Button type="primary" onClick={pushAgentConfigRadius}>
//           Push configuration
//         </Button>
//       </div>
//     </Card>
//   );
// }

type InstanceProps = {
  namespaceId: string;
  instanceId: string;
};

function DockerSystemInfo({ namespaceId, instanceId }: InstanceProps) {
  const api = useAdminApi();
  const { data } = useRequest(
    async () => {
      return await api?.getAgentDockerSystemInformation({
        namespaceId,
        id: instanceId,
      });
    },
    {
      refreshDeps: [namespaceId, instanceId],
    }
  );

  return <JsonDataDisplay data={data} />;
}

function DockerImagesList({ namespaceId, instanceId }: InstanceProps) {
  const api = useAdminApi();
  const { data } = useRequest(
    async () => {
      return await api?.agentDockerImageList({
        namespaceId,
        id: instanceId,
      });
    },
    {
      refreshDeps: [namespaceId, instanceId],
    }
  );

  return <JsonDataDisplay data={data} />;
}

function DockerNetworksList({ namespaceId, instanceId }: InstanceProps) {
  const api = useAdminApi();
  const { data } = useRequest(
    async () => {
      return await api?.listAgentDockerNetowks({
        namespaceId,
        id: instanceId,
      });
    },
    {
      refreshDeps: [namespaceId, instanceId],
    }
  );

  return <JsonDataDisplay data={data} />;
}

function AgentDashboard({ instanceId }: { instanceId: string }) {
  const api = useAdminApi();
  const { namespaceId } = useNamespace();
  const { openDrawer } = useContext(DrawerContext);
  const { data: instance } = useRequest(
    async () => {
      if (namespaceId) {
        return api?.getAgentInstance({
          namespaceId: namespaceId,
          id: instanceId,
        });
      }
    },
    {
      refreshDeps: [instanceId, namespaceId],
    }
  );

  const { data: endpointConfig } = useRequest(
    async () => {
      return (await api?.getAgentConfig({
        configName: AgentConfigName.AgentConfigNameEndpoint,
        namespaceId: namespaceId,
      })) as AgentConfigEndpoint;
    },
    {
      refreshDeps: [instanceId, namespaceId],
    }
  );

  const { run: deleteInstance } = useRequest(
    async (instanceId: string) => {
      if (namespaceId) {
        return api?.deleteAgentInstance({
          namespaceId: namespaceId,
          id: instanceId,
        });
      }
    },
    {
      manual: true,
    }
  ); // const { namespaceId: namespaceIdentifier, namespaceKind } =
  //   useContext(NamespaceContext);
  // const { instanceId, hasToken, getAccessToken, setAccessToken } = useContext(
  //   ProxyAuthTokenContext
  // );
  // const { openDrawer } = useContext(DataDrawerContext);

  // const { data } = useRequest(
  //   async () => {
  //     if (instanceId) {
  //       const result = await api.getAgentInstance({
  //         namespaceKind,
  //         namespaceId: namespaceIdentifier,
  //         resourceId: instanceId,
  //       });

  //       return result;
  //     }
  //   },
  //   {
  //     refreshDeps: [namespaceKind, namespaceIdentifier, instanceId],
  //   }
  // );

  // const { run: acquireToken } = useRequest(
  //   async () => {
  //     if (instanceId) {
  //       const result = await api.createAgentInstanceProxyAuthToken({
  //         namespaceId: namespaceIdentifier,
  //         namespaceKind,
  //         resourceId: instanceId,
  //       });
  //       setAccessToken(result.accessToken);
  //       return result;
  //     }
  //   },
  //   {
  //     manual: true,
  //   }
  // );

  // const getAgentDiagnostics = useMemoizedFn(async () => {
  //   return await api.getAgentDiagnostics({
  //     namespaceId: namespaceIdentifier,
  //     namespaceKind,
  //     resourceId: instanceId,
  //     xCryptocatProxyAuthorization: getAccessToken(),
  //   });
  // });

  return (
    <>
      <Card title="Agent information">
        <dl className="dl">
          <div>
            <dt>Instance ID</dt>
            <dd className="font-mono">{instance?.id}</dd>
          </div>
          <div>
            <dt>Build ID</dt>
            <dd className="font-mono">{instance?.buildId}</dd>
          </div>
          <div>
            <dt>Config version</dt>
            <dd className="font-mono">
              {instance?.configVersion}{" "}
              {instance?.configVersion &&
                instance.configVersion === endpointConfig?.version && (
                  <Tag color="success">Up to date</Tag>
                )}
            </dd>
          </div>
          <div>
            <dt>State</dt>
            <dd className="capitalize">{instance?.state}</dd>
          </div>
        </dl>
        <div>
          <Button
            danger
            onClick={() => {
              deleteInstance(instanceId);
            }}
          >
            Delete
          </Button>
        </div>
      </Card>
      <ContainersTableCard instanceId={instanceId} />
      <Card title="View">
        <div className="flex flex-col gap-4 items-start">
          <Button
            type="link"
            onClick={() => {
              openDrawer(
                <DockerImagesList
                  instanceId={instanceId}
                  namespaceId={namespaceId}
                />,
                {
                  title: "Docker images",
                  size: "large",
                }
              );
            }}
          >
            List Docker images
          </Button>
          <Button
            type="link"
            onClick={() => {
              openDrawer(
                <DockerNetworksList
                  instanceId={instanceId}
                  namespaceId={namespaceId}
                />,
                {
                  title: "Docker networks",
                  size: "large",
                }
              );
            }}
          >
            List Docker networks
          </Button>
          <Button
            type="link"
            onClick={() => {
              openDrawer(
                <DockerSystemInfo
                  instanceId={instanceId}
                  namespaceId={namespaceId}
                />,
                {
                  title: "Docker system information",
                  size: "large",
                }
              );
            }}
          >
            Docker system info
          </Button>
          {/* <Button
            type="link"
            onClick={() => {
              openDrawer({
                title: "Agent diagnostics",
                onGetData: getAgentDiagnostics,
              });
            }}
            disabled={!hasToken}
          >
            Agent request diagnostics
          </Button> */}
        </div>
      </Card>
      <Card title="Pull Image">
        {endpointConfig && (
          <DockerPullImageForm
            instanceId={instanceId}
            endpointConfig={endpointConfig}
          />
        )}
      </Card>

      {/* {hasToken && (
        <>
          <Card title="Docker pull image">
            <DockerPullImageForm />
          </Card>
          <RadiusConfigurationCard api={api} />
          {data && (
            <Card title="Launch container">
              <LaunchContainerForm
                mode={data?.mode === "server" ? "launcher" : "server"}
              />
            </Card>
          )}
        </>
      )} */}
    </>
  );
}

export default function AgentDashboardPage() {
  const { instanceId } = useParams<{ instanceId: string }>();
  const { agent } = useContext(AgentContext);

  return (
    <>
      <Typography.Title>Agent Dashboard: {agent?.displayName}</Typography.Title>
      <NamespaceContext.Provider
        value={{
          namespaceId: agent?.servicePrincipalId ?? "",
          namespaceKind: NamespaceProvider.NamespaceProviderServicePrincipal,
        }}
      >
        {instanceId && agent?.servicePrincipalId && <AgentDashboard instanceId={instanceId} />}
      </NamespaceContext.Provider>
    </>
  );
}
