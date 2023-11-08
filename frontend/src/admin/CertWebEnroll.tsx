import { useRequest } from "ahooks";
import { Alert, Button, Steps, Typography } from "antd";
import {
  useContext,
  useEffect,
  useLayoutEffect,
  useMemo,
  useState,
} from "react";
import { useAppAuthContext } from "../auth/AuthProvider";
import {
  AdminApi,
  CertPolicy,
  Certificate,
  JsonWebKey,
  JsonWebKeySignatureAlgorithm,
  JsonWebSignatureKey,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";
import AnchorLink from "antd/es/anchor/AnchorLink";
import { Link } from "../components/Link";
import { toPEMBlock } from "../utils/encodingUtils";
import { base64UrlEncodedToStdEncoded } from "../utils/encodingUtils";

function useKeyGenAlgParams(certPolicy: CertPolicy | undefined) {
  return useMemo((): RsaHashedKeyGenParams | EcKeyGenParams | undefined => {
    if (!certPolicy) {
      return undefined;
    }
    switch (certPolicy.keySpec.alg) {
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmRS256:
        return {
          name: "RSASSA-PKCS1-v1_5",
          hash: "SHA-256",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmRS384:
        return {
          name: "RSASSA-PKCS1-v1_5",
          hash: "SHA-384",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmRS512:
        return {
          name: "RSASSA-PKCS1-v1_5",
          hash: "SHA-512",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmES256:
        return {
          name: "ECDSA",
          namedCurve: "P-256",
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmES384:
        return {
          name: "ECDSA",
          namedCurve: "P-384",
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmES512:
        return {
          name: "ECDSA",
          namedCurve: "P-521",
        };

      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmPS256:
        return {
          name: "RSA-PSS",
          hash: "SHA-256",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmPS384:
        return {
          name: "RSA-PSS",
          hash: "SHA-384",
          modulusLength: certPolicy.keySpec.keySize!,
          publicExponent: new Uint8Array([1, 0, 1]),
        };
      case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmPS512:
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

type ProofJwtClaims = {
  aud: string;
  iss: string;
  nbf: number;
  exp: number;
};

function getProofJwtClaims(aud: string, issuer: string): ProofJwtClaims {
  const now = Math.floor(Date.now() / 1000);
  return {
    aud,
    iss: issuer,
    nbf: now,
    exp: now + 600,
  };
}

function encodeBase64Url(data: string) {
  const stdEncoded = btoa(data);
  return stdEncoded.replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
}

function getSignatureCryptoAlg(
  alg: JsonWebKeySignatureAlgorithm
): AlgorithmIdentifier | RsaPssParams | EcdsaParams | undefined {
  switch (alg) {
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmRS256:
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmRS384:
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmRS512:
      return "RSASSA-PKCS1-v1_5";
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmES256:
      return {
        hash: "SHA-256",
        name: "ECDSA",
      };
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmES384:
      return {
        hash: "SHA-384",
        name: "ECDSA",
      };
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmES512:
      return {
        hash: "SHA-512",
        name: "ECDSA",
      };
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmPS256:
      return {
        hash: "SHA-256",
        name: "RSA-PSS",
      };
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmPS384:
      return {
        hash: "SHA-384",
        name: "RSA-PSS",
      };
    case JsonWebKeySignatureAlgorithm.JsonWebKeySignatureAlgorithmPS512:
      return {
        hash: "SHA-512",
        name: "RSA-PSS",
      };
  }
  return undefined;
}

async function signProofJwt(
  proofJwtClaims: ProofJwtClaims,
  jwtAlg: JsonWebKeySignatureAlgorithm,
  signAlg: AlgorithmIdentifier | RsaPssParams | EcdsaParams,
  key: CryptoKey
) {
  const encoder = new TextEncoder();
  const toBeSigned = [
    encodeBase64Url(
      JSON.stringify({
        alg: jwtAlg,
        typ: "JWT",
      })
    ),
    encodeBase64Url(JSON.stringify(proofJwtClaims)),
  ].join(".");
  const data = encoder.encode(toBeSigned);
  const signatureBuf = await window.crypto.subtle.sign(signAlg, key, data);
  const encodedSignature = encodeBase64Url(
    String.fromCharCode(...new Uint8Array(signatureBuf))
  );
  return toBeSigned + "." + encodedSignature;
}

export function CertWebEnroll({
  certPolicy,
}: {
  certPolicy: CertPolicy | undefined;
}) {
  const alg = useKeyGenAlgParams(certPolicy);
  const { account } = useAppAuthContext();
  const { namespaceId, namespaceKind } = useContext(NamespaceContext);

  const api = useAuthedClient(AdminApi);

  const { data, run, loading } = useRequest(
    async (): Promise<[Blob, Certificate] | undefined> => {
      if (alg && account && certPolicy && certPolicy.keySpec.alg) {
        const keypair = await window.crypto.subtle.generateKey(alg, true, [
          "sign",
          "verify",
        ]);

        const proofJwtClaims = getProofJwtClaims(
          namespaceId,
          account.localAccountId
        );

        const proofJwt = await signProofJwt(
          proofJwtClaims,
          certPolicy.keySpec.alg,
          getSignatureCryptoAlg(certPolicy.keySpec.alg)!,
          keypair.privateKey
        );

        const enrollResp = await api.enrollCertificate({
          enrollCertificateRequest: {
            enrollmentType: "group-memeber",
            proof: proofJwt,
            publicKey: (await window.crypto.subtle.exportKey(
              "jwk",
              keypair.publicKey
            )) as JsonWebSignatureKey,
          },
          namespaceId,
          namespaceKind,
          resourceId: certPolicy.id,
        });

        const privateKeyDer = await window.crypto.subtle.exportKey(
          "pkcs8",
          keypair.privateKey
        );

        return [
          new Blob(
            [
              [
                toPEMBlock(
                  btoa(String.fromCharCode(...new Uint8Array(privateKeyDer))),
                  "PRIVATE KEY"
                ),
                ...(enrollResp.jwk.x5c ?? []).map((cert) =>
                  toPEMBlock(base64UrlEncodedToStdEncoded(cert), "CERTIFICATE")
                ),
              ].join("\n"),
            ],
            {
              type: "application/x-pem-file",
            }
          ),
          enrollResp,
        ];
      }
    },
    { manual: true }
  );

  const [blobUrl, setBlobUrl] = useState<string>();

  const [certBlob, enrolledCert] = data ?? [];

  useEffect(() => {
    if (certBlob) {
      const url = URL.createObjectURL(certBlob);
      setBlobUrl(url);
      return () => {
        URL.revokeObjectURL(url);
      };
    }
  }, [certBlob]);

  return (
    <div className="space-y-4">
      <Alert type="info" message="Private key does not leave your browser" />
      {!enrolledCert && (
        <Button onClick={run} type="primary" loading={loading}>
          Begin enroll
        </Button>
      )}
      {enrolledCert && (
        <div className="ring-1 ring-neutral-400 rounded-md p-4 space-y-4">
          <div className="text-lg font-semibold">Download certificate</div>
          <Alert
            type="warning"
            message="This is the only time you zcan download the certificate. Please save it in a safe place."
          />
          {blobUrl && (
            <Button
              type="primary"
              href={blobUrl}
              download={`${enrolledCert.id}.pem`}
            >
              Download certificate
            </Button>
          )}
        </div>
      )}
    </div>
  );
}
