import { useRequest } from "ahooks";
import { Button } from "antd";
import { useContext, useMemo } from "react";
import { useAppAuthContext } from "../auth/AuthProvider";
import {
  AdminApi,
  CertPolicy,
  JsonWebKey,
  JsonWebKeySignatureAlgorithm,
  JsonWebSignatureKey,
} from "../generated";
import { useAuthedClient } from "../utils/useCertsApi";
import { NamespaceContext } from "./contexts/NamespaceContext";

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

  const { data, run } = useRequest(
    async () => {
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

        await api.enrollCertificate({
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

        console.log(
          JSON.stringify(
            await window.crypto.subtle.exportKey("jwk", keypair.publicKey)
          )
        );

        console.log(
          JSON.stringify(
            await window.crypto.subtle.exportKey("jwk", keypair.privateKey)
          )
        );
        console.log(proofJwt);
        return undefined;
      }
    },
    { manual: true }
  );

  return (
    <div>
      <p>Private key will not be sent</p>
      <Button onClick={run} type="primary">
        Create key pair
      </Button>
    </div>
  );
}
