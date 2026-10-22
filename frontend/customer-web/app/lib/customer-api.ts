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

export interface Domain { id: string; accountId: string; name: string; kind: string; status: string; documentRoot?: string; createdAt: string; }
export interface DNSZone { id: string; accountId: string; domainId: string; name: string; status: string; createdAt: string; }
export interface DNSRecord { id: string; zoneId: string; name: string; type: string; ttl: number; value: string; }

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

function parseDomains(value: unknown): Domain[] {
  if (!Array.isArray(value)) throw new Error("The domains API returned an invalid collection.");
  return value.map((item) => {
    if (!isRecord(item)) throw new Error("The domains API returned an invalid domain.");
    return { id: readRequiredString(item.id, "domain id"), accountId: readRequiredString(item.account_id, "domain account id"), name: readRequiredString(item.name, "domain name"), kind: readRequiredString(item.kind, "domain kind"), status: readRequiredString(item.status, "domain status"), documentRoot: typeof item.document_root === "string" ? item.document_root : undefined, createdAt: readRequiredString(item.created_at, "domain created time") };
  });
}

function parseZones(value: unknown): DNSZone[] {
  if (!Array.isArray(value)) throw new Error("The DNS zones API returned an invalid collection.");
  return value.map((item) => {
    if (!isRecord(item)) throw new Error("The DNS zones API returned an invalid zone.");
    return { id: readRequiredString(item.id, "zone id"), accountId: readRequiredString(item.account_id, "zone account id"), domainId: readRequiredString(item.domain_id, "zone domain id"), name: readRequiredString(item.name, "zone name"), status: readRequiredString(item.status, "zone status"), createdAt: readRequiredString(item.created_at, "zone created time") };
  });
}

function parseRecords(value: unknown): DNSRecord[] {
  if (!Array.isArray(value)) throw new Error("The DNS records API returned an invalid collection.");
  return value.map((item) => {
    if (!isRecord(item) || typeof item.ttl !== "number") throw new Error("The DNS records API returned an invalid record.");
    return { id: readRequiredString(item.id, "record id"), zoneId: readRequiredString(item.zone_id, "record zone id"), name: readRequiredString(item.name, "record name"), type: readRequiredString(item.type, "record type"), ttl: item.ttl, value: readRequiredString(item.value, "record value") };
  });
}

export function listCustomerDomains(accountId: string, signal?: AbortSignal) { return getCollection(`/api/v1/domains?account_id=${encodeURIComponent(accountId)}`, parseDomains, signal); }
export function listCustomerZones(accountId: string, signal?: AbortSignal) { return getCollection(`/api/v1/dns/zones?account_id=${encodeURIComponent(accountId)}`, parseZones, signal); }
export function listCustomerRecords(zoneId: string, signal?: AbortSignal) { return getCollection(`/api/v1/dns/zones/${encodeURIComponent(zoneId)}/records`, parseRecords, signal); }

async function postCommand<T>(path: string, body: unknown, parse: (value: unknown) => T): Promise<T> {
  const response = await fetch(`${apiBaseUrl}${path}`, { method: "POST", headers: { Accept: "application/json", "Content-Type": "application/json" }, body: JSON.stringify(body) });
  const payload: unknown = await response.json().catch(() => null);
  if (!response.ok || !isRecord(payload) || !("data" in payload)) throw new Error(isRecord(payload) && typeof payload.error === "string" ? payload.error : `Request failed with status ${response.status}.`);
  return parse(payload.data);
}

export function createCustomerDomain(command: { account_id: string; name: string }) { return postCommand("/api/v1/domains", command, (value) => parseDomains([value])[0]); }
export function createCustomerZone(command: { account_id: string; domain_id: string }) { return postCommand("/api/v1/dns/zones", command, (value) => parseZones([value])[0]); }
export function createCustomerRecord(zoneId: string, command: { name: string; type: string; ttl: number; value: string }) { return postCommand(`/api/v1/dns/zones/${encodeURIComponent(zoneId)}/records`, command, (value) => parseRecords([value])[0]); }

export function getCustomerDashboardData(signal?: AbortSignal): Promise<CustomerDashboardData> {
  return Promise.all([
    getCollection("/api/v1/accounts", parseAccounts, signal),
    getCollection("/api/v1/servers", parseServers, signal),
  ]).then(([accounts, servers]) => ({ accounts, servers }));
}
