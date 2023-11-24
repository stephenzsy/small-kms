import { useRequest, useUnmountedRef } from "ahooks";
import { Button, Checkbox } from "antd";
import { useEffect, useState } from "react";
import { JsonWebSignatureAlgorithm } from "../../generated";
import {
  AdminApi,
  Certificate,
  CertificatePolicy,
  ExchangePKCS12Result,
  JsonWebKey,
  NamespaceProvider,
} from "../../generated/apiv2";
import { useAdminApi, useAuthedClientV2 } from "../../utils/useCertsApi";
import { useNamespace } from "../contexts/NamespaceContextRouteProvider";
import {
  base64StdEncodedToUrlEncoded,
  base64UrlDecodeBuffer,
  base64UrlEncodeBuffer,
  base64UrlEncodedToStdEncoded,
  toPEMBlock,
} from "../../utils/encodingUtils";

class EnrollmentSession {
  constructor(
    public readonly policyNamespaceProvider: NamespaceProvider,
    public readonly policyNamespaceID: string,
    public readonly certPolicy: CertificatePolicy,
    public readonly keyPair?: CryptoKeyPair,
    public readonly enrollCertResponse?: Certificate
  ) {}

  public async getPemBlob(): Promise<Blob | undefined> {
    if (!this.keyPair || !this.enrollCertResponse?.jwk) {
      return undefined;
    }
    const privateKeyDer = await window.crypto.subtle.exportKey(
      "pkcs8",
      this.keyPair.privateKey
    );

    return new Blob(
      [
        [
          toPEMBlock(
            btoa(String.fromCharCode(...new Uint8Array(privateKeyDer))),
            "PRIVATE KEY"
          ),
          ...(this.enrollCertResponse.jwk.x5c ?? []).map((cert) =>
            toPEMBlock(base64UrlEncodedToStdEncoded(cert), "CERTIFICATE")
          ),
        ].join("\n"),
      ],
      {
        type: "application/x-pem-file",
      }
    );
  }

  private get keyAlgorithm(): string {
    switch (this.certPolicy.keySpec.alg) {
      case JsonWebSignatureAlgorithm.Rs256:
      case JsonWebSignatureAlgorithm.Rs384:
      case JsonWebSignatureAlgorithm.Rs512:
        return "RSASSA-PKCS1-v1_5";
      case JsonWebSignatureAlgorithm.Es256:
      case JsonWebSignatureAlgorithm.Es384:
      case JsonWebSignatureAlgorithm.Es512:
        return "ECDSA";
      case JsonWebSignatureAlgorithm.Ps256:
      case JsonWebSignatureAlgorithm.Ps384:
      case JsonWebSignatureAlgorithm.Ps512:
        return "RSA-PSS";
      default:
        throw new Error("Unsupported algorithm");
    }
  }

  private get signAlgorithm():
    | AlgorithmIdentifier
    | RsaPssParams
    | EcdsaParams {
    switch (this.certPolicy.keySpec.alg) {
      case JsonWebSignatureAlgorithm.Rs256:
      case JsonWebSignatureAlgorithm.Rs384:
      case JsonWebSignatureAlgorithm.Rs512:
        return {
          name: "RSASSA-PKCS1-v1_5",
        };
      case JsonWebSignatureAlgorithm.Es256:
        return {
          name: "ECDSA",
          hash: "SHA-256",
        };
      case JsonWebSignatureAlgorithm.Es384:
        return {
          name: "ECDSA",
          hash: "SHA-384",
        };
      case JsonWebSignatureAlgorithm.Es512:
        return {
          name: "ECDSA",
          hash: "SHA-512",
        };
      case JsonWebSignatureAlgorithm.Ps256:
      case JsonWebSignatureAlgorithm.Ps384:
      case JsonWebSignatureAlgorithm.Ps512:
        return {
          name: "RSA-PSS",
        };
      default:
        throw new Error("Unsupported algorithm");
    }
  }

