"use client";

import { useEffect, useState } from "react";

import { listHostingAccounts, listManagedServers, type HostingAccount, type ManagedServer } from "../lib/api-client";
import { AdminShell } from "../components/admin-shell";
import { Icon } from "../components/icons";
import { LinkButton, MetricCard, Panel, ProgressBar, StatusBadge } from "../components/ui";

type DashboardState = { accounts: HostingAccount[]; servers: ManagedServer[] };

function statusTone(status: string) {
  if (["active", "connected", "healthy"].includes(status.toLowerCase())) return "success" as const;
  if (["pending", "enrolling", "unknown"].includes(status.toLowerCase())) return "warning" as const;
  if (["suspended", "offline", "revoked", "degraded"].includes(status.toLowerCase())) return "danger" as const;
  return "neutral" as const;
}

function displayStatus(status: string) {
  return status.replace(/[_-]+/g, " ").replace(/\b\w/g, (character) => character.toUpperCase());
}

function formatObservedAt(value?: string) {
  if (!value) return "Not reported";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "Not reported" : date.toLocaleString();
}

function DashboardLoading() {
  return <div className="dashboard-state" role="status"><Icon name="activity" size={20} /><strong>Loading control-plane data…</strong><span>Fetching hosting accounts and managed servers.</span></div>;
}

function DashboardError({ message, onRetry }: { message: string; onRetry: () => void }) {
  return <div className="dashboard-state dashboard-state-error" role="alert"><Icon name="x-circle" size={20} /><strong>Unable to load the administration overview</strong><span>{message}</span><button className="button button-secondary" onClick={onRetry} type="button">Try again</button></div>;
}

