import { headers } from "next/headers";

type XMsClientPrincipalEntry = {
  readonly typ: "roles";
  readonly val: string;
};

interface XMsClientPrincipal {
  readonly claims: readonly XMsClientPrincipalEntry[];
}

interface XMsClient {
  readonly principalName?: string;
  readonly principalId?: string;
  readonly principal?: XMsClientPrincipal;
}

export interface MsAuthClient {
  readonly principalName: string | undefined;
  readonly isAdmin: boolean;
}

export class MsAuthServer {
  public constructor(private readonly config: XMsClient) {}

  public get isAdmin(): boolean {
    return (
      this.config.principal?.claims.some(
        (x) => x.typ === "roles" && x.val === "App.Admin"
      ) ?? false
    );
  }

  public get principalName(): string | undefined {
    return this.config.principalName;
  }

  public readonly client: MsAuthClient = {
    principalName: this.principalName,
    isAdmin: this.isAdmin,
  };
}

export function getMsAuth(): MsAuthServer {
  if (process.env.USE_STUB_AUTH === "admin") {
    return new MsAuthServer({
      principalName: "admin@example.com",
      principalId: "00000000-0000-0000-0000-000000000002",
      principal: {
        claims: [{ typ: "roles", val: "App.Admin" }],
      },
    });
  }

  let principal = undefined;
  try {
    const principalEncoded =
      headers().get("x-ms-client-principal") ?? undefined;

    principal =
      principalEncoded &&
      JSON.parse(Buffer.from(principalEncoded, "base64").toString("utf-8"));
  } catch {}
  return new MsAuthServer({
    principalName: headers().get("x-ms-client-principal-name") ?? undefined,
    principalId: headers().get("x-ms-client-principal-id") ?? undefined,
    principal,
  });
}
