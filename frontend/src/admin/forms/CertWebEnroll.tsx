import { useUnmountedRef } from "ahooks";
import { Button } from "antd";
import { useEffect, useState } from "react";
import { JsonWebSignatureAlgorithm } from "../../generated";
import {
  AdminApi,
  Certificate,
  CertificatePolicy,
  JsonWebKey,
  NamespaceProvider,
} from "../../generated/apiv2";
import {
  base64UrlDecodeBuffer,
  base64UrlEncodeBuffer,
} from "../../utils/encodingUtils";
import { useAuthedClientV2 } from "../../utils/useCertsApi";
import { useNamespace } from "../contexts/NamespaceContextRouteProvider";

class EnrollmentSession {
  constructor(
    public readonly policyNamespaceProvider: NamespaceProvider,
    public readonly policyNamespaceID: string,
    public readonly certPolicy: CertificatePolicy,
    public readonly keyPair?: CryptoKeyPair,
    public readonly enrollCertResponse?: Certificate,
    public readonly signedCertResponse?: Certificate
  ) {}

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
      } else if (!session.enrollCertResponse) {
        session.withEnrollmentResponse(api).then((s) => {
          if (s !== session) {
            if (!unmounted.current) {
              setSession((p) => (p == session ? s : p));
            }
          }
        });
      }
    }
  }, [session]);

  console.log(session);

  return (
    <div>
      <Button
        onClick={() => {
          setSession(
            new EnrollmentSession(namespaceProvider, namespaceId, certPolicy)
          );
        }}
      >
        Begin enrollment
      </Button>
    </div>
  );
}
