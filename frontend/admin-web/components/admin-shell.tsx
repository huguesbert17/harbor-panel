import type { ReactNode } from "react";

import { Icon, type IconName } from "./icons";

type NavigationItem = { href: string; icon: IconName; label: string };
type NavigationGroup = { items: NavigationItem[]; label: string };

const navigationGroups: NavigationGroup[] = [
  { label: "Overview", items: [{ href: "/", icon: "home", label: "Dashboard" }] },
  { label: "Hosting", items: [{ href: "/accounts", icon: "users", label: "Accounts" }, { href: "/packages", icon: "package", label: "Packages" }, { href: "/domains", icon: "globe", label: "Domains" }, { href: "/dns", icon: "dns", label: "DNS zones" }] },
  { label: "Infrastructure", items: [{ href: "/servers", icon: "server", label: "Servers" }, { href: "/services", icon: "sliders", label: "Services" }, { href: "/backups", icon: "backup", label: "Backups" }] },
  { label: "Security & access", items: [{ href: "/security", icon: "shield", label: "Security" }, { href: "/activity", icon: "activity", label: "Activity log" }] },
];

export function AdminShell({ children, activeHref = "/" }: { children: ReactNode; activeHref?: string }) {
  return <div className="admin-shell">
    <aside className="sidebar" aria-label="Administrator navigation">
      <div className="brand-lockup"><span className="brand-mark" aria-hidden="true">H</span><span><strong>HarborPanel</strong><small>Administrator</small></span></div>
      <nav className="sidebar-nav">{navigationGroups.map((group) => <div className="nav-group" key={group.label}><p className="nav-group-label">{group.label}</p>{group.items.map((item) => { const isActive = item.href === activeHref; return <a aria-current={isActive ? "page" : undefined} className={`nav-item ${isActive ? "nav-item-active" : ""}`} href={item.href} key={item.href}><Icon name={item.icon} size={17} /><span>{item.label}</span></a>; })}</div>)}</nav>
      <div className="sidebar-footer"><a className="nav-item" href="/documentation"><Icon name="help" size={17} /><span>Documentation</span></a><a className="nav-item" href="/settings"><Icon name="settings" size={17} /><span>Settings</span></a><div className="sidebar-server"><span className="server-pulse" aria-hidden="true" /><div><strong>web-01</strong><span>Production</span></div><Icon name="chevron-down" size={15} /></div></div>
    </aside>
    <div className="admin-main">
      <header className="topbar"><div className="mobile-brand"><span className="brand-mark">H</span><strong>HarborPanel</strong></div><div className="topbar-context"><span className="context-dot" /> Production / web-01 <Icon name="chevron-down" size={14} /></div><div className="topbar-actions"><label className="search-field"><Icon name="search" size={17} /><span className="sr-only">Search administration</span><input aria-label="Search administration" placeholder="Search" type="search" /></label><a className="icon-button" href="/activity" aria-label="View activity log"><Icon name="activity" size={18} /></a><a className="user-menu" href="/settings"><span className="avatar">AD</span><span className="user-name">Admin</span><Icon name="chevron-down" size={14} /></a></div></header>
      <main>{children}</main>
    </div>
  </div>;
}
