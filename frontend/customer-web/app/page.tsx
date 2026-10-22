import { navigation } from "@harbor-panel/ui";

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

const accountStats = [
  ["Websites", "4", "3 online", "positive"],
  ["Databases", "4 / 10", "6 available", "neutral"],
  ["Mailboxes", "3 / 25", "22 available", "neutral"],
  ["SSL certificates", "3", "All valid", "positive"],
] as const;

const resourceUsage = [
  ["Disk space", "3.4 GB", "10 GB", 34, "blue"],
  ["Monthly bandwidth", "42.8 GB", "100 GB", 43, "violet"],
  ["Mailboxes", "3", "25", 12, "green"],
  ["Databases", "4", "10", 40, "amber"],
] as const;

const sites = [
  ["example.com", "Primary website", "Online", "Valid", "Today, 09:42"],
  ["store.example.com", "Online store", "Online", "Valid", "Yesterday, 18:15"],
  ["blog.example.com", "Blog", "Attention", "Valid", "Oct 20, 14:08"],
] as const;

const recentActivity = [
  ["Backup completed", "example.com · Today, 09:42", "success"],
  ["Certificate renewed", "store.example.com · Yesterday, 18:15", "success"],
  ["New mailbox created", "hello@example.com · Oct 20, 14:08", "info"],
] as const;

function ToolIcon({ label }: { label: string }) {
  return (
    <span className="tool-icon" aria-hidden="true">
      {iconByLabel[label] ?? "•"}
    </span>
  );
}

export default function CustomerDashboard() {
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
            <strong>Alex Morgan</strong>
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
          <span>All systems operational</span>
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
            <a className="user" href="/profile" aria-label="Open profile for admin@example.com">
              <span className="user-email">admin@example.com</span>
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
              <p className="eyebrow">Business plan · Account #HM-1048</p>
              <h1 id="page-title">Good morning, Alex</h1>
              <p className="subtitle">Manage your websites and hosting services from one place.</p>
            </div>
            <div className="welcome-actions">
              <span className="health-badge"><span className="status-dot" aria-hidden="true" /> Account healthy</span>
              <a className="primary-button" href="/domains/new"><span aria-hidden="true">+</span> Add domain</a>
            </div>
          </section>

          <section className="stats-grid" aria-label="Account summary">
            {accountStats.map(([label, value, detail, tone]) => (
              <div className="stat-card" key={label}>
                <div className="stat-label">{label}</div>
                <div className="stat-value-row">
                  <strong>{value}</strong>
                  <span className={`stat-indicator ${tone}`}>{tone === "positive" ? "●" : ""}</span>
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
              <div className="usage-list">
                {resourceUsage.map(([label, used, limit, percentage, tone]) => (
                  <div className="usage-row" key={label}>
                    <div className="usage-row-heading">
                      <span>{label}</span>
                      <strong>{used} <span className="usage-limit">of {limit}</span></strong>
                    </div>
                    <div className="progress-track" role="progressbar" aria-label={`${label} usage`} aria-valuemin={0} aria-valuemax={100} aria-valuenow={percentage}>
                      <span className={`progress-fill ${tone}`} style={{ width: `${percentage}%` }} />
                    </div>
                  </div>
                ))}
              </div>
              <div className="usage-note"><span className="status-dot" aria-hidden="true" /> You have room to grow on your current plan.</div>
            </section>

            <section className="card panel activity-panel" aria-labelledby="activity-title">
              <div className="panel-header">
                <div>
                  <h2 id="activity-title">Recent activity</h2>
                  <p>Latest changes to your account</p>
                </div>
                <a className="text-link" href="/activity" aria-label="View all account activity">View all</a>
              </div>
              <div className="activity-list">
                {recentActivity.map(([title, detail, tone]) => (
                  <div className="activity-item" key={title}>
                    <span className={`activity-icon ${tone}`} aria-hidden="true">{tone === "success" ? "✓" : "i"}</span>
                    <div><strong>{title}</strong><p>{detail}</p></div>
                  </div>
                ))}
              </div>
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
            <div className="table-wrap">
              <table>
                <caption className="sr-only">Websites and domain status</caption>
                <thead>
                  <tr>
                    <th scope="col">Domain</th>
                    <th scope="col">Site</th>
                    <th scope="col">Status</th>
                    <th scope="col">SSL/TLS</th>
                    <th scope="col">Last activity</th>
                    <th scope="col"><span className="sr-only">Actions</span></th>
                  </tr>
                </thead>
                <tbody>
                  {sites.map(([domain, site, status, ssl, updated]) => (
                    <tr key={domain}>
                      <td><a className="domain-link" href={`/domains/${domain}`}>{domain}</a></td>
                      <td className="table-muted">{site}</td>
                      <td><span className={`status-pill ${status === "Online" ? "online" : "attention"}`}><span aria-hidden="true">●</span> {status}</span></td>
                      <td><span className="ssl-status"><span aria-hidden="true">✓</span> {ssl}</span></td>
                      <td className="table-muted">{updated}</td>
                      <td><a className="row-action" href={`/domains/${domain}`}>Manage</a></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
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