  private get hashAlgorithm(): AlgorithmIdentifier {
    switch (this.certPolicy.keySpec.alg) {
      case JsonWebSignatureAlgorithm.Rs256:
      case JsonWebSignatureAlgorithm.Es256:
      case JsonWebSignatureAlgorithm.Ps256:
        return "SHA-256";
      case JsonWebSignatureAlgorithm.Rs384:
      case JsonWebSignatureAlgorithm.Es384:
      case JsonWebSignatureAlgorithm.Ps384:
        return "SHA-384";
      case JsonWebSignatureAlgorithm.Rs512:
      case JsonWebSignatureAlgorithm.Es512:
      case JsonWebSignatureAlgorithm.Ps512:
        return "SHA-512";
      default:
        throw new Error("Unsupported algorithm");
    }
  }

  private get keyGenParams(): RsaHashedKeyGenParams | EcKeyGenParams {
    switch (this.certPolicy.keySpec.alg) {
      case JsonWebSignatureAlgorithm.Rs256:
      case JsonWebSignatureAlgorithm.Rs384:
      case JsonWebSignatureAlgorithm.Rs512:
      case JsonWebSignatureAlgorithm.Ps256:
      case JsonWebSignatureAlgorithm.Ps384:
      case JsonWebSignatureAlgorithm.Ps512:
        return {
          name: this.keyAlgorithm,
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: this.hashAlgorithm,
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
      default:
        throw new Error("Unsupported algorithm");
    }
  }

  public async withGeneratedKeyPair() {
    if (this.keyPair) {
      return this;
    }
    const keyPair = await crypto.subtle.generateKey(this.keyGenParams, true, [
      "sign",
      "verify",
    ]);
    return new EnrollmentSession(
      this.policyNamespaceProvider,
      this.policyNamespaceID,
      this.certPolicy,
      keyPair
    );
  }

  public async withEnrollmentResponse(api: AdminApi) {
    if (!this.keyPair) {
      return this;
    }
    const resp = await api.enrollCertificate({
      namespaceId: this.policyNamespaceID,
      namespaceProvider: this.policyNamespaceProvider,
      id: this.certPolicy.id,
      enrollCertificateRequest: {
        publicKey: (await crypto.subtle.exportKey(
          "jwk",
          this.keyPair.publicKey
        )) as JsonWebKey,
        withOneTimePkcs12Key: true,
      },
    });
    return new EnrollmentSession(
      this.policyNamespaceProvider,
      this.policyNamespaceID,
      this.certPolicy,
      this.keyPair,
      resp
    );
  }
}

function putUint32(array: Uint8Array, offset: number, value: number) {
  array[offset] = value >>> 24;
  array[offset + 1] = value >>> 16;
  array[offset + 2] = value >>> 8;
  array[offset + 3] = value;
  return offset + 4;
}

function putBuffer(array: Uint8Array, offset: number, value: ArrayBuffer) {
  array.set(new Uint8Array(value), offset);
  return offset + value.byteLength;
}

async function ecdhEsKdf(z: ArrayBuffer, alg: string): Promise<ArrayBuffer> {
  const finalBufferLen = 4 + z.byteLength + 4 + alg.length + 4 * 3;
  const finalBuffer = new Uint8Array(finalBufferLen);
  let offset = 0;
  offset = putUint32(finalBuffer, offset, 1);
  offset = putBuffer(finalBuffer, offset, z);
  offset = putUint32(finalBuffer, offset, alg.length);
  offset = putBuffer(finalBuffer, offset, new TextEncoder().encode(alg));
  offset = putUint32(finalBuffer, offset, 0);
  offset = putUint32(finalBuffer, offset, 0);
  offset = putUint32(finalBuffer, offset, z.byteLength * 8);
  console.log(finalBuffer, offset);
  return crypto.subtle.digest("SHA-256", finalBuffer);
}

async function ecdhEsDeriveA256GCM(remoteKey: CryptoKey, localKey: CryptoKey) {
  const derivedBits = await window.crypto.subtle.deriveBits(
    {
      name: "ECDH",
      public: remoteKey,
    },
    localKey,
    256
  );
  const ecdhEsCEK = await ecdhEsKdf(derivedBits, "A256GCM");
  return await crypto.subtle.importKey(
    "raw",
    ecdhEsCEK,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
}

class EnrollmentPKCS12ExchangeSession {
  constructor(
    public readonly session: EnrollmentSession,
    public readonly legacy: boolean,
    public readonly pkcs12Blob?: Blob,
    public readonly result?: ExchangePKCS12Result
  ) {}

  public async fetchPKCS12(api: AdminApi, passwordProtected: boolean) {
    if (!this.session.enrollCertResponse) {
      throw new Error("no base session");
    }
    const [payload, encKey] = await this.preparePayload();
    const [ns, certId] =
      this.session.enrollCertResponse.identififier.split("/");
    const [namespaceProvider, namespaceId] = ns.split(":");
    const resp = await api.exchangePKCS12({
      namespaceId,
      namespaceProvider: namespaceProvider as NamespaceProvider,
      id: certId,
      exchangePKCS12Request: {
        payload,
        passwordProtected,
        legacy: this.legacy,
      },
    });
    const resPayload = resp.payload.split(".");
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
        additionalData: new TextEncoder().encode(resPayload[0]),
      },
      encKey,
      encrypted
    );
    return new EnrollmentPKCS12ExchangeSession(
      this.session,
      this.legacy,
      new Blob([decrypted], {
        type: "application/x-pkcs12",
      }),
      resp
    );
  }