export default function AdminDashboard() {
  const [dashboard, setDashboard] = useState<DashboardState | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [requestKey, setRequestKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    setIsLoading(true);
    setError(null);

    Promise.all([listHostingAccounts(), listManagedServers()])
      .then(([accounts, servers]) => {
        if (!cancelled) setDashboard({ accounts, servers });
      })
      .catch((requestError: unknown) => {
        if (!cancelled) setError(requestError instanceof Error ? requestError.message : "An unexpected API error occurred.");
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false);
      });

    return () => { cancelled = true; };
  }, [requestKey]);

  const accounts = dashboard?.accounts ?? [];
  const servers = dashboard?.servers ?? [];
  const primaryServer = servers[0];
  const activeAccounts = accounts.filter((account) => account.status === "active").length;
  const connectedServers = servers.filter((server) => server.connectionStatus === "connected").length;

  return <AdminShell>
    <div className="page-content">
      <div className="page-heading">
        <div><div className="breadcrumb"><span>Administration</span><Icon name="chevron-down" size={13} /><strong>Dashboard</strong></div><h1>Server overview</h1><p className="page-subtitle">Manage your hosting platform, accounts, and infrastructure from one place.</p></div>
        <div className="heading-actions"><a className="button button-secondary" href="/servers"><Icon name="server" size={16} /> View server</a><a className="button button-primary" href="/accounts/new"><span className="button-plus">+</span> Create account</a></div>
      </div>

      {isLoading ? <DashboardLoading /> : error ? <DashboardError message={error} onRetry={() => setRequestKey((value) => value + 1)} /> : <>
        <section aria-label="Platform summary" className="metrics-grid">
          <MetricCard detail={`${activeAccounts} active now`} icon="users" label="Hosting accounts" tone="blue" value={String(accounts.length)} />
          <MetricCard detail="Not reported by API yet" icon="globe" label="Active domains" tone="green" value="—" />
          <MetricCard detail="Not reported by API yet" icon="database" label="Storage used" tone="purple" value="—" />
          <MetricCard detail={`${connectedServers} connected`} icon="activity" label="Managed servers" tone="orange" value={String(servers.length)} />
        </section>

        <div className="dashboard-grid dashboard-grid-primary">
        <Panel action={<LinkButton href="/accounts">View all accounts</LinkButton>} className="accounts-panel" eyebrow="Account management" title="Recent hosting accounts">
          {accounts.length === 0 ? <div className="panel-empty"><Icon name="users" size={20} /><strong>No hosting accounts yet</strong><span>Accounts created through the control plane will appear here.</span></div> : <div className="table-wrap"><table className="data-table"><caption className="sr-only">Recent hosting accounts</caption><thead><tr><th scope="col">Account</th><th scope="col">Primary domain</th><th scope="col">Plan</th><th scope="col">Status</th><th scope="col">Usage</th><th scope="col">Created</th></tr></thead><tbody>{accounts.map((item) => <tr key={item.id}><th scope="row"><a className="table-primary-link" href={`/accounts/${item.id}`}>{item.username}</a></th><td className="table-muted">Not reported</td><td className="table-muted">Not reported</td><td><StatusBadge label={displayStatus(item.status)} tone={statusTone(item.status)} /></td><td className="table-muted">Not reported</td><td className="table-muted">Not reported</td></tr>)}</tbody></table></div>}
        </Panel>

        <Panel action={<LinkButton href="/servers">View servers</LinkButton>} className="activity-panel" eyebrow="Control plane" title="Managed servers">
          {servers.length === 0 ? <div className="panel-empty"><Icon name="server" size={20} /><strong>No managed servers yet</strong><span>Register a server to begin managing infrastructure.</span></div> : <div className="activity-list">{servers.map((server) => <div className="activity-item" key={server.serverId}><div className={`activity-icon activity-icon-${statusTone(server.connectionStatus)}`}><Icon name="server" size={16} /></div><div className="activity-copy"><strong>{server.serverId}</strong><span>{displayStatus(server.registrationStatus)} · {server.capabilities.length} capabilities</span></div><time>{formatObservedAt(server.lastSeenAt ?? server.observedAt)}</time></div>)}</div>}
        </Panel>
        </div>

        <div className="dashboard-grid dashboard-grid-secondary">
        <Panel action={<LinkButton href="/servers/web-01">Manage server</LinkButton>} eyebrow="Infrastructure" title="Server resources">
          <div className="resource-body"><div className="resource-intro"><div className="resource-name"><span className={`server-status-dot ${primaryServer?.connectionStatus === "connected" ? "" : "server-status-dot-muted"}`} /> <strong>{primaryServer?.serverId ?? "No server selected"}</strong>{primaryServer ? <StatusBadge label={displayStatus(primaryServer.connectionStatus)} tone={statusTone(primaryServer.connectionStatus)} /> : null}</div><span className="resource-meta">Resource metrics are not reported by the current API</span></div>{primaryServer ? <><div className="resource-bars"><ProgressBar detail="Not reported" label="Storage" value={0} tone="blue" /><ProgressBar detail="Not reported" label="Memory" value={0} tone="green" /><ProgressBar detail="Not reported" label="CPU load" value={0} tone="orange" /></div><div className="service-strip">{primaryServer.capabilities.length === 0 ? <span>Capabilities not reported</span> : primaryServer.capabilities.map((capability) => <span key={capability.name}><i className={`service-dot ${capability.status !== "available" ? "service-dot-muted" : ""}`} /> {displayStatus(capability.name)}</span>)}</div></> : <div className="panel-empty panel-empty-inline"><span>Register a managed server to see its health and capabilities.</span></div>}</div>
        </Panel>

        <Panel eyebrow="Operations" title="Needs attention">
          <div className="attention-list"><div className="attention-item"><span className="attention-icon attention-icon-blue"><Icon name="activity" size={16} /></span><span><strong>Live operational alerts are coming soon</strong><small>Jobs, certificates, and package thresholds will appear here as their modules are connected.</small></span></div></div>
        </Panel>
        </div>
      </>}
    </div>
  </AdminShell>;
}
