"use client";

import { useEffect, useState } from "react";

import { listHostingAccounts, listManagedServers, type HostingAccount, type ManagedServer } from "../lib/api-client";
import { AdminShell } from "../components/admin-shell";
import { Icon, type IconName } from "../components/icons";
import { LinkButton, Panel, ProgressBar, StatusBadge } from "../components/ui";

type DashboardState = { accounts: HostingAccount[]; servers: ManagedServer[] };

const favorites = [
  ["List Accounts", "/accounts", "users", "This interface lists your server’s accounts and lets you perform certain actions on them."],
  ["Process Manager", "/services", "activity", "View and manage your server’s active processes."],
  ["Create a New Account", "/accounts/new", "users", "Create a new hosting account on this server."],
  ["DNS Zone Manager", "/dns", "dns", "Manage DNS zones and records for your accounts."],
  ["Terminal", "/terminal", "terminal", "Open an in-browser terminal for server administration."],
  ["Mail Queue Manager", "/mail/queue", "settings", "Review and manage queued mail delivery."],
];

const toolGroups = [
  { label: "Server Configuration", tools: ["Basic WebHost Manager Setup", "Change Root Password", "Server Profile", "Server Time", "Terminal", "Tweak Settings", "Update Preferences"] },
  { label: "Support", tools: ["Contact Support", "Documentation", "Support Center"] },
  { label: "Networking Setup", tools: ["Nameserver Configuration", "Network Setup", "IP Functions", "Hostname Configuration"] },
  { label: "Security Center", tools: ["Security Advisor", "SSL/TLS Manager", "Firewall", "Access Control"] },
  { label: "Service Configuration", tools: ["Service Manager", "Service Configuration", "Mail Queue Manager", "PHP-FPM Manager"] },
  { label: "Backup", tools: ["Backup Configuration", "Backup Restoration", "Backup Destination"] },
  { label: "Server Status", tools: ["Service Status", "Process Manager", "System Health", "Daily Process Log"] },
  { label: "Account Information", tools: ["List Accounts", "Account Information", "Account Statistics", "Account Transfers"] },
  { label: "Account Functions", tools: ["Create a New Account", "Terminate an Account", "Modify an Account", "Manage Account Suspensions"] },
  { label: "Multi Account Functions", tools: ["Modify Multiple Accounts", "Change Multiple Sites", "Manage Account Upgrades"] },
  { label: "Packages", tools: ["Create Package", "Feature Manager", "Manage Packages", "Reseller Packages"] },
  { label: "DNS Functions", tools: ["DNS Zone Manager", "Add an A Entry", "Edit DNS Zone", "Synchronize DNS Records"] },
  { label: "Database Services", tools: ["Manage Database Services", "Configure PostgreSQL", "Configure MySQL"] },
  { label: "Email", tools: ["Mail Queue Manager", "Mail Delivery Reports", "Mail Configurations"] },
  { label: "SSL/TLS", tools: ["Install an SSL Certificate", "Manage SSL Hosts", "SSL Storage Manager"] },
];

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
  const [openTools, setOpenTools] = useState<Record<string, boolean>>({ "Server Configuration": true });

  return <AdminShell>
    <div className="page-content">
      {isLoading ? <DashboardLoading /> : error ? <DashboardError message={error} onRetry={() => setRequestKey((value) => value + 1)} /> : <>
        <section className="whm-notice" aria-label="System notice"><div className="notice-icon"><Icon name="activity" size={17} /></div><div><strong>Welcome to your control plane</strong><span>Use search, favorites, or the tool groups below to manage this server.</span></div><button className="notice-dismiss" type="button" aria-label="Dismiss notice">Dismiss</button></section>
        <div className="dashboard-grid dashboard-grid-launchpad">
          <Panel eyebrow="Quick access" title="Favorites"><div className="favorites-grid">{favorites.map(([label, href, icon, description]) => <a className="favorite-card" href={href} key={label}><span className="favorite-icon"><Icon name={icon as IconName} size={17} /></span><span className="favorite-copy"><strong>{label}</strong><small>{description}</small></span><Icon name="arrow-up-right" size={14} /></a>)}</div></Panel>
          <Panel eyebrow="Server statistics" title="Statistics"><div className="server-stat-list"><div><span>IPv4 address</span><strong>192.0.2.10</strong></div><div><span>User analytics ID</span><strong>Not enrolled</strong></div><div><span>WP Toolkit</span><strong className="stat-health"><span className="server-status-dot" /> No security issues found</strong></div></div></Panel>
        </div>
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

        <Panel className="tools-panel" eyebrow="Administration catalog" title="Tools">
          <div className="tool-groups">{toolGroups.map((group) => { const isOpen = openTools[group.label] ?? false; return <div className="tool-group" key={group.label}><button className="tool-group-toggle" aria-expanded={isOpen} onClick={() => setOpenTools((current) => ({ ...current, [group.label]: !isOpen }))} type="button"><span><Icon name="chevron-down" size={15} />{group.label}</span><small>{group.tools.length} tools</small></button>{isOpen ? <div className="tool-items">{group.tools.map((tool) => <a href="#" key={tool}><Icon name="arrow-up-right" size={13} />{tool}</a>)}</div> : null}</div>; })}</div>
        </Panel>
      </>}
    </div>
  </AdminShell>;
}
