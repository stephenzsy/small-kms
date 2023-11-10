import { StopIcon, TrashIcon, XMarkIcon } from "@heroicons/react/24/solid";
import { useBoolean, useMemoizedFn, useRequest } from "ahooks";
import { Button, Card, Drawer, Form, Input, Table, Typography } from "antd";
import { useForm } from "antd/es/form/Form";
import { ColumnType } from "antd/es/table";
import {
  PropsWithChildren,
  createContext,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";
import { useParams } from "react-router-dom";
import {
  JsonDataDisplay,
  JsonDataDisplayProps,
} from "../components/JsonDataDisplay";
import {
  AdminApi,
  AgentMode,
  LaunchAgentRequest,
  PullImageRequest,
  SecretMount,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import {
  ProxyAuthTokenContext,
  ProxyAuthTokenContextProvider,
} from "./ProxyAuthTokenContext";
import { NamespaceContext } from "./contexts/NamespaceContext";

type DockerPullImageFormState = {
  imageTag: string;
  imageRepo?: string;
};

function DockerPullImageForm() {
  const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
  const [form] = useForm<DockerPullImageFormState>();
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const api = useAuthedClient(AdminApi);

  const { run: pullImage, loading } = useRequest(
    async (req: PullImageRequest) => {
      await api.agentPullImage({
        namespaceId: namespaceIdentifier,
        namespaceKind,
        resourceId: instanceId,
        pullImageRequest: req,
        xCryptocatProxyAuthorization: getAccessToken(),
      });
    },
    { manual: true }
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
      <Form.Item<DockerPullImageFormState>
        name="imageRepo"
        label="Image repository"
      >
        <Input />
      </Form.Item>
      <Form.Item<DockerPullImageFormState>
        name="imageTag"
        label="Image tag"
        required
      >
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

type LaunchContainerFormState = {
  containerName: string;
  exposedPortSpecs: string[];
  hostBinds: string[];
  imageTag: string;
  listenerAddress: string;
  pushEndpoint: string;
  networkName?: string;
  secrets: SecretMount[];
  msEntraIdClientCertSecretName?: string;
  env?: string[];
};

function LaunchContainerForm({ mode }: { mode: AgentMode }) {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const [form] = useForm<LaunchContainerFormState>();
  const api = useAuthedClient(AdminApi);
  const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
  const { run, loading } = useRequest(
    async (req: LaunchAgentRequest) => {
      await api.agentLaunchAgent({
        namespaceId: namespaceIdentifier,
        namespaceKind,
        resourceId: instanceId,
        launchAgentRequest: req,
        xCryptocatProxyAuthorization: getAccessToken(),
      });
    },
    { manual: true }
  );

  const onFinish = useMemoizedFn((values: LaunchContainerFormState) => {
    if (
      !values.imageTag ||
      !values.containerName ||
      !values.listenerAddress ||
      !values.pushEndpoint
    ) {
      return;
    }
    run({
      containerName: values.containerName,
      mode,
      exposedPortSpecs: values.exposedPortSpecs,
      hostBinds: values.hostBinds,
      pushEndpoint: values.pushEndpoint,
      imageTag: values.imageTag,
      listenerAddress: values.listenerAddress,
      networkName: values.networkName,
      secrets: values.secrets,
      msEntraIdClientCertSecretName: values.msEntraIdClientCertSecretName,
      env: values.env,
    });
  });

  return (
    <>
      <div className="mb-4">Mode: {mode}</div>
      <Form<LaunchContainerFormState>
        form={form}
        layout="vertical"
        initialValues={{
          containerName:
            mode === AgentMode.AgentModeLauncher
              ? "cryptocat-agent-launcher"
              : "cryptocat-agent",
          exposedPortSpecs: ["11443:11443"],
          hostBinds: [
            "/opt/smallkms/config:/opt/smallkms/config:rw",
            "/var/run/docker.sock:/var/run/docker.sock:rw",
          ],
          imageTag: "",
          listenerAddress: ":10443",
          networkName: "",
          secrets: [
            {
              source: "/opt/smallkms/sp-client-cert.pem",
              targetName: "aad-client-creds.pem",
            },
          ],
        }}
        onFinish={onFinish}
      >
        <Form.Item<LaunchContainerFormState>
          name="containerName"
          label="Container name"
          required
        >
          <Input placeholder="cryptocat-agent" />
        </Form.Item>
        <Form.Item<LaunchContainerFormState>
          name="imageTag"
          label="Image tag"
          required
        >
          <Input placeholder="latest" />
        </Form.Item>
        <Form.Item<LaunchContainerFormState>
          name="listenerAddress"
          label="Listener address"
          required
        >
          <Input placeholder=":11443" />
        </Form.Item>
        <Form.Item<LaunchContainerFormState>
          name="pushEndpoint"
          label="Push endpoint"
          required
        >
          <Input placeholder="https://localhost:11443" />
        </Form.Item>
        <Form.Item<LaunchContainerFormState>
          name="msEntraIdClientCertSecretName"
          label="Microsoft Entra ID client certificate secret name"
        >
          <Input />
        </Form.Item>
        <Form.Item<LaunchContainerFormState>
          name="networkName"
          label="Network name"
        >
          <Input placeholder="" />
        </Form.Item>
        <Form.List name={"exposedPortSpecs"}>
          {(subFields, subOpt) => {
            return (
              <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md">
                <div className="text-lg font-semibold">Port bindings</div>
                {subFields.map((subField) => (
                  <div key={subField.key} className="flex items-center gap-4">
                    <Form.Item
                      noStyle
                      name={subField.name}
                      className="flex-auto"
                    >
                      <Input placeholder={"11443:11443"} />
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
                  Add exposed port
                </Button>
              </div>
            );
          }}
        </Form.List>
        <Form.List name={"hostBinds"}>
          {(subFields, subOpt) => {
            return (
              <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
                <div className="text-lg font-semibold">Host bindings</div>
                {subFields.map((subField) => (
                  <div key={subField.key} className="flex items-center gap-4">
                    <Form.Item
                      noStyle
                      name={subField.name}
                      className="flex-auto"
                    >
                      <Input placeholder={"source:target"} />
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
                  Add host binding
                </Button>
              </div>
            );
          }}
        </Form.List>
        <Form.List name={"secrets"}>
          {(subFields, subOpt) => {
            return (
              <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
                <div className="text-lg font-semibold">Secret bindings</div>
                {subFields.map((subField) => (
                  <div key={subField.key} className="flex items-center gap-4">
                    <Form.Item
                      noStyle
                      name={[subField.name, "targetName"]}
                      className="flex-auto"
                      label="Name"
                    >
                      <Input
                        placeholder={"source:target"}
                        addonBefore={"Name"}
                      />
                    </Form.Item>
                    <Form.Item
                      noStyle
                      name={[subField.name, "source"]}
                      className="flex-auto"
                      label="Source"
                    >
                      <Input
                        placeholder={"source:target"}
                        addonBefore={"Source"}
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
                  Add secret binding
                </Button>
              </div>
            );
          }}
        </Form.List>
        <Form.List name={"env"}>
          {(subFields, subOpt) => {
            return (
              <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
                <div className="text-lg font-semibold">
                  Enviornment variables
                </div>
                {subFields.map((subField) => (
                  <div key={subField.key} className="flex items-center gap-4">
                    <Form.Item
                      noStyle
                      name={[subField.name]}
                      className="flex-auto"
                    >
                      <Input placeholder={"FOO=BAR"} />
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
                  Add environment variable
                </Button>
              </div>
            );
          }}
        </Form.List>
        <Form.Item className="mt-4">
          <Button type="primary" htmlType="submit" loading={loading}>
            Launch
          </Button>
        </Form.Item>
      </Form>
    </>
  );
}

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

function ContainersTableCard({ api }: { api: AdminApi }) {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const { instanceId, hasToken, getAccessToken } = useContext(
    ProxyAuthTokenContext
  );
  const { openDrawer } = useContext(DataDrawerContext);
  const { data: containers, run: listContainers } = useRequest(
    async (): Promise<DockerContainer[] | undefined> => {
      return (await api.agentDockerContainerList({
        namespaceId: namespaceIdentifier,
        namespaceKind,
        resourceId: instanceId,
        xCryptocatProxyAuthorization: getAccessToken(),
      })) as DockerContainer[];
    },
    { manual: true }
  );

  const onInspect = useMemoizedFn((id: string) => {
    openDrawer({
      title: `Inspect ${id}`,
      onGetData: async () => {
        return await api.agentDockerContainerInspect({
          namespaceId: namespaceIdentifier,
          namespaceKind,
          resourceId: instanceId,
          xCryptocatProxyAuthorization: getAccessToken(),
          containerId: id,
        });
      },
    });
  });

  const { run: onStop } = useRequest(
    async (containerId: string): Promise<void> => {
      return await api.agentDockerContainerStop({
        namespaceId: namespaceIdentifier,
        namespaceKind,
        resourceId: instanceId,
        containerId,
        xCryptocatProxyAuthorization: getAccessToken(),
      });
    },
    { manual: true }
  );

  const { run: onRemove } = useRequest(
    async (containerId: string): Promise<void> => {
      return await api.agentDockerContainerRemove({
        namespaceId: namespaceIdentifier,
        namespaceKind,
        resourceId: instanceId,
        containerId,
        xCryptocatProxyAuthorization: getAccessToken(),
      });
    },
    { manual: true }
  );

  const columns = useDockerContainerColumns(onInspect, onStop, onRemove);
  return (
    <Card
      title="Containers"
      extra={
        <Button type="link" onClick={listContainers} disabled={!hasToken}>
          List containers
        </Button>
      }
    >
      <Table<DockerContainer>
        columns={columns}
        dataSource={containers}
        rowKey={(c) => c.Id}
      />
      <div>
        <Button
          type="link"
          size="small"
          onClick={() =>
            openDrawer({
              data: containers,
              title: "Containers",
            })
          }
        >
          View JSON
        </Button>
      </div>
    </Card>
  );
}

function RadiusConfigurationCard({ api }: { api: AdminApi }) {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
  const { run: pushAgentConfigRadius } = useRequest(
    async () => {
      if (instanceId) {
        const result = await api.pushAgentConfigRadius({
          namespaceKind,
          namespaceId: namespaceIdentifier,
          resourceId: instanceId,
          xCryptocatProxyAuthorization: getAccessToken(),
        });

        return result;
      }
    },
    {
      manual: true,
    }
  );
  return (
    <Card title="RADIUS configuration">
      <div>
        <Button type="primary" onClick={pushAgentConfigRadius}>
          Push configuration
        </Button>
      </div>
    </Card>
  );
}

function AgentDashboard({ api }: { api: AdminApi }) {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const { instanceId, hasToken, getAccessToken, setAccessToken } = useContext(
    ProxyAuthTokenContext
  );
  const { openDrawer } = useContext(DataDrawerContext);

  const { data } = useRequest(
    async () => {
      if (instanceId) {
        const result = await api.getAgentInstance({
          namespaceKind,
          namespaceId: namespaceIdentifier,
          resourceId: instanceId,
        });

        return result;
      }
    },
    {
      refreshDeps: [namespaceKind, namespaceIdentifier, instanceId],
    }
  );

  const { run: acquireToken } = useRequest(
    async () => {
      if (instanceId) {
        const result = await api.createAgentInstanceProxyAuthToken({
          namespaceId: namespaceIdentifier,
          namespaceKind,
          resourceId: instanceId,
        });
        setAccessToken(result.accessToken);
        return result;
      }
    },
    {
      manual: true,
    }
  );

  const getAgentDiagnostics = useMemoizedFn(async () => {
    return await api.getAgentDiagnostics({
      namespaceId: namespaceIdentifier,
      namespaceKind,
      resourceId: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });
  const getDockerInfo = useMemoizedFn(async () => {
    return await api.agentDockerInfo({
      namespaceId: namespaceIdentifier,
      namespaceKind,
      resourceId: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });

  const getDockerImages = useMemoizedFn(async () => {
    return await api.agentDockerImageList({
      namespaceId: namespaceIdentifier,
      namespaceKind,
      resourceId: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });

  const listDockerNetworks = useMemoizedFn(async () => {
    return await api.agentDockerNetworkList({
      namespaceId: namespaceIdentifier,
      namespaceKind,
      resourceId: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });

  return (
    <>
      <Card title="Agent proxy information">
        <JsonDataDisplay data={data} />
        <div className="mt-6">
          <Button type="primary" onClick={acquireToken}>
            Authorize
          </Button>
        </div>
      </Card>
      <ContainersTableCard api={api} />
      <Card title="View">
        <div className="flex flex-col gap-4 items-start">
          <Button
            type="link"
            onClick={() => {
              openDrawer({
                title: "Docker images",
                onGetData: getDockerImages,
              });
            }}
            disabled={!hasToken}
          >
            List Docker images
          </Button>
          <Button
            type="link"
            onClick={() => {
              openDrawer({
                title: "Docker networks",
                onGetData: listDockerNetworks,
              });
            }}
            disabled={!hasToken}
          >
            List Docker networks
          </Button>
          <Button
            type="link"
            onClick={() => {
              openDrawer({
                title: "Docker system info",
                onGetData: getDockerInfo,
              });
            }}
            disabled={!hasToken}
          >
            Docker system info
          </Button>
          <Button
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
          </Button>
        </div>
      </Card>

      {hasToken && (
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
      )}
    </>
  );
}

type DrawerHandler<T> = {
  title: string;
  data?: T;
  onGetData?: () => Promise<T>;
  toJson?: JsonDataDisplayProps<T>["toJson"];
};

type DrawerContextValue = {
  openDrawer: <T>(handler: DrawerHandler<T>) => void;
};

const DataDrawerContext = createContext<DrawerContextValue>({
  openDrawer: () => {},
});
function DrawerProvider(props: PropsWithChildren<{}>) {
  const [drawerOpen, { setTrue: setDrawerOpen, setFalse: closeDrawer }] =
    useBoolean(false);
  const [handler, setHandler] = useState<DrawerHandler<any>>();
  const onGetDataRef = useRef<() => Promise<any>>(() => Promise.resolve());
  const { data, loading, run } = useRequest(
    () => {
      return onGetDataRef.current();
    },
    { manual: true }
  );

  const openDrawer = useMemoizedFn((handler: DrawerHandler<any>) => {
    setHandler(handler);
    setDrawerOpen();
    if (handler.onGetData) {
      onGetDataRef.current = handler.onGetData ?? (() => Promise.resolve());
    }
    run();
  });

  const { title, toJson, data: handlerData, onGetData } = handler ?? {};
  return (
    <DataDrawerContext.Provider
      value={{
        openDrawer,
      }}
    >
      {props.children}
      <Drawer
        title={title}
        open={drawerOpen}
        placement="right"
        size="large"
        onClose={closeDrawer}
        extra={
          onGetData && (
            <Button onClick={run} type="primary">
              Refresh
            </Button>
          )
        }
      >
        <JsonDataDisplay<any>
          data={data ?? handlerData}
          toJson={toJson}
          loading={loading}
        />
      </Drawer>
    </DataDrawerContext.Provider>
  );
}

export default function AgentDashboardPage() {
  const { namespaceId: namespaceIdentifier, namespaceKind } =
    useContext(NamespaceContext);
  const { instanceId } = useParams<{ instanceId: string }>();

  const api = useAuthedClient(AdminApi);

  return (
    <>
      <Typography.Title>Agent Dashboard</Typography.Title>
      {instanceId && (
        <ProxyAuthTokenContextProvider instanceId={instanceId}>
          <DrawerProvider>
            <AgentDashboard api={api} />
          </DrawerProvider>
        </ProxyAuthTokenContextProvider>
      )}
    </>
  );
}
