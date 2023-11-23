import { useMemoizedFn, useRequest } from "ahooks";
import { Alert, Button, Form, Input, Select } from "antd";
import { DefaultOptionType } from "antd/es/select";
import Link from "antd/es/typography/Link";
import { useEffect, useMemo, useState } from "react";
import { useAppAuthContext } from "../auth/AuthProvider";
import {
  AdminApi,
  CertPolicy,
  Certificate,
  JsonWebSignatureAlgorithm,
  Key,
  KeyRef,
  NamespaceKind,
  ResourceKind,
} from "../generated";
import {
  base64StdEncodedToUrlEncoded,
  base64UrlDecodeBuffer,
  base64UrlEncodeBuffer,
  base64UrlEncodedToStdEncoded,
  toPEMBlock,
} from "../utils/encodingUtils";
import { useAuthedClient } from "../utils/useCertsApi";
import { useSystemAppRequest } from "./forms/useSystemAppRequest";

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

class EnrollmentSessionWithPkcs12Blob extends EnrollmentSession {
  constructor(
    readonly session: EnrollmentSession,
    public readonly pkcs12Blob: Blob
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
  const { account } = useAppAuthContext();

  const { data, run, loading } = useRequest(
    async (): Promise<EnrollmentSessionWithPemBlob | undefined> => {
      if (alg && account && certPolicy && certPolicy.keySpec.alg) {
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
              <DownloadPkcs12Section session={data} />
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

function DownloadPkcs12Section({ session }: { session: EnrollmentSession }) {
  const api = useAuthedClient(AdminApi);
  const [selectedKey, setSelectedKey] = useState<string>();

  const { data: wrapKeyWithId } = useRequest(
    async (): Promise<[Key, string] | undefined> => {
      if (!selectedKey) {
        return;
      }
      const [partitionKey, resourceId] = selectedKey.split("/");
      const [nsKind, nsId] = partitionKey.split(":");
      return [
        await api.getKey({
          namespaceId: nsId,
          namespaceKind: nsKind as NamespaceKind,
          resourceId,
        }),
        selectedKey,
      ];
    },
    {
      refreshDeps: [selectedKey],
    }
  );

  const [wrapKey, wrapKeyLocator] = wrapKeyWithId ?? [];

  const [password, setPassword] = useState<string>();

  const { data: pkcs12Session, run: exchangePkcs12 } = useRequest(
    async (
      session: EnrollmentSession,
      wrapKey: Key,
      keyLocator: string,
      password: string,
      legacy?: boolean
    ) => {
      const [payload, encKey] = await session.preparePayload(wrapKey, password);
      const resp = await api.exchangePKCS12({
        exchangePKCS12Request: {
          keyLocator,
          payload,
          legacy,
        },
        namespaceId: "me",
        namespaceKind: NamespaceKind.NamespaceKindUser,
        resourceId: session.enrollResponse.id,
      });
      const reqPayload = payload.split(".");
      const resPayload = resp.payload.split(".");
      if (reqPayload[1] !== resPayload[1]) {
        throw new Error("Invalid response payload[0]");
      }
      const iv = base64UrlDecodeBuffer(resPayload[2]);
      const ciphertext = base64UrlDecodeBuffer(resPayload[3]);
      const tag = base64UrlDecodeBuffer(resPayload[4]);
      const encrypted = new Uint8Array(ciphertext.byteLength + tag.byteLength);
      encrypted.set(new Uint8Array(ciphertext), 0);
      encrypted.set(new Uint8Array(tag), ciphertext.byteLength);
      const decrypted = await crypto.subtle.decrypt(
        {
          name: "AES-GCM",
          iv: iv,
          additionalData: new TextEncoder().encode(reqPayload[0]),
        },
        encKey,
        encrypted
      );
      return new EnrollmentSessionWithPkcs12Blob(
        session,
        new Blob([decrypted], {
          type: "application/x-pkcs12",
        })
      );
    },
    { manual: true }
  );

  const [pkcs12BlobUrl, setPkcs12BlobUrl] = useState<string>();
  useEffect(() => {
    if (!pkcs12Session) {
      return;
    }
    const url = URL.createObjectURL(pkcs12Session.pkcs12Blob);
    setPkcs12BlobUrl(url);
    return () => {
      URL.revokeObjectURL(url);
      setPkcs12BlobUrl(undefined);
    };
  }, [pkcs12Session]);
  return (
    <div className="ring-1 ring-neutral-400 rounded-md p-4 space-y-4">
      <div className="mb-4">
        <p>
          Use this option if you would like us to convert to PKCS#12 format.
          Your private key will be encrypted with the selected key and sent to
          our server to prepare for this certificate bundle.
        </p>
        <p>
          <span className="text-red-600 font-semibold text">
            Please delete this file after successful installation.
          </span>
        </p>
      </div>
      <KeySelector value={selectedKey} onChange={setSelectedKey} />
      <Form.Item label="Password">
        <Input
          type="password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
          }}
        />
      </Form.Item>
      <div className="flex flex-row gap-4 items-center">
        <Button
          type="primary"
          onClick={() => {
            if (wrapKey && wrapKeyLocator) {
              exchangePkcs12(session, wrapKey, wrapKeyLocator, password || "");
            }
          }}
        >
          Get link for certificate bundle (.p12) - Modern
        </Button>
        <span>
          Uses AES-256-CBC with PBKDF2 to encrypt the private key. Select legacy
          option if encounting problems with loading the certificate bundle.
        </span>
      </div>
      <div className="flex flex-row gap-4 items-center">
        <Button
          type="primary"
          onClick={() => {
            if (wrapKey && wrapKeyLocator) {
              exchangePkcs12(
                session,
                wrapKey,
                wrapKeyLocator,
                password || "",
                true
              );
            }
          }}
        >
          Get link for certificate bundle (.p12) - Legacy
        </Button>
        <span>
          Uses legacy algorithm RC2_CBC or 3DES_CBC to encrypt the private key.
          Works with macOS and iOS (requires Safari).
        </span>
      </div>
      {pkcs12BlobUrl && (
        <Link
          href={pkcs12BlobUrl}
          download={`${session.enrollResponse.id}.p12`}
        >
          Download .p12
        </Link>
      )}
    </div>
  );
}

function KeySelector({
  value,
  onChange,
}: {
  value?: string;
  onChange?: (value: string) => void;
}) {
  const { data: systemApp } = useSystemAppRequest("backend");

  const nsId = systemApp?.servicePrincipalId;

  const api = useAuthedClient(AdminApi);
  const { data: keyPolicies, loading: keyPoliciesLoading } = useRequest(
    async () => {
      return await api.listKeyPolicies({
        namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
        namespaceId: nsId!,
      });
    },
    {
      refreshDeps: [nsId],
      ready: !!nsId,
    }
  );

  const keyPoliciesOptions = useMemo(() => {
    return keyPolicies?.map(
      (p): DefaultOptionType => ({
        label: (
          <span>
            {p.displayName} (<span className="font-mono">{p.id}</span>)
          </span>
        ),
        value: p.id,
      })
    );
  }, [keyPolicies]);

  const [selectedPolicyId, setSelectedPolicyId] = useState<string>();

  const { data: keysWithPolicyID, loading: keysLoading } = useRequest(
    async (): Promise<[KeyRef[], string]> => {
      return [
        await api.listKeys({
          namespaceKind: NamespaceKind.NamespaceKindServicePrincipal,
          namespaceId: nsId!,
          policyId: selectedPolicyId!,
        }),
        selectedPolicyId!,
      ];
    },
    {
      refreshDeps: [nsId, selectedPolicyId],
      ready: !!nsId && !!selectedPolicyId,
    }
  );

  const keyOptions = useMemo(() => {
    const [keys, policyId] = keysWithPolicyID ?? [];
    if (!policyId || policyId !== selectedPolicyId) {
      return undefined;
    }
    return keys?.map(
      (k): DefaultOptionType => ({
        label: <span className="font-mono">{k.id}</span>,
        value: k.id,
      })
    );
  }, [keysWithPolicyID, selectedPolicyId]);

  const selectedKeyId = value?.split("/")[1];

  const onSelectChange = useMemoizedFn((keyId) => {
    onChange?.(
      NamespaceKind.NamespaceKindServicePrincipal +
        ":" +
        nsId +
        ":" +
        ResourceKind.ResourceKindKey +
        "/" +
        keyId
    );
  });

  return (
    <div>
      <Form.Item label="Select encryption key policy">
        <Select
          options={keyPoliciesOptions}
          loading={keyPoliciesLoading}
          value={selectedPolicyId}
          onChange={setSelectedPolicyId}
        />
      </Form.Item>
      <Form.Item label="Select encryption key">
        <Select
          disabled={!selectedPolicyId}
          options={keyOptions}
          loading={keysLoading}
          value={selectedKeyId}
          onChange={onSelectChange}
        />
      </Form.Item>
    </div>
  );
}