  private async preparePayload(): Promise<[string, CryptoKey]> {
    if (!this.session.enrollCertResponse?.oneTimePkcs12Key) {
      throw new Error("Missing oneTimePkcs12Key");
    }
    if (!this.session.keyPair) {
      throw new Error("Missing keyPair");
    }
    console.log("enter");
    const remoteKey = await crypto.subtle.importKey(
      "jwk",
      this.session.enrollCertResponse?.oneTimePkcs12Key as JsonWebKey,
      { name: "ECDH", namedCurve: "P-384" },
      false,
      []
    );
    console.log("remote key ok");
    const ephemeralKey = await window.crypto.subtle.generateKey(
      {
        name: "ECDH",
        namedCurve: "P-384",
      },
      false,
      ["deriveKey", "deriveBits"]
    );
    const encKey = await ecdhEsDeriveA256GCM(
      remoteKey,
      ephemeralKey.privateKey
    );

    const joseHeader = {
      alg: "ECDH-ES",
      enc: "A256GCM",
      epk: await crypto.subtle.exportKey("jwk", ephemeralKey.publicKey),
    };

    const joseEncoded = base64StdEncodedToUrlEncoded(
      btoa(JSON.stringify(joseHeader))
    );

    const iv = crypto.getRandomValues(new Uint8Array(12));
    const plainTextBuffer = new TextEncoder().encode(
      JSON.stringify(
        await crypto.subtle.exportKey("jwk", this.session.keyPair.privateKey)
      )
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
        "",
        base64UrlEncodeBuffer(iv),
        base64UrlEncodeBuffer(ciphertext),
        base64UrlEncodeBuffer(tag),
      ].join("."),
      encKey,
    ];
  }
}

function ExchangePKCS12Control({
  baseSession,
  legacy,
  passwordProtected,
  label,
}: {
  baseSession: EnrollmentSession | undefined;
  legacy: boolean;
  passwordProtected: boolean;
  label: React.ReactNode;
}) {
  const api = useAdminApi();
  const { run, data } = useRequest(
    async () => {
      if (!api || !baseSession) return;
      const session = await new EnrollmentPKCS12ExchangeSession(
        baseSession,
        legacy
      ).fetchPKCS12(api, passwordProtected);
      return session;
    },
    { manual: true }
  );

  const [blobUrl, setBlobUrl] = useState<string>();

  useEffect(() => {
    if (data?.pkcs12Blob) {
      const url = URL.createObjectURL(data.pkcs12Blob);
      setBlobUrl(url);
      return () => {
        URL.revokeObjectURL(url);
        setBlobUrl(undefined);
      };
    }
  }, [data, setBlobUrl]);

  return (
    <div>
      <Button type="primary" onClick={run}>
        {label}
      </Button>
      {!legacy && (
        <div>If having issues importing certificate, try legacy options</div>
      )}
      {blobUrl && (
        <>
          <div>
            Download PKCS#12 file:{" "}
            <a
              href={blobUrl}
              download={`${data?.session.enrollCertResponse?.id}${
                legacy ? "-legacy" : ""
              }.p12`}
            >
              {data?.session.enrollCertResponse?.id}
              {legacy ? "-legacy" : ""}.p12
            </a>
          </div>
          <div>
            Password:{" "}
            <span className="font-mono">{data?.result?.password}</span>
          </div>
        </>
      )}
    </div>
  );
}

