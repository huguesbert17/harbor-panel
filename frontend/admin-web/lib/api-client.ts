export type HostingAccount = {
  id: string;
  username: string;
  status: string;
};

export type ServerCapability = {
  name: string;
  version?: string;
  status: string;
};

export type ManagedServer = {
  serverId: string;
  registrationStatus: string;
  connectionStatus: string;
  agentVersion?: string;
  lastSeenAt?: string;
  observedAt: string;
  capabilities: ServerCapability[];
};
export type Domain = { id: string; accountId: string; name: string; kind: string; status: string; documentRoot?: string; createdAt: string };
export type DNSZone = { id: string; accountId: string; domainId: string; name: string; status: string; createdAt: string };
export type DNSRecord = { id: string; zoneId: string; name: string; type: string; ttl: number; value: string };

type CollectionResponse<T> = {
  data?: T[];
};

type ApiErrorResponse = {
  error?: string;
};

async function get<T>(path: string): Promise<T> {
  let response: Response;

  try {
    response = await fetch(`/api/control-plane${path}`, { headers: { Accept: "application/json" }, cache: "no-store" });
  } catch {
    throw new Error("The control-plane API could not be reached. Check that the API is running and the admin API URL is configured.");
  }

  const payload = await response.json().catch(() => null) as ApiErrorResponse | null;
  if (!response.ok) {
    throw new Error(payload?.error || `The control-plane API returned HTTP ${response.status}.`);
  }

  return payload as T;
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const response = await fetch(`/api/control-plane${path}`, { method: "POST", headers: { Accept: "application/json", "Content-Type": "application/json" }, body: JSON.stringify(body) });
  const payload = await response.json().catch(() => null) as (T & ApiErrorResponse) | null;
  if (!response.ok) throw new Error(payload?.error || `The control-plane API returned HTTP ${response.status}.`);
  return payload as T;
}

function parseAccounts(payload: CollectionResponse<{ id: string; username: string; status: string }>): HostingAccount[] {
  if (!Array.isArray(payload.data)) {
    throw new Error("The accounts response was not in the expected format.");
  }

  return payload.data.map((account) => ({ id: account.id, username: account.username, status: account.status }));
}

function parseServers(payload: CollectionResponse<{
  server_id: string;
  registration_status: string;
  connection_status: string;
  agent_version?: string;
  last_seen_at?: string;
  observed_at: string;
  capabilities?: ServerCapability[];
}>): ManagedServer[] {
  if (!Array.isArray(payload.data)) {
    throw new Error("The servers response was not in the expected format.");
  }

  return payload.data.map((server) => ({
    serverId: server.server_id,
    registrationStatus: server.registration_status,
    connectionStatus: server.connection_status,
    agentVersion: server.agent_version,
    lastSeenAt: server.last_seen_at,
    observedAt: server.observed_at,
    capabilities: Array.isArray(server.capabilities) ? server.capabilities : [],
  }));
}

export async function listHostingAccounts() {
  return parseAccounts(await get<CollectionResponse<{ id: string; username: string; status: string }>>("/accounts"));
}

export async function listManagedServers() {
  return parseServers(await get<CollectionResponse<{
    server_id: string;
    registration_status: string;
    connection_status: string;
    agent_version?: string;
    last_seen_at?: string;
    observed_at: string;
    capabilities?: ServerCapability[];
  }>>("/servers"));
}

export async function loadAdminDashboard() {
  const [accounts, servers] = await Promise.all([listHostingAccounts(), listManagedServers()]);
  return { accounts, servers };
}

export async function listDomains(accountId?: string) { const payload = await get<CollectionResponse<{ id: string; account_id: string; name: string; kind: string; status: string; document_root?: string; created_at: string }>>(`/domains${accountId ? `?account_id=${encodeURIComponent(accountId)}` : ""}`); if (!Array.isArray(payload.data)) throw new Error("The domains response was not in the expected format."); return payload.data.map((item) => ({ id: item.id, accountId: item.account_id, name: item.name, kind: item.kind, status: item.status, documentRoot: item.document_root, createdAt: item.created_at })); }
export async function listZones(accountId?: string) { const payload = await get<CollectionResponse<{ id: string; account_id: string; domain_id: string; name: string; status: string; created_at: string }>>(`/dns/zones${accountId ? `?account_id=${encodeURIComponent(accountId)}` : ""}`); if (!Array.isArray(payload.data)) throw new Error("The DNS zones response was not in the expected format."); return payload.data.map((item) => ({ id: item.id, accountId: item.account_id, domainId: item.domain_id, name: item.name, status: item.status, createdAt: item.created_at })); }
export async function listRecords(zoneId: string) { const payload = await get<CollectionResponse<{ id: string; zone_id: string; name: string; type: string; ttl: number; value: string }>>(`/dns/zones/${encodeURIComponent(zoneId)}/records`); if (!Array.isArray(payload.data)) throw new Error("The DNS records response was not in the expected format."); return payload.data.map((item) => ({ id: item.id, zoneId: item.zone_id, name: item.name, type: item.type, ttl: item.ttl, value: item.value })); }
export async function createDomain(body: { account_id: string; name: string }) { return post<{ data: Domain }>("/domains", body); }
export async function createZone(body: { account_id: string; domain_id: string }) { return post<{ data: DNSZone }>("/dns/zones", body); }
export async function createRecord(zoneId: string, body: { name: string; type: string; ttl: number; value: string }) { return post<{ data: DNSRecord }>(`/dns/zones/${encodeURIComponent(zoneId)}/records`, body); }

