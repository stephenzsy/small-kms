import { useForm } from "antd/es/form/Form";
import {
  AgentConfigRadiusFields,
  RadiusServerConfig,
  RadiusServerListenConfig,
  RadiusServerListenerType,
} from "../../generated";
import { useRadiusConfigPatch } from "../contexts/RadiusConfigPatchContext";
import { useEffect } from "react";
import { Button, Form, Input, Radio } from "antd";

export function RadiusConfigServersForm() {
  const { run, data } = useRadiusConfigPatch();
  const [form] = useForm<Pick<AgentConfigRadiusFields, "servers">>();

  useEffect(() => {
    if (data && data.servers) {
      form.setFieldsValue({
        servers: data.servers,
      });
    }
  }, [data, form]);
  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(v) => {
        run({
          servers: v.servers,
        });
      }}
    >
      <Form.List name={"servers"}>
        {(subFields, subOpt) => {
          return (
            <div className="flex flex-col gap-4">
              {subFields.map((subField) => (
                <div
                  key={subField.key}
                  className=" ring-1 ring-neutral-400 p-4 rounded-md"
                >
                  <Form.Item<RadiusServerConfig[]>
                    name={[subField.name, "name"]}
                    label="Name"
                    required
                  >
                    <Input placeholder={"default"} />
                  </Form.Item>
                  <Form.List name={[subField.name, "listeners"]}>
                    {(subFields, subOpt) => {
                      return (
                        <div className="flex flex-col gap-4">
                          <div className="text-lg font-medium">Listeners</div>
                          {subFields.map((subField) => (
                            <div
                              key={subField.key}
                              className=" ring-1 ring-neutral-400 p-4 rounded-md"
                            >
                              <Form.Item<RadiusServerListenConfig[]>
                                name={[subField.name, "type"]}
                                label="Type"
                                required
                              >
                                <Radio.Group
                                  options={[
                                    {
                                      label: "Auth",
                                      value:
                                        RadiusServerListenerType.RadiusServerListenerTypeAuth,
                                    },
                                    {
                                      label: "Accounting",
                                      value:
                                        RadiusServerListenerType.RadiusServerListenerTypeAcct,
                                    },
                                  ]}
                                />
                              </Form.Item>
                              <Form.Item<RadiusServerListenConfig[]>
                                name={[subField.name, "ipaddr"]}
                                label="IP Address"
                                required
                              >
                                <Input placeholder={"*"} />
                              </Form.Item>

                              <Button
                                danger
                                onClick={() => {
                                  subOpt.remove(subField.name);
                                }}
                              >
                                Remove
                              </Button>
                            </div>
                          ))}
                          <Button
                            type="dashed"
                            onClick={() => subOpt.add()}
                            block
                          >
                            Add listener
                          </Button>
                        </div>
                      );
                    }}
                  </Form.List>
                  <div className="mt-4">
                    <Button
                      danger
                      onClick={() => {
                        subOpt.remove(subField.name);
                      }}
                    >
                      Remove
                    </Button>
                  </div>
                </div>
              ))}
              <Button type="dashed" onClick={() => subOpt.add()} block>
                Add server configuration
              </Button>
            </div>
          );
        }}
      </Form.List>
      <Form.Item className="mt-6">
        <Button htmlType="submit" type="primary">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}
