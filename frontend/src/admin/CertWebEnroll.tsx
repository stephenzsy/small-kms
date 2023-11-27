import { useRequest } from "ahooks";
import { Alert, Button } from "antd";
import { useEffect, useMemo, useState } from "react";
import {
  CertPolicy,
  Certificate,
  JsonWebSignatureAlgorithm,
  Key,
} from "../generated";
import {
  base64StdEncodedToUrlEncoded,
  base64UrlEncodeBuffer,
  base64UrlEncodedToStdEncoded,
  toPEMBlock,
} from "../utils/encodingUtils";

function useKeyGenAlgParams(certPolicy: CertPolicy | undefined) {
  return useMemo((): RsaHashedKeyGenParams | EcKeyGenParams | undefined => {
    if (!certPolicy) {
      return undefined;
    }
    switch (certPolicy.keySpec.alg) {
      case JsonWebSignatureAlgorithm.Rs256:
        return {
          name: "RSASSA-PKCS1-v1_5",
          hash: "SHA-256",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebSignatureAlgorithm.Rs384:
        return {
          name: "RSASSA-PKCS1-v1_5",
          hash: "SHA-384",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebSignatureAlgorithm.Rs512:
        return {
          name: "RSASSA-PKCS1-v1_5",
          hash: "SHA-512",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebSignatureAlgorithm.Es256:
        return {
          name: "ECDSA",
          namedCurve: "P-256",
        };
      case JsonWebSignatureAlgorithm.Es384:
        return {
          name: "ECDSA",
          namedCurve: "P-384",
        };
      case JsonWebSignatureAlgorithm.Es512:
        return {
          name: "ECDSA",
          namedCurve: "P-521",
        };

      case JsonWebSignatureAlgorithm.Ps256:
        return {
          name: "RSA-PSS",
          hash: "SHA-256",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebSignatureAlgorithm.Ps384:
        return {
          name: "RSA-PSS",
          hash: "SHA-384",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebSignatureAlgorithm.Ps512:
        return {
          name: "RSA-PSS",
          hash: "SHA-512",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
    }
    return undefined;
  }, [certPolicy]);
}

class EnrollmentSession {
  constructor(
    readonly keypair: CryptoKeyPair,
    readonly proofJwt: string,
    readonly enrollResponse: Certificate
  ) {}

  public async getPemBlob(): Promise<Blob> {
    const privateKeyDer = await window.crypto.subtle.exportKey(
      "pkcs8",
      this.keypair.privateKey
    );

    return new Blob(
      [
        [
          toPEMBlock(
            btoa(String.fromCharCode(...new Uint8Array(privateKeyDer))),
            "PRIVATE KEY"
          ),
          ...(this.enrollResponse.jwk.x5c ?? []).map((cert) =>
            toPEMBlock(base64UrlEncodedToStdEncoded(cert), "CERTIFICATE")
          ),
        ].join("\n"),
      ],
      {
        type: "application/x-pem-file",
      }
    );
  }

  public async preparePayload(
    wrapKey: Key,
    password: string
  ): Promise<[string, CryptoKey]> {
    const encKey = await crypto.subtle.generateKey(
      {
        name: "AES-GCM",
        length: 256,
      },
      true,
      ["encrypt", "decrypt"]
    );

    const wrapPublickey = await crypto.subtle.importKey(
      "jwk",
      wrapKey,
      { name: "RSA-OAEP", hash: "SHA-256" },
      true,
      ["wrapKey"]
    );

    const joseHeader = {
      alg: "RSA-OAEP-256",
      kid: wrapKey.kid,
      enc: "A256GCM",
    };

    const joseEncoded = base64StdEncodedToUrlEncoded(
      btoa(JSON.stringify(joseHeader))
    );
    const encryptedKey = await crypto.subtle.wrapKey(
      "raw",
      encKey,
      wrapPublickey,
      {
        name: "RSA-OAEP",
      }
    );
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const plainTextBuffer = new TextEncoder().encode(
      JSON.stringify({
        password,
        privateKey: await crypto.subtle.exportKey(
          "jwk",
          this.keypair.privateKey
        ),
      })
    );
    const encrypted = await crypto.subtle.encrypt(
      {
        name: "AES-GCM",
        iv,
        additionalData: new TextEncoder().encode(joseEncoded),
      },
      encKey,
      plainTextBuffer
    );
    const ciphertext = encrypted.slice(0, -16);
    const tag = encrypted.slice(-16);
    return [
      [
        joseEncoded,
        base64UrlEncodeBuffer(encryptedKey),
        base64UrlEncodeBuffer(iv),
        base64UrlEncodeBuffer(ciphertext),
        base64UrlEncodeBuffer(tag),
      ].join("."),
      encKey,
    ];
  }
}

class EnrollmentSessionWithPemBlob extends EnrollmentSession {
  constructor(
    readonly session: EnrollmentSession,
    public readonly pemBlob: Blob
  ) {
    super(session.keypair, session.proofJwt, session.enrollResponse);
  }
}

export function CertWebEnroll({
  certPolicy,
}: {
  certPolicy: CertPolicy | undefined;
}) {
  const alg = useKeyGenAlgParams(certPolicy);

  const { data, run, loading } = useRequest(
    async (): Promise<EnrollmentSessionWithPemBlob | undefined> => {
      if (alg && certPolicy && certPolicy.keySpec.alg) {
        // const enrollResp = await api.enrollCertificate({
        //   enrollCertificateRequest: {
        //     enrollmentType: "group-memeber",
        //     proof: proofJwt,
        //     publicKey: (await window.crypto.subtle.exportKey(
        //       "jwk",
        //       keypair.publicKey
        //     )) as JsonWebSignatureKey,
        //   },
        //   namespaceId,
        //   namespaceKind,
        //   resourceId: certPolicy.id,
        // });
        // const session = new EnrollmentSession(keypair, proofJwt, enrollResp);
        // return new EnrollmentSessionWithPemBlob(
        //   session,
        //   await session.getPemBlob()
        // );
        return undefined;
      }
    },
    { manual: true }
  );

  const [blobUrl, setBlobUrl] = useState<string>();

  useEffect(() => {
    if (data?.pemBlob) {
      const url = URL.createObjectURL(data.pemBlob);
      setBlobUrl(url);
      return () => {
        URL.revokeObjectURL(url);
        setBlobUrl(undefined);
      };
    }
  }, [data?.pemBlob]);

  return (
    <div className="space-y-4">
      <Alert
        type="info"
        message="Private key does not leave your browser unencrypted"
      />
      {data ? (
        <div className="ring-1 ring-neutral-400 rounded-md p-4 space-y-4">
          <div className="text-lg font-semibold">Download certificate</div>
          <Alert
            type="warning"
            message="This is the only time you can download the certificate with the private key. Please save it in a safe place."
          />
          {blobUrl && (
            <>
              <div>
                <Button
                  type="primary"
                  href={blobUrl}
                  download={`${data.enrollResponse.id}.pem`}
                >
                  Download certificate (.pem)
                </Button>
              </div>
              <div className="ring-1 ring-neutral-400 rounded-md p-4">
                You can use the following command to convert the certificate
                bundle to PKCS#12 (.p12) format to be installed to Windows
                Certificate Manager or macOS Keychain.
                <pre className="bg-neutral-200 p-2 rounded-md">
                  openssl pkcs12 -export -inkey{" "}
                  {`${data.enrollResponse.id}.pem`} -in{" "}
                  {`${data.enrollResponse.id}.pem`} -out{" "}
                  {`${data.enrollResponse.id}.p12`}
                </pre>
              </div>
            </>
          )}
        </div>
      ) : (
        <Button onClick={run} type="primary" loading={loading}>
          Begin enroll
        </Button>
      )}
    </div>
  );
}

// function KeySelector({
//   value,
//   onChange,
// }: {
//   value?: string;
//   onChange?: (value: string) => void;
// }) {
//   const { data: systemApp } = useSystemAppRequest("backend");

//   const nsId = systemApp?.servicePrincipalId;

//   const api = useAuthedClient(AdminApi);
//   const { data: keyPolicies, loading: keyPoliciesLoading } = useRequest(
//     async () => {
//       return await api.listKeyPolicies({
//         namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
//         namespaceId: nsId!,
//       });
//     },
//     {
//       refreshDeps: [nsId],
//       ready: !!nsId,
//     }
//   );

//   const keyPoliciesOptions = useMemo(() => {
//     return keyPolicies?.map(
//       (p): DefaultOptionType => ({
//         label: (
//           <span>
//             {p.displayName} (<span className="font-mono">{p.id}</span>)
//           </span>
//         ),
//         value: p.id,
//       })
//     );
//   }, [keyPolicies]);

//   const [selectedPolicyId, setSelectedPolicyId] = useState<string>();

//   const { data: keysWithPolicyID, loading: keysLoading } = useRequest(
//     async (): Promise<[KeyRef[], string]> => {
//       return [
//         await api.listKeys({
//           namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
//           namespaceId: nsId!,
//           policyId: selectedPolicyId!,
//         }),
//         selectedPolicyId!,
//       ];
//     },
//     {
//       refreshDeps: [nsId, selectedPolicyId],
//       ready: !!nsId && !!selectedPolicyId,
//     }
//   );

//   const keyOptions = useMemo(() => {
//     const [keys, policyId] = keysWithPolicyID ?? [];
//     if (!policyId || policyId !== selectedPolicyId) {
//       return undefined;
//     }
//     return keys?.map(
//       (k): DefaultOptionType => ({
//         label: <span className="font-mono">{k.id}</span>,
//         value: k.id,
//       })
//     );
//   }, [keysWithPolicyID, selectedPolicyId]);

//   const selectedKeyId = value?.split("/")[1];

//   const onSelectChange = useMemoizedFn((keyId) => {
//     onChange?.(
//       NamespaceKind.NamespaceKindServicePrincipal +
//         ":" +
//         nsId +
//         ":" +
//         ResourceKind.ResourceKindKey +
//         "/" +
//         keyId
//     );
//   });

//   return (
//     <div>
//       <Form.Item label="Select encryption key policy">
//         <Select
//           options={keyPoliciesOptions}
//           loading={keyPoliciesLoading}
//           value={selectedPolicyId}
//           onChange={setSelectedPolicyId}
//         />
//       </Form.Item>
//       <Form.Item label="Select encryption key">
//         <Select
//           disabled={!selectedPolicyId}
//           options={keyOptions}
//           loading={keysLoading}
//           value={selectedKeyId}
//           onChange={onSelectChange}
//         />
//       </Form.Item>
//     </div>
//   );
// }
