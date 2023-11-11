import { XMarkIcon } from "@heroicons/react/24/outline";
import { Button, Form, Input } from "antd";
import { useForm } from "antd/es/form/Form";
import { useEffect } from "react";
import { AgentContainerConfiguration } from "../../generated";
import { AgentServerConfigFormState } from "../RadiusConfigPage";
import { useRadiusConfigPatch } from "../contexts/RadiusConfigPatchContext";

export function RadiusConfigContainerForm() {
  const [form] = useForm<AgentServerConfigFormState>();
  const { run, data } = useRadiusConfigPatch();
  useEffect(() => {
    if (data?.container) {
      form.setFieldsValue(data.container);
    }
  }, [data?.container]);

  return (
    <Form<AgentContainerConfiguration>
      form={form}
      layout="vertical"
      onFinish={(values) => {
        run({
          container: values,
        });
      }}
    >
      <Form.Item<AgentServerConfigFormState>
        name="containerName"
        label="Azure Container container name"
        required
      >
        <Input placeholder="radius" />
      </Form.Item>

      <Form.Item<AgentServerConfigFormState>
        name="imageRepo"
        label="Image repository"
        required
      >
        <Input placeholder="example.com/radius" />
      </Form.Item>

      <Form.Item<AgentServerConfigFormState>
        name="imageTag"
        label="Image tag"
        required
      >
        <Input placeholder="latest" />
      </Form.Item>

      <Form.List name={"exposedPortSpecs"}>
        {(subFields, subOpt) => {
          return (
            <div className="flex flex-col gap-4 ring-1 ring-neutral-400 p-4 rounded-md mt-6">
              <div className="text-lg font-semibold">Port bindings</div>
              {subFields.map((subField) => (
                <div key={subField.key} className="flex items-center gap-4">
                  <Form.Item
                    noStyle
                    name={[subField.name]}
                    className="flex-auto"
                  >
                    <Input placeholder={"1812/udp"} />
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
                Add port binding
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
                  <Form.Item noStyle name={subField.name} className="flex-auto">
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
                    <Input placeholder={"source:target"} addonBefore={"Name"} />
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
              <div className="text-lg font-semibold">Enviornment variables</div>
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
        <Button htmlType="submit" type="primary">
          Submit
        </Button>
      </Form.Item>
    </Form>
  );
}
