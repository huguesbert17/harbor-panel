"use client";

import { useEffect, useMemo, useState } from "react";
import { navigation } from "@harbor-panel/ui";
import { CustomerDashboardData, getCustomerDashboardData, HostingAccount } from "./lib/customer-api";

type CustomerNavItem = (typeof navigation.customer)[number];

const iconByLabel: Record<string, string> = {
  Dashboard: "⌂",
  Domains: "◈",
  "File Manager": "▤",
  Databases: "◆",
  Email: "✉",
  "DNS Zone Editor": "⌁",
  "SSL/TLS": "✓",
  Backups: "↻",
};

const navGroups: { label: string; items: readonly CustomerNavItem[] }[] = [
  { label: "Overview", items: [navigation.customer[0]] },
  {
    label: "Websites & domains",
    items: [navigation.customer[1], navigation.customer[2], navigation.customer[5], navigation.customer[6]],
  },
  { label: "Applications & mail", items: [navigation.customer[3], navigation.customer[4], navigation.customer[7]] },
];

const toolGroups = [
  {
    label: "Websites & domains",
    description: "Build, connect, and secure your sites.",
    tools: [
      ["Domains", "Manage domains, subdomains, and redirects.", "/domains"],
      ["File Manager", "Upload and organize your website files.", "/files"],
      ["DNS Zone Editor", "Control records for your domains.", "/dns"],
      ["SSL/TLS", "Secure your websites with certificates.", "/ssl"],
    ],
  },
  {
    label: "Applications & data",
    description: "Manage the services your applications depend on.",
    tools: [
      ["Databases", "Create databases and database users.", "/databases"],
      ["Email", "Manage mailboxes, aliases, and forwarding.", "/email"],
      ["Backups", "Create and restore account backups.", "/backups"],
    ],
  },
] as const;

type DashboardState =
  | { status: "loading" }
  | { status: "ready"; data: CustomerDashboardData }
  | { status: "error"; message: string };

function accountStatusTone(account?: HostingAccount) {
  if (!account) return "neutral";
  return account.status === "active" ? "positive" : "warning";
}

function accountStatusLabel(account?: HostingAccount) {
  if (!account) return "Not assigned";
  return account.status.charAt(0).toUpperCase() + account.status.slice(1);
}

function formatServerStatus(data: CustomerDashboardData) {
  if (data.servers.length === 0) return "Awaiting telemetry";
  const connected = data.servers.filter((server) => server.connectionStatus === "connected").length;
  return `${connected} of ${data.servers.length} connected`;
}

function formatApiError(message: string) {
  return message.includes("fetch") || message.includes("Failed")
    ? "The control-plane API could not be reached. Check the API URL and try again."
    : message;
}

function ToolIcon({ label }: { label: string }) {
  return (
    <span className="tool-icon" aria-hidden="true">
      {iconByLabel[label] ?? "•"}
    </span>
  );
}

