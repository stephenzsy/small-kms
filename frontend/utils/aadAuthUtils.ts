import { headers } from "next/headers";

type XMsClientPrincipalEntry = {
  readonly typ: "roles";
  readonly val: string;
};

interface XMsClientPrincipalParsed {
  readonly claims: readonly XMsClientPrincipalEntry[];
}

interface XMsClient {
  readonly principalName?: string;
  readonly principalId?: string;
  readonly principal?: string;
  readonly principalParsed?: XMsClientPrincipalParsed;
}

export interface MsAuthClient {
  readonly principalName: string | undefined;
  readonly isAdmin: boolean;
}

export class MsAuthServer {
  public constructor(private readonly config: XMsClient) {}

  public get isAdmin(): boolean {
    return (
      this.config.principalParsed?.claims.some(
        (x) => x.typ === "roles" && x.val === "App.Admin"
      ) ?? false
    );
  }

  public get principalName(): string | undefined {
    return this.config.principalName;
  }

  public get principalId(): string | undefined {
    return this.config.principalId;
  }

  public readonly client: MsAuthClient = {
    principalName: this.principalName,
    isAdmin: this.isAdmin,
  };

  public get principal(): string | undefined {
    return this.config.principal;
  }
}

export function getMsAuth(): MsAuthServer {
  if (process.env.USE_STUB_AUTH === "admin") {
    const principalParsed: XMsClientPrincipalParsed = {
      claims: [{ typ: "roles", val: "App.Admin" }],
    };
    return new MsAuthServer({
      principalName: "admin@example.com",
      principalId: "00000000-0000-0000-0000-000000000002",
      principal: Buffer.from(JSON.stringify(principalParsed), "utf-8").toString(
        "base64"
      ),
      principalParsed,
    });
  }

  const principal = headers().get("x-ms-client-principal") ?? undefined;

  let principalParsed: XMsClientPrincipalParsed | undefined;
  try {
    principalParsed =
      principal &&
      JSON.parse(Buffer.from(principal, "base64").toString("utf-8"));
  } catch {}
  return new MsAuthServer({
    principalName: headers().get("x-ms-client-principal-name") ?? undefined,
    principalId: headers().get("x-ms-client-principal-id") ?? undefined,
    principal,
    principalParsed,
  });
}