export function CertWebEnroll({
  certPolicy,
}: {
  certPolicy: CertificatePolicy;
}) {
  const [session, setSession] = useState<EnrollmentSession | undefined>(
    undefined
  );
  const { namespaceProvider, namespaceId } = useNamespace();
  const api = useAuthedClientV2(AdminApi);
  const unmounted = useUnmountedRef();
  const [pemBlob, setPemBlob] = useState<Blob | undefined>(undefined);
  const [pemBlobUrl, setPemBlobUrl] = useState<string | undefined>(undefined);

  useEffect(() => {
    if (session) {
      if (!session.keyPair) {
        session.withGeneratedKeyPair().then((s) => {
          if (s !== session) {
            if (!unmounted.current) {
              setSession((p) => (p == session ? s : p));
            }
          }
        });
      } else if (!session.enrollCertResponse) {
        session
          .withEnrollmentResponse(api)
          .then((s) => {
            if (s !== session) {
              if (!unmounted.current) {
                setSession((p) => (p == session ? s : p));
                return s.getPemBlob();
              }
            }
          })
          .then((pemBlob) => {
            if (pemBlob) {
              if (!unmounted.current) {
                setPemBlob(pemBlob);
              }
            }
          });
      }
    }
  }, [session, api, unmounted]);

  useEffect(() => {
    if (pemBlob) {
      const url = URL.createObjectURL(pemBlob);
      setPemBlobUrl(URL.createObjectURL(pemBlob));
      return () => {
        URL.revokeObjectURL(url);
        setPemBlobUrl(undefined);
      };
    }
  }, [pemBlob]);

  const [pkcs12PasswordProtected, setPkcs12PasswordProtected] =
    useState<boolean>(true);

  return (
    <div className="space-y-4">
      {!session && (
        <Button
          type="primary"
          onClick={() => {
            setSession(
              new EnrollmentSession(namespaceProvider, namespaceId, certPolicy)
            );
          }}
        >
          Begin enrollment
        </Button>
      )}
      {pemBlobUrl && (
        <div className="ring p-4 rounded-md ring-green-600 space-y-4 bg-green-50">
          <span className="text-red-600">
            This is the only time you can download the private key.
          </span>
          <br />
          Download certificate:{" "}
          <a
            href={pemBlobUrl}
            download={`${session?.enrollCertResponse?.id}.pem`}
          >
            {session?.enrollCertResponse?.id}.pem
          </a>
        </div>
      )}
      {session?.enrollCertResponse && (
        <div className="ring-1 p-4 rounded-md ring-yellow-500 bg-yellow-50 space-y-4">
          <div>
            Select following options to download file can be imported to Windows
            and MacOS. Your private key will be encrypted and sent to the server
            to create the bundle PKCS12 file.
          </div>
          <div className="text-red-600">
            Note: Your private key will leave your browser.
          </div>
          <Checkbox
            checked={pkcs12PasswordProtected}
            onChange={(e) => {
              setPkcs12PasswordProtected(e.target.checked);
            }}
          >
            Password protected
          </Checkbox>
          <ExchangePKCS12Control
            baseSession={session}
            legacy={false}
            passwordProtected={pkcs12PasswordProtected}
            label="Get PKCS 12 file (Modern)"
          />
          <ExchangePKCS12Control
            baseSession={session}
            legacy={true}
            passwordProtected={pkcs12PasswordProtected}
            label="Get PKCS 12 file (Legacy)"
          />
          {/* <div>
            <Button type="primary">Get PKCS 12 file (Modern encryption)</Button>
          </div>
          <div className="flex items-center flex-wrap gap-y-1 gap-x-2">
            <Button type="primary">Get PKCS 12 file (Legacy encryption)</Button>{" "}
            Works with macOS 14.1.1 and iOS 17.1.1 and later
          </div> */}
          <div>
            Alternatively you may use OpenSSL to convert the downloaded PEM file
            to PKCS12 file
            <code className="block bg-neutral-50 ring-1 rounded-md p-4">
              openssl pkcs12 -export -in (certificate.pem) -inkey
              (certificate.pem) -out certificate.p12 [-legacy]
            </code>
          </div>
        </div>
      )}
    </div>
  );
}
