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

