"use client";

import { useEffect, useMemo, useState } from "react";
import { FiArchive, FiBarChart2, FiCode, FiDatabase, FiFile, FiGlobe, FiHardDrive, FiLock, FiMail, FiSettings, FiShield, FiTool, FiUsers } from "react-icons/fi";
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
    label: "Email",
    description: "Manage mailboxes and delivery.",
    icon: FiMail,
    tools: [
      ["Email Accounts", "Create and manage mailboxes.", "/email"], ["Forwarders", "Route mail to another address.", "/email/forwarders"], ["Email Routing", "Choose how mail is delivered.", "/email/routing"], ["Autoresponders", "Send automatic replies.", "/email/autoresponders"], ["Email Filters", "Organize incoming messages.", "/email/filters"], ["Spam Filters", "Protect your inbox from spam.", "/email/spam"],
    ],
  },
  {
    label: "Files",
    description: "Upload, restore, and protect your files.",
    icon: FiFile,
    tools: [
      ["File Manager", "Upload and organize website files.", "/files"], ["Disk Usage", "Review storage usage by directory.", "/files/disk-usage"], ["Backup", "Create or restore an account backup.", "/backups"], ["Git Version Control", "Deploy code from a repository.", "/files/git"],
    ],
  },
  {
    label: "Databases",
    description: "Create and manage application data.",
    icon: FiDatabase,
    tools: [
      ["Manage My Databases", "Create databases and users.", "/databases"], ["Database Wizard", "Set up a database step by step.", "/databases/wizard"], ["phpMyAdmin", "Work with database tables and queries.", "/databases/phpmyadmin"],
    ],
  },
  {
    label: "Domains",
    description: "Connect, route, and secure your websites.",
    icon: FiGlobe,
    tools: [
      ["Domains", "Manage domains, subdomains, and aliases.", "/domains"], ["Redirects", "Send visitors to another URL.", "/domains/redirects"], ["Zone Editor", "Control DNS records.", "/dns"], ["SSL/TLS", "Secure websites with certificates.", "/ssl"],
    ],
  },
  {
    label: "Metrics",
    description: "Understand traffic and resource usage.",
    icon: FiBarChart2,
    tools: [
      ["Visitors", "Review recent visitors.", "/metrics/visitors"], ["Errors", "Inspect recent website errors.", "/metrics/errors"], ["Bandwidth", "Review transfer usage.", "/metrics/bandwidth"], ["Raw Access", "Download raw access logs.", "/metrics/raw-access"],
    ],
  },
  {
    label: "Security",
    description: "Protect your account and websites.",
    icon: FiShield,
    tools: [
      ["SSL/TLS Certificates", "Install and manage certificates.", "/ssl"], ["SSH Access", "Manage secure shell access.", "/security/ssh"], ["IP Blocker", "Block unwanted visitors.", "/security/ip-blocker"], ["Manage API Tokens", "Control programmatic access.", "/security/api-tokens"],
    ],
  },
  {
    label: "Software",
    description: "Configure runtimes and applications.",
    icon: FiCode,
    tools: [
      ["MultiPHP Manager", "Choose the PHP version for domains.", "/software/php"], ["MultiPHP INI Editor", "Tune PHP runtime settings.", "/software/php-ini"], ["Optimize Website", "Improve content delivery.", "/software/optimize"],
    ],
  },
  {
    label: "Preferences",
    description: "Customize your account experience.",
    icon: FiSettings,
    tools: [
      ["Account Preferences", "Manage account settings.", "/profile"], ["Password & Security", "Update account security.", "/profile/security"], ["Contact Information", "Keep contact details current.", "/profile/contact"], ["User Manager", "Manage additional users.", "/profile/users"],
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
  const Icon = toolIconByLabel[label] ?? FiTool;
  return (
    <span className="tool-icon" aria-hidden="true">
      <Icon />
    </span>
  );
}

const toolIconByLabel: Record<string, typeof FiTool> = {
  "Email Accounts": FiUsers, Forwarders: FiMail, "Email Routing": FiSettings, Autoresponders: FiMail, "Email Filters": FiTool, "Spam Filters": FiShield,
  "File Manager": FiFile, "Disk Usage": FiHardDrive, Backup: FiArchive, "Git Version Control": FiCode,
  "Manage My Databases": FiDatabase, "Database Wizard": FiSettings, phpMyAdmin: FiDatabase,
  Domains: FiGlobe, Redirects: FiGlobe, "Zone Editor": FiGlobe, "SSL/TLS": FiLock, Visitors: FiUsers, Errors: FiTool, Bandwidth: FiBarChart2, "Raw Access": FiFile,
  "SSL/TLS Certificates": FiLock, "SSH Access": FiCode, "IP Blocker": FiShield, "Manage API Tokens": FiTool, "MultiPHP Manager": FiCode, "MultiPHP INI Editor": FiSettings, "Optimize Website": FiBarChart2,
  "Account Preferences": FiSettings, "Password & Security": FiLock, "Contact Information": FiUsers, "User Manager": FiUsers,
};

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

          <div className="customer-dashboard-layout">
          <div className="tool-groups">
            {toolGroups.map((group) => (
              <section className="tool-group" key={group.label} aria-labelledby={`${group.label}-title`}>
                <div className="tool-group-header">
                  <h3 id={`${group.label}-title`}><group.icon aria-hidden="true" /> {group.label}</h3>
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
          <aside className="account-rail" aria-label="Account information">
            <section className="card account-panel">
              <h2>General Information</h2>
              <dl>
                <div><dt>Current User</dt><dd>{account?.username ?? "example"}</dd></div>
                <div><dt>Primary Domain</dt><dd><a href="/domains">example.com ↗</a></dd></div>
                <div><dt>SSL Certificate</dt><dd className="rail-warning"><FiShield /> Self-signed <a href="/ssl">Add SSL Certificate</a></dd></div>
                <div><dt>Shared IP Address</dt><dd>192.0.2.10</dd></div>
                <div><dt>Home Directory</dt><dd>/home/{account?.username ?? "example"}</dd></div>
              </dl>
            </section>
            <section className="card account-panel">
              <h2>Statistics</h2>
              <dl>
                <div><dt>Disk Usage</dt><dd>1.95 MB / ∞</dd></div>
                <div><dt>Bandwidth</dt><dd>56.67 KB / ∞</dd></div>
                <div><dt>Email Accounts</dt><dd>5 / ∞</dd></div>
                <div><dt>Databases</dt><dd>2 / ∞</dd></div>
              </dl>
            </section>
          </aside>
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
