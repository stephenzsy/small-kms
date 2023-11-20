import { useUnmountedRef } from "ahooks";
import { Button } from "antd";
import { useEffect, useState } from "react";
import { JsonWebSignatureAlgorithm } from "../../generated";
import {
  AdminApi,
  CertificatePolicy,
  JsonWebSignatureKey,
} from "../../generated/apiv2";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useNamespace } from "../contexts/NamespaceContextRouteProvider";

class EnrollmentSession {
  constructor(
    public readonly certPolicy: CertificatePolicy,
    public readonly keyPair?: CryptoKeyPair
  ) {}

  public async withGeneratedKeyPair() {
    if (this.keyPair) {
      return this;
    }
    let generateKeyParams: RsaHashedKeyGenParams | EcKeyGenParams;
    switch (this.certPolicy.keySpec.alg) {
      case JsonWebSignatureAlgorithm.Rs256:
        generateKeyParams = {
          name: "RSASSA-PKCS1-v1_5",
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: "SHA-256",
        };

        break;
      case JsonWebSignatureAlgorithm.Rs384:
        generateKeyParams = {
          name: "RSASSA-PKCS1-v1_5",
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: "SHA-384",
        };
        break;
      case JsonWebSignatureAlgorithm.Rs512:
        generateKeyParams = {
          name: "RSASSA-PKCS1-v1_5",
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: "SHA-512",
        };
        break;
      case JsonWebSignatureAlgorithm.Es256:
        generateKeyParams = {
          name: "ECDSA",
          namedCurve: "P-256",
        };
        break;
      case JsonWebSignatureAlgorithm.Es384:
        generateKeyParams = {
          name: "ECDSA",
          namedCurve: "P-384",
        };
        break;
      case JsonWebSignatureAlgorithm.Es512:
        generateKeyParams = {
          name: "ECDSA",
          namedCurve: "P-521",
        };
        break;
      case JsonWebSignatureAlgorithm.Ps256:
        generateKeyParams = {
          name: "RSA-PSS",
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: "SHA-256",
        };
        break;
      case JsonWebSignatureAlgorithm.Ps384:
        generateKeyParams = {
          name: "RSA-PSS",
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: "SHA-384",
        };
        break;
      case JsonWebSignatureAlgorithm.Ps512:
        generateKeyParams = {
          name: "RSA-PSS",
          modulusLength: this.certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: "SHA-512",
        };
        break;
      default:
        throw new Error("Unsupported algorithm");
    }
    const keyPair = await crypto.subtle.generateKey(generateKeyParams, true, [
      "sign",
      "verify",
    ]);
    return new EnrollmentSession(this.certPolicy, keyPair);
  }
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
      } else {
        (async () =>
          api.enrollCertificate({
            namespaceId,
            namespaceProvider,
            id: certPolicy.id,
            enrollCertificateRequest: {
              publicKey: (await crypto.subtle.exportKey(
                "jwk",
                session.keyPair!.publicKey
              )) as JsonWebSignatureKey,
            },
          }))();
      }
    }
  }, [session]);

  console.log(session);

  return (
    <div>
      <Button
        onClick={() => {
          setSession(new EnrollmentSession(certPolicy));
        }}
      >
        Begin enrollment
      </Button>
    </div>
  );
}
