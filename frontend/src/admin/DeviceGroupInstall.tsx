import { useMemoizedFn } from "ahooks";
import React from "react";
import { Button } from "../components/Button";
import { AdminApi } from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { InputField } from "./InputField";

export function DeviceGroupInstall({ namespaceId }: { namespaceId: string }) {
  const adminApi = useAuthedClient(AdminApi);

  const [groupId, setGroupId] = React.useState("");
  const [templateId, setTemplateId] = React.useState("");
  const [windowsScripts, setWindowsScripts] = React.useState<string[]>([]);

  const onGenerateWindowsScript = useMemoizedFn(() => {
    if (groupId && templateId) {
      /*
      const req: CertificateEnrollmentRequest = {
        appId: linkInfo.applicationClientId,
        deviceNamespaceId: namespaceId,
        linkId: linkInfo.ref.id,
        servicePrincipalId: linkInfo.servicePrincipalId,
        type: "device-linked-service-principal",
      };*/
      setWindowsScripts([
        [
          // `$env:AZURE_CLIENT_ID = '${linkInfo.applicationClientId}}'`,
          // "$env:AZURE_TENANT_ID = '...'",
          // `$exe v2 certificate-templates enroll post \\`,
          // `    --namespace-id ${groupId} \\`,
          // `    --template-id ${templateId} \\`,
          // `    --body '${JSON.stringify(req)} > .\\enroll.json'`,
        ].join("\n"),
        [
          "$jwtClaims = (Get-Content -Raw .\\enroll.json |",
          "    ConvertFrom-Json |",
          "    select -ExpandProperty jwtClaims);",
          "$stdEncodedjwtClaims = $jwtClaims.Replace('-','+').Replace('_','/');",
          "if ($stdEncodedjwtClaims.Length % 4 -eq 2) { $stdEncodedjwtClaims += '==' }",
          "if ($stdEncodedjwtClaims.Length % 4 -eq 3) { $stdEncodedjwtClaims += '=' }",
          "[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String( $stdEncodedjwtClaims)) | ConvertFrom-Json | ConvertTo-Json",
        ].join("\n"),
        [
          '$header = \'{"typ":"JWT","alg":"RS256"}\'',
          "$headerEncoded = [System.Convert]::ToBase64String($header.ToCharArray()).Replace('+','-').Replace('/','_').Replace('=','');",
          '"$headerEncoded.$jwtClaims" > .\\signingjwt.txt;',
          "Start-Process pwsh.exe -Verb RunAs",
        ].join("\n"),
        [
          '$cert = (New-SelfSignedCertificate -Subject "CN=TMP" -KeyExportPolicy NonExportable -KeyAlgorithm RSA -KeyLength 2048);',
          "$privateKey = [System.Security.Cryptography.X509Certificates.RSACertificateExtensions]::GetRSAPrivateKey($cert);",
          "$signatureBytes = $privateKey.SignData($data, [System.Security.Cryptography.HashAlgorithmName]::SHA256, [System.Security.Cryptography.RSASignaturePadding]::Pkcs1)",
          '$signatureEncoded = [System.Convert]::ToBase64String($signatureBytes).Replace("+","-").Replace("/","_").Replace("=","");',
          '"$dataRaw.$signatureEncoded" > .\\signedJwt.txt;',
        ].join("\n"),
        [
          "Get-ChildItem Microsoft.PowerShell.Security\\Certificate::LocalMachine\\MY\\...",
        ].join("\n"),
      ]);
    }
  });
  return (
    <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
      <div className="pad-4 space-y-4">
        <InputField
          labelContent="Group ID for template"
          value={groupId}
          onChange={setGroupId}
        />
        <InputField
          labelContent="Template ID"
          value={templateId}
          onChange={setTemplateId}
        />
      </div>
      <Button
        className="mt-4"
        variant="primary"
        onClick={onGenerateWindowsScript}
      >
        Generate windows script
      </Button>
      {windowsScripts && (
        <div>
          <div className="ring-1 ring-neutral-500 max-w-full overflow-auto p-4">
            <div>Step 1. Request a new certificate</div>
            <pre className="max-w-full overflow-auto">{windowsScripts[0]}</pre>
          </div>
          <div className="ring-1 ring-neutral-500 max-w-full overflow-auto p-4">
            <div>Step 2. Review certificate details</div>
            <pre className="max-w-full overflow-auto">{windowsScripts[1]}</pre>
          </div>
          <div className="ring-1 ring-neutral-500 max-w-full overflow-auto p-4">
            <div>
              Step 3. If Everything looks ok, create header section and launch
              elevated window
            </div>
            <pre className="max-w-full overflow-auto">{windowsScripts[2]}</pre>
          </div>
          <div className="ring-1 ring-neutral-500 max-w-full overflow-auto p-4">
            <div>
              Step 4. In elevated shell, create private key and sign the request
            </div>
            <pre className="max-w-full overflow-auto">{windowsScripts[3]}</pre>
          </div>
          <div className="ring-1 ring-neutral-500 max-w-full overflow-auto p-4">
            <div>
              Step 4. Come back and copy the thumbprint from the certificate
            </div>
            <pre className="max-w-full overflow-auto">{windowsScripts[4]}</pre>
          </div>
        </div>
      )}
    </section>
  );
}