export default function CustomerDashboard() {
  const [dashboardState, setDashboardState] = useState<DashboardState>({ status: "loading" });

  useEffect(() => {
    const controller = new AbortController();
    getCustomerDashboardData(controller.signal)
      .then((data) => setDashboardState({ status: "ready", data }))
      .catch((error: unknown) => {
        if (!controller.signal.aborted) {
          setDashboardState({ status: "error", message: error instanceof Error ? error.message : "Unable to load dashboard data." });
        }
      });

    return () => controller.abort();
  }, []);

  const account = dashboardState.status === "ready" ? dashboardState.data.accounts[0] : undefined;
  const servers = dashboardState.status === "ready" ? dashboardState.data.servers : [];
  const availableServices = useMemo(
    () => servers.reduce((total, server) => total + server.capabilities.filter((capability) => capability.status === "available").length, 0),
    [servers],
  );
  const accountStats = [
    ["Hosting account", dashboardState.status === "loading" ? "…" : account ? "1" : "—", account ? account.username : "No account returned", accountStatusTone(account)],
    ["Account status", dashboardState.status === "loading" ? "…" : accountStatusLabel(account), account ? `Account ID ${account.id}` : "Connect an account to continue", accountStatusTone(account)],
    ["Managed servers", dashboardState.status === "loading" ? "…" : String(servers.length), dashboardState.status === "ready" ? formatServerStatus(dashboardState.data) : "Loading server status", servers.length > 0 ? "positive" : "neutral"],
    ["Available services", dashboardState.status === "loading" ? "…" : String(availableServices), dashboardState.status === "ready" ? "From server capabilities" : "Loading capabilities", availableServices > 0 ? "positive" : "neutral"],
  ] as const;
  const healthLabel = dashboardState.status === "loading"
    ? "Checking account health"
    : dashboardState.status === "error"
      ? "API unavailable"
      : account && account.status !== "active"
        ? `Account ${account.status}`
        : servers.length === 0
          ? "Awaiting server telemetry"
          : "Account healthy";
  const healthTone = dashboardState.status === "error" || (account && account.status !== "active") ? "warning" : "positive";

  return (
    <div className="shell">
      <a className="skip-link" href="#main-content">
        Skip to main content
      </a>

      <aside className="sidebar" aria-label="Customer account navigation">
        <div className="brand-block">
          <a className="brand" href="/" aria-label="HarborPanel customer dashboard">
            <span className="brand-mark" aria-hidden="true">
              H
            </span>
            <span>HarborPanel</span>
          </a>
          <div className="account-switcher">
            <span className="account-switcher-label">Account</span>
            <strong>{account?.username ?? (dashboardState.status === "loading" ? "Loading account…" : "Customer account")}</strong>
            <span className="account-switcher-caret" aria-hidden="true">
              ⌄
            </span>
          </div>
        </div>

        <nav className="sidebar-nav">
          {navGroups.map((group) => (
            <div className="nav-group" key={group.label}>
              <p className="section-label">{group.label}</p>
              {group.items.map((item) => (
                <a
                  className={`nav-item ${item.href === "/" ? "active" : ""}`}
                  href={item.href}
                  key={item.href}
                  aria-current={item.href === "/" ? "page" : undefined}
                >
                  <span className="nav-icon" aria-hidden="true">
                    {iconByLabel[item.label]}
                  </span>
                  <span>{item.label}</span>
                </a>
              ))}
            </div>
          ))}
          <div className="nav-group nav-group-support">
            <p className="section-label">Support</p>
            <a className="nav-item" href="/documentation">
              <span className="nav-icon" aria-hidden="true">?</span>
              <span>Documentation</span>
            </a>
            <a className="nav-item" href="/support">
              <span className="nav-icon" aria-hidden="true">↗</span>
              <span>Contact support</span>
            </a>
          </div>
        </nav>
        <div className="sidebar-footer">
          <span className="status-dot" aria-hidden="true" />
          <span>{dashboardState.status === "ready" ? formatServerStatus(dashboardState.data) : dashboardState.status === "error" ? "API unavailable" : "Checking systems…"}</span>
        </div>
      </aside>

      <main className="main" id="main-content">
        <header className="topbar">
          <div className="breadcrumb" aria-label="Breadcrumb">
            <span className="breadcrumb-muted">Customer account</span>
            <span aria-hidden="true">/</span>
            <strong>Dashboard</strong>
          </div>
          <div className="topbar-actions">
            <a className="topbar-link" href="/support">Help center</a>
            <span className="user-divider" aria-hidden="true" />
            <a className="user" href="/profile" aria-label="Open customer profile">
              <span className="user-email">{account?.username ?? "customer account"}</span>
              <span className="avatar" aria-hidden="true">AM</span>
            </a>
          </div>
        </header>

        <nav className="mobile-nav" aria-label="Mobile account navigation">
          {navigation.customer.slice(0, 5).map((item) => (
            <a className={item.href === "/" ? "active" : ""} href={item.href} key={item.href} aria-current={item.href === "/" ? "page" : undefined}>
              <span aria-hidden="true">{iconByLabel[item.label]}</span>
              {item.label}
            </a>
          ))}
        </nav>

        <div className="content">
          <section className="welcome-row" aria-labelledby="page-title">
            <div>
              <p className="eyebrow">Customer account{account ? ` · ${account.id}` : ""}</p>
              <h1 id="page-title">Welcome to your hosting account</h1>
              <p className="subtitle">Manage your websites and hosting services from one place.</p>
            </div>
            <div className="welcome-actions">
              <span className={`health-badge ${healthTone}`} aria-live="polite"><span className="status-dot" aria-hidden="true" /> {healthLabel}</span>
              <a className="primary-button" href="/domains/new"><span aria-hidden="true">+</span> Add domain</a>
            </div>
          </section>

          <section className="stats-grid" aria-label="Account summary">
            {accountStats.map(([label, value, detail, tone]) => (
              <div className="stat-card" key={label}>
                <div className="stat-label">{label}</div>
                <div className="stat-value-row">
                  <strong>{value}</strong>
                  <span className={`stat-indicator ${tone}`}>{tone === "positive" ? "●" : tone === "warning" ? "!" : ""}</span>
                </div>
                <div className="stat-detail">{detail}</div>
              </div>
            ))}
          </section>

          <section className="section-heading" id="tools" aria-labelledby="tools-title">
            <div>
              <h2 id="tools-title">Quick access</h2>
              <p>Common tools for your account</p>
            </div>
            <a className="text-link" href="/features">View all features <span aria-hidden="true">→</span></a>
          </section>

          <div className="tool-groups">
            {toolGroups.map((group) => (
              <section className="tool-group" key={group.label} aria-labelledby={`${group.label}-title`}>
                <div className="tool-group-header">
                  <h3 id={`${group.label}-title`}>{group.label}</h3>
                  <p>{group.description}</p>
                </div>
                <div className="tools-grid">
                  {group.tools.map(([label, copy, href]) => (
                    <a className="card tool" href={href} key={label}>
                      <ToolIcon label={label} />
                      <span className="tool-content">
                        <span className="tool-title">{label}</span>
                        <span className="tool-copy">{copy}</span>
                      </span>
                      <span className="tool-arrow" aria-hidden="true">→</span>
                    </a>
                  ))}
                </div>
              </section>
            ))}
          </div>

          <div className="dashboard-columns">
            <section className="card panel resource-panel" aria-labelledby="usage-title">
              <div className="panel-header">
                <div>
                  <h2 id="usage-title">Resource usage</h2>
                  <p>Usage resets on the first day of each month.</p>
                </div>
                <a className="text-link" href="/usage">View details</a>
              </div>
              {dashboardState.status === "loading" ? (
                <div className="state-card loading-state" role="status">Loading usage and quota information…</div>
              ) : dashboardState.status === "error" ? (
                <div className="state-card error-state" role="alert">{formatApiError(dashboardState.message)}</div>
              ) : (
                <div className="state-card empty-state">
                  <strong>{account ? "Usage data is not connected yet" : "No hosting account is assigned"}</strong>
                  <p>{account ? "Quota and bandwidth details will appear when the usage API is available." : "Ask an administrator to assign a hosting account before managing services."}</p>
                </div>
              )}
              <div className="usage-note"><span className="status-dot" aria-hidden="true" /> Account limits will be shown here when available.</div>
            </section>

            <section className="card panel activity-panel" aria-labelledby="activity-title">
              <div className="panel-header">
                <div>
                  <h2 id="activity-title">Recent activity</h2>
                  <p>Latest changes to your account</p>
                </div>
                <a className="text-link" href="/activity" aria-label="View all account activity">View all</a>
              </div>
              {dashboardState.status === "loading" ? (
                <div className="state-card loading-state" role="status">Loading recent activity…</div>
              ) : dashboardState.status === "error" ? (
                <div className="state-card error-state" role="alert">Activity is unavailable while the API is offline.</div>
              ) : (
                <div className="state-card empty-state"><strong>No recent activity</strong><p>The activity feed will appear when account events are available.</p></div>
              )}
            </section>
          </div>

          <section className="card panel domains-panel" aria-labelledby="sites-title">
            <div className="panel-header">
              <div>
                <h2 id="sites-title">Websites &amp; domains</h2>
                <p>Monitor the sites connected to your account.</p>
              </div>
              <a className="text-link" href="/domains">Manage domains <span aria-hidden="true">→</span></a>
            </div>
            {dashboardState.status === "loading" ? (
              <div className="state-card loading-state" role="status">Loading domains…</div>
            ) : dashboardState.status === "error" ? (
              <div className="state-card error-state" role="alert">Domains are unavailable while the API is offline.</div>
            ) : (
              <div className="state-card empty-state"><strong>No domains to display</strong><p>The domains API is not connected yet. Use the Domains tool to manage this area when it becomes available.</p></div>
            )}
          </section>

          <footer className="page-footer">
            <span>HarborPanel · Business plan</span>
            <span><a href="/status">System status</a> · <a href="/support">Support</a></span>
          </footer>
        </div>
      </main>
    </div>
  );
}
