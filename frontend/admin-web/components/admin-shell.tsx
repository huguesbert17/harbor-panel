import { useState, type ReactNode } from "react";

import { Icon, type IconName } from "./icons";

type NavigationItem = { href: string; icon: IconName; label: string };
type NavigationGroup = { items: NavigationItem[]; label: string };

const navigationGroups: NavigationGroup[] = [
  { label: "Home", items: [{ href: "/", icon: "home", label: "Dashboard" }] },
  { label: "Server configuration", items: [{ href: "/servers", icon: "server", label: "Servers" }, { href: "/services", icon: "sliders", label: "Services" }, { href: "/settings", icon: "settings", label: "Settings" }] },
  { label: "Accounts", items: [{ href: "/accounts", icon: "users", label: "List accounts" }, { href: "/accounts/new", icon: "users", label: "Create account" }, { href: "/packages", icon: "package", label: "Packages" }] },
  { label: "Networking & security", items: [{ href: "/domains", icon: "globe", label: "Domains" }, { href: "/dns", icon: "dns", label: "DNS zones" }, { href: "/security", icon: "shield", label: "Security center" }] },
  { label: "Operations", items: [{ href: "/backups", icon: "backup", label: "Backups" }, { href: "/activity", icon: "activity", label: "Activity log" }] },
];

export function AdminShell({ children, activeHref = "/" }: { children: ReactNode; activeHref?: string }) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  return <div className={`admin-shell ${sidebarOpen ? "" : "admin-shell-sidebar-collapsed"}`}>
    <aside className="sidebar" aria-label="Administrator navigation">
      <div className="brand-lockup"><span className="brand-mark" aria-hidden="true">H</span><span><strong>HarborPanel</strong><small>Administrator</small></span></div>
      <div className="sidebar-controls"><button type="button" onClick={() => setSidebarOpen(true)}><span>↓</span> Expand</button><button type="button" onClick={() => setSidebarOpen(false)}><span>↑</span> Collapse</button></div>
      <label className="sidebar-search"><Icon name="search" size={16} /><span className="sr-only">Search tools</span><input aria-label="Search tools" placeholder="Search tools (Ctrl /)" /></label>
      <nav className="sidebar-nav">{navigationGroups.map((group) => <div className="nav-group" key={group.label}><p className="nav-group-label">{group.label}</p>{group.items.map((item) => { const isActive = item.href === activeHref; return <a aria-current={isActive ? "page" : undefined} className={`nav-item ${isActive ? "nav-item-active" : ""}`} href={item.href} key={item.href}><Icon name={item.icon} size={17} /><span>{item.label}</span></a>; })}</div>)}</nav>
      <div className="sidebar-footer"><a className="nav-item" href="/documentation"><Icon name="help" size={17} /><span>Documentation</span></a><a className="nav-item" href="/settings"><Icon name="settings" size={17} /><span>Settings</span></a><div className="sidebar-server"><span className="server-pulse" aria-hidden="true" /><div><strong>web-01</strong><span>Production</span></div><Icon name="chevron-down" size={15} /></div></div>
    </aside>
    <div className="admin-main">
      <header className="admin-header"><div className="server-facts"><div><span>Username</span><strong>root</strong></div><div><span>Hostname</span><strong>host.example.com</strong></div><div><span>Operating system</span><strong>AlmaLinux v9.8.0</strong></div><div><span>HarborPanel version</span><strong>0.1.0</strong></div><div className="server-facts-load"><span>Load averages</span><strong>↻ &nbsp;0.00&nbsp;&nbsp;0.05&nbsp;&nbsp;0.07</strong></div></div><div className="topbar"><div className="mobile-brand"><span className="brand-mark">H</span><strong>HarborPanel</strong></div><div className="topbar-context"><span className="context-dot" /> Production / web-01 <Icon name="chevron-down" size={14} /></div><div className="topbar-actions"><label className="search-field"><Icon name="search" size={17} /><span className="sr-only">Search tools and accounts</span><input aria-label="Search tools and accounts" placeholder="Search tools and accounts (/)" type="search" /></label><a className="icon-button" href="/activity" aria-label="View notifications and activity"><Icon name="activity" size={18} /></a><a className="user-menu" href="/settings"><span className="avatar">AD</span><span className="user-name">Admin</span><Icon name="chevron-down" size={14} /></a></div></div></header>
      <main>{children}</main>
    </div>
  </div>;
}
