import type { ReactNode } from "react";

import { Icon, type IconName } from "./icons";

export type StatusTone = "success" | "warning" | "danger" | "neutral" | "info";

export function Panel({ action, children, className = "", eyebrow, title }: { action?: ReactNode; children: ReactNode; className?: string; eyebrow?: string; title: string }) {
  return <section className={`panel ${className}`}><div className="panel-header"><div>{eyebrow ? <p className="panel-eyebrow">{eyebrow}</p> : null}<h2 className="panel-title">{title}</h2></div>{action ? <div className="panel-action">{action}</div> : null}</div>{children}</section>;
}

export function StatusBadge({ label, tone = "neutral" }: { label: string; tone?: StatusTone }) {
  const icon: Record<StatusTone, IconName> = { success: "circle-check", warning: "clock", danger: "x-circle", neutral: "circle-info", info: "circle-info" };
  return <span className={`status-badge status-${tone}`}><Icon name={icon[tone]} size={14} />{label}</span>;
}

export function MetricCard({ detail, icon, label, tone = "blue", value }: { detail: string; icon: IconName; label: string; tone?: "blue" | "green" | "orange" | "purple"; value: string }) {
  return <article className="metric-card"><div className={`metric-icon metric-icon-${tone}`}><Icon name={icon} size={19} /></div><p className="metric-label">{label}</p><p className="metric-value">{value}</p><p className="metric-detail">{detail}</p></article>;
}

export function ProgressBar({ label, value, detail, tone = "blue" }: { detail: string; label: string; tone?: "blue" | "green" | "orange"; value: number }) {
  return <div className="progress-row"><div className="progress-meta"><span>{label}</span><strong>{detail}</strong></div><div aria-label={`${label}: ${value}% used`} aria-valuemax={100} aria-valuemin={0} aria-valuenow={value} className="progress-track" role="progressbar"><span className={`progress-fill progress-fill-${tone}`} style={{ width: `${value}%` }} /></div></div>;
}

export function LinkButton({ children, href = "#" }: { children: ReactNode; href?: string }) {
  return <a className="link-button" href={href}>{children}<Icon name="arrow-up-right" size={15} /></a>;
}
