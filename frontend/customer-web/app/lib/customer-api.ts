export type HostingAccountStatus = "active" | "pending" | "suspended" | string;

export interface HostingAccount {
  id: string;
  username: string;
  status: HostingAccountStatus;
}

export type ServerRegistrationStatus =
  | "pending"
  | "enrolling"
  | "active"
  | "degraded"
  | "offline"
  | "revoked"
  | string;

export type ServerConnectionStatus = "unknown" | "connected" | "disconnected" | string;

export interface ServerCapability {
  name: string;
  version?: string;
  status: "available" | "unavailable" | "unknown" | string;
}

export interface ServerStatus {
  serverId: string;
  registrationStatus: ServerRegistrationStatus;
  connectionStatus: ServerConnectionStatus;
  agentVersion?: string;
  lastSeenAt?: string;
  observedAt: string;
  capabilities: ServerCapability[];
}

export interface CustomerDashboardData {
  accounts: HostingAccount[];
  servers: ServerStatus[];
}

const apiBaseUrl = (process.env.NEXT_PUBLIC_HARBOR_API_BASE_URL ?? "http://localhost:8080").replace(/\/$/, "");

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function readRequiredString(value: unknown, field: string): string {
  if (typeof value !== "string" || value.length === 0) {
    throw new Error(`The API returned an invalid ${field}.`);
  }
  return value;
}

function parseAccounts(value: unknown): HostingAccount[] {
  if (!Array.isArray(value)) {
    throw new Error("The accounts API returned an invalid collection.");
  }

  return value.map((item) => {
    if (!isRecord(item)) {
      throw new Error("The accounts API returned an invalid account.");
    }

    return {
      id: readRequiredString(item.id, "account id"),
      username: readRequiredString(item.username, "account username"),
      status: readRequiredString(item.status, "account status"),
    };
  });
}

function parseServers(value: unknown): ServerStatus[] {
  if (!Array.isArray(value)) {
    throw new Error("The servers API returned an invalid collection.");
  }

  return value.map((item) => {
    if (!isRecord(item)) {
      throw new Error("The servers API returned an invalid server.");
    }

    const capabilities = item.capabilities;
    if (!Array.isArray(capabilities)) {
      throw new Error("The servers API returned invalid capabilities.");
    }

    return {
      serverId: readRequiredString(item.server_id, "server id"),
      registrationStatus: readRequiredString(item.registration_status, "server registration status"),
      connectionStatus: readRequiredString(item.connection_status, "server connection status"),
      agentVersion: typeof item.agent_version === "string" ? item.agent_version : undefined,
      lastSeenAt: typeof item.last_seen_at === "string" ? item.last_seen_at : undefined,
      observedAt: readRequiredString(item.observed_at, "server observation time"),
      capabilities: capabilities.map((capability) => {
        if (!isRecord(capability)) {
          throw new Error("The servers API returned an invalid capability.");
        }

        return {
          name: readRequiredString(capability.name, "capability name"),
          version: typeof capability.version === "string" ? capability.version : undefined,
          status: readRequiredString(capability.status, "capability status"),
        };
      }),
    };
  });
}

async function getCollection<T>(path: string, parse: (value: unknown) => T, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    headers: { Accept: "application/json" },
    signal,
  });

  const body: unknown = await response.json().catch(() => null);
  if (!response.ok) {
    const message = isRecord(body) && typeof body.error === "string" ? body.error : `Request failed with status ${response.status}.`;
    throw new Error(message);
  }

  if (!isRecord(body) || !("data" in body)) {
    throw new Error("The API returned an invalid response envelope.");
  }

  return parse(body.data);
}

export function getCustomerDashboardData(signal?: AbortSignal): Promise<CustomerDashboardData> {
  return Promise.all([
    getCollection("/api/v1/accounts", parseAccounts, signal),
    getCollection("/api/v1/servers", parseServers, signal),
  ]).then(([accounts, servers]) => ({ accounts, servers }));
}
