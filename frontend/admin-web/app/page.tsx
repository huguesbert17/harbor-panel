import { AdminShell } from "../components/admin-shell";
import { Icon } from "../components/icons";
import { LinkButton, MetricCard, Panel, ProgressBar, StatusBadge } from "../components/ui";

const accounts = [
  { account: "acme-web", domain: "acme.example", plan: "Business", status: "Active", usage: "2.4 GB", created: "Today, 09:42" },
  { account: "northstar", domain: "northstar.example", plan: "Starter", status: "Active", usage: "814 MB", created: "Yesterday" },
  { account: "studio-lumen", domain: "studio-lumen.example", plan: "Business", status: "Pending", usage: "—", created: "Oct 21, 16:18" },
  { account: "brightleaf", domain: "brightleaf.example", plan: "Pro", status: "Suspended", usage: "6.1 GB", created: "Oct 20, 11:07" },
];

const activities = [
  { icon: "circle-check" as const, title: "Certificate issued", detail: "acme.example", time: "8 min ago", tone: "success" },
  { icon: "backup" as const, title: "Backup completed", detail: "northstar", time: "42 min ago", tone: "info" },
  { icon: "users" as const, title: "Hosting account created", detail: "studio-lumen", time: "1 hr ago", tone: "blue" },
  { icon: "sliders" as const, title: "PHP runtime updated", detail: "web-01 · PHP 8.3", time: "3 hr ago", tone: "orange" },
];

export default function AdminDashboard() {
  return <AdminShell>
    <div className="page-content">
      <div className="page-heading">
        <div><div className="breadcrumb"><span>Administration</span><Icon name="chevron-down" size={13} /><strong>Dashboard</strong></div><h1>Server overview</h1><p className="page-subtitle">Manage your hosting platform, accounts, and infrastructure from one place.</p></div>
        <div className="heading-actions"><a className="button button-secondary" href="/servers"><Icon name="server" size={16} /> View server</a><a className="button button-primary" href="/accounts/new"><span className="button-plus">+</span> Create account</a></div>
      </div>

      <section aria-label="Platform summary" className="metrics-grid">
        <MetricCard detail="8 added this month" icon="users" label="Hosting accounts" tone="blue" value="128" />
        <MetricCard detail="12 added this month" icon="globe" label="Active domains" tone="green" value="342" />
        <MetricCard detail="34% of 1.2 TB capacity" icon="database" label="Storage used" tone="purple" value="418 GB" />
        <MetricCard detail="All services operational" icon="activity" label="Server health" tone="orange" value="99.98%" />
      </section>

      <div className="dashboard-grid dashboard-grid-primary">
        <Panel action={<LinkButton href="/accounts">View all accounts</LinkButton>} className="accounts-panel" eyebrow="Account management" title="Recent hosting accounts">
          <div className="table-wrap"><table className="data-table"><caption className="sr-only">Recent hosting accounts</caption><thead><tr><th scope="col">Account</th><th scope="col">Primary domain</th><th scope="col">Plan</th><th scope="col">Status</th><th scope="col">Usage</th><th scope="col">Created</th></tr></thead><tbody>{accounts.map((item) => <tr key={item.account}><th scope="row"><a className="table-primary-link" href={`/accounts/${item.account}`}>{item.account}</a></th><td>{item.domain}</td><td>{item.plan}</td><td><StatusBadge label={item.status} tone={item.status === "Active" ? "success" : item.status === "Pending" ? "warning" : "danger"} /></td><td>{item.usage}</td><td className="table-muted">{item.created}</td></tr>)}</tbody></table></div>
        </Panel>

        <Panel action={<LinkButton href="/activity">All activity</LinkButton>} className="activity-panel" eyebrow="Audit trail" title="Recent activity">
          <div className="activity-list">{activities.map((item) => <div className="activity-item" key={`${item.title}-${item.detail}`}><div className={`activity-icon activity-icon-${item.tone}`}><Icon name={item.icon} size={16} /></div><div className="activity-copy"><strong>{item.title}</strong><span>{item.detail}</span></div><time>{item.time}</time></div>)}</div>
        </Panel>
      </div>

      <div className="dashboard-grid dashboard-grid-secondary">
        <Panel action={<LinkButton href="/servers/web-01">Manage server</LinkButton>} eyebrow="Infrastructure" title="Server resources">
          <div className="resource-body"><div className="resource-intro"><div className="resource-name"><span className="server-status-dot" /> <strong>web-01</strong><StatusBadge label="Healthy" tone="success" /></div><span className="resource-meta">Ubuntu 24.04 LTS · 4 vCPU</span></div><div className="resource-bars"><ProgressBar detail="34% · 418 GB of 1.2 TB" label="Storage" value={34} tone="blue" /><ProgressBar detail="58% · 2.3 GB of 4 GB" label="Memory" value={58} tone="green" /><ProgressBar detail="22% · 1.1 of 4 cores" label="CPU load" value={22} tone="orange" /></div><div className="service-strip"><span><i className="service-dot" /> Nginx</span><span><i className="service-dot" /> PHP-FPM</span><span><i className="service-dot" /> MariaDB</span><span><i className="service-dot" /> PowerDNS</span></div></div>
        </Panel>

        <Panel eyebrow="Operations" title="Needs attention">
          <div className="attention-list"><a className="attention-item" href="/certificates"><span className="attention-icon attention-icon-warning"><Icon name="clock" size={16} /></span><span><strong>3 certificates expire soon</strong><small>Review certificates expiring within 30 days</small></span><Icon name="arrow-up-right" size={16} /></a><a className="attention-item" href="/jobs"><span className="attention-icon attention-icon-blue"><Icon name="activity" size={16} /></span><span><strong>2 provisioning jobs queued</strong><small>Account and certificate operations are waiting</small></span><Icon name="arrow-up-right" size={16} /></a><a className="attention-item" href="/packages"><span className="attention-icon attention-icon-danger"><Icon name="circle-info" size={16} /></span><span><strong>1 package near its limit</strong><small>Business plan is at 88% account allocation</small></span><Icon name="arrow-up-right" size={16} /></a></div>
        </Panel>
      </div>
    </div>
  </AdminShell>;
}
