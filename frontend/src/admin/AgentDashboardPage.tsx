import {
  useBoolean,
  useLatest,
  useMemoizedFn,
  useRequest,
  useUpdate,
} from "ahooks";
import { Button, Card, Drawer, Form, Input, Table, Typography } from "antd";
import { useParams } from "react-router-dom";
import {
  AdminApi,
  AgentMode,
  LaunchAgentRequest,
  NamespaceKind,
  PullImageRequest,
  SecretMount,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import {
  MutableRefObject,
  PropsWithChildren,
  createContext,
  useContext,
  useRef,
  useState,
} from "react";
import {
  JsonDataDisplay,
  JsonDataDisplayProps,
} from "../components/JsonDataDisplay";
import { useForm } from "antd/es/form/Form";
import { ColumnType } from "antd/es/table";
import {
  ProxyAuthTokenContext,
  ProxyAuthTokenContextProvider,
} from "./ProxyAuthTokenContext";
import { XMarkIcon } from "@heroicons/react/24/outline";

type DockerPullImageFormState = {
  imageTag: string;
};

function DockerPullImageForm() {
  const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
  const [form] = useForm<DockerPullImageFormState>();
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const api = useAuthedClient(AdminApi);

  const { run: pullImage } = useRequest(
    async (req: PullImageRequest) => {
      await api.agentPullImage({
        namespaceIdentifier,
        namespaceKind,
        resourceIdentifier: instanceId,
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
          imageTag: s.imageTag,
        });
      }}
    >
      <Form.Item<DockerPullImageFormState>
        name="imageTag"
        label="Image tag"
        required
      >
        <Input placeholder="latest" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
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
  networkName?: string;
  secrets: SecretMount[];
};

function LaunchContainerForm({ mode }: { mode: AgentMode }) {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const [form] = useForm<LaunchContainerFormState>();
  const api = useAuthedClient(AdminApi);
  const { instanceId, getAccessToken } = useContext(ProxyAuthTokenContext);
  const { run } = useRequest(
    async (req: LaunchAgentRequest) => {
      await api.agentLaunchAgent({
        namespaceIdentifier,
        namespaceKind,
        resourceIdentifier: instanceId,
        launchAgentRequest: req,
        xCryptocatProxyAuthorization: getAccessToken(),
      });
    },
    { manual: true }
  );

  const onFinish = useMemoizedFn((values: LaunchContainerFormState) => {
    if (!values.imageTag || !values.containerName || !values.listenerAddress) {
      return;
    }
    run({
      containerName: values.containerName,
      mode,
      exposedPortSpecs: values.exposedPortSpecs,
      hostBinds: values.hostBinds,
      imageTag: values.imageTag,
      listenerAddress: values.listenerAddress,
      networkName: values.networkName,
      secrets: values.secrets,
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
          <Input placeholder=":10443" />
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
                  Add hsot binding
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
        <Form.Item className="mt-4">
          <Button type="primary" htmlType="submit">
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
  Ports: {
    IP: string;
    PrivatePort: number;
    PublicPort: number;
    Type: string;
  }[];
  Names: string[];
};

function useDockerContainerColumns(
  onInspect: (id: string) => void
): ColumnType<DockerContainer>[] {
  return [
    {
      title: "Id",
      dataIndex: "Id",
      key: "Id",
      render: (id: DockerContainer["Id"]) => {
        return (
          <span className="font-mono">
            {id.substring(0, 12)}{" "}
            <Button
              type="link"
              onClick={() => {
                onInspect(id);
              }}
            >
              Inspect
            </Button>
          </span>
        );
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
  ];
}

function ContainersTableCard({ api }: { api: AdminApi }) {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const { instanceId, hasToken, getAccessToken } = useContext(
    ProxyAuthTokenContext
  );
  const { openDrawer } = useContext(DataDrawerContext);
  const { data: containers, run: listContainers } = useRequest(
    async (): Promise<DockerContainer[] | undefined> => {
      return (await api.agentDockerContainerList({
        namespaceIdentifier,
        namespaceKind,
        resourceIdentifier: instanceId,
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
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
          xCryptocatProxyAuthorization: getAccessToken(),
          containerId: id,
        });
      },
    });
  });

  const columns = useDockerContainerColumns(onInspect);
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

function AgentDashboard({ api }: { api: AdminApi }) {
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
  const { instanceId, hasToken, getAccessToken, setAccessToken } = useContext(
    ProxyAuthTokenContext
  );
  const { openDrawer } = useContext(DataDrawerContext);

  const { data } = useRequest(
    async () => {
      if (instanceId) {
        const result = await api.getAgentInstance({
          namespaceKind,
          namespaceIdentifier,
          resourceIdentifier: instanceId,
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
          namespaceIdentifier,
          namespaceKind,
          resourceIdentifier: instanceId,
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
      namespaceIdentifier,
      namespaceKind,
      resourceIdentifier: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });
  const getDockerInfo = useMemoizedFn(async () => {
    return await api.agentDockerInfo({
      namespaceIdentifier,
      namespaceKind,
      resourceIdentifier: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });

  const getDockerImages = useMemoizedFn(async () => {
    return await api.agentDockerImageList({
      namespaceIdentifier,
      namespaceKind,
      resourceIdentifier: instanceId,
      xCryptocatProxyAuthorization: getAccessToken(),
    });
  });

  const listDockerNetworks = useMemoizedFn(async () => {
    return await api.agentDockerNetworkList({
      namespaceIdentifier,
      namespaceKind,
      resourceIdentifier: instanceId,
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
  const { namespaceIdentifier, namespaceKind } = useContext(NamespaceContext);
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
