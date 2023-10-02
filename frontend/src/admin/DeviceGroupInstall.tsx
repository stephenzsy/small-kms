import { useMemoizedFn, useRequest } from "ahooks";
import { Button } from "../components/Button";
import { AdminApi, ServicePrincipalLinkedDevice } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { InputField } from "./InputField";
import React from "react";

export function DeviceGroupInstall(props: {
  namespaceId: string;
  linkInfo?: ServicePrincipalLinkedDevice;
}) {
  const adminApi = useAuthedClient(AdminApi);

  const [groupId, setGroupId] = React.useState("");
  const [windowsScript, setWindowsScript] = React.useState("");

  const onGenerateWindowsScript = useMemoizedFn(() => {
    setWindowsScript(
      [
        "$groupId = '" + groupId + "'",
        "$namespaceId = '" + props.namespaceId + "'",
        "",
      ].join("\n")
    );
  });
  return (
    <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
      <InputField
        labelContent="Group ID for template"
        value={groupId}
        onChange={setGroupId}
      />
      <Button
        className="mt-4"
        variant="primary"
        onClick={onGenerateWindowsScript}
      >
        Generate windows script
      </Button>
      <pre>{JSON.stringify(props.linkInfo, undefined, 2)}</pre>
      <output>
        <pre>{windowsScript}</pre>
      </output>
    </section>
  );
}
