import { DefaultAzureCredential } from "@azure/identity";
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

class ServiceCredential {
  private static readonly azCredential = new DefaultAzureCredential({
    managedIdentityClientId: process.env.MANAGED_IDENTITY_CLIENT_ID,
  });

  public static readonly accessToken = async (): Promise<string> => {
    return (
      await ServiceCredential.azCredential.getToken(process.env.API_SCOPE!)
    ).token;
  };
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

  public readonly accessToken = ServiceCredential.accessToken;

  public async getAuthHeaders(
    allowNonAdmin: boolean = false
  ): Promise<Record<string, string>> {
    if (allowNonAdmin || this.isAdmin) {
      return {
        Authorization: `Bearer ${await ServiceCredential.accessToken()}`,
        "X-Caller-Principal-Name": this.principalName!,
        "X-Caller-Principal-Id": this.principalId!,
      };
    }
    return {};
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
