import type { ReactNode, SVGProps } from "react";

export type IconName =
  | "activity" | "arrow-up-right" | "backup" | "chevron-down" | "circle-check" | "circle-info" | "clock" | "database" | "dns" | "globe" | "grid" | "help" | "home" | "lock" | "menu" | "package" | "search" | "server" | "settings" | "shield" | "sliders" | "users" | "x-circle";

type IconProps = SVGProps<SVGSVGElement> & { name: IconName; size?: number };

const paths: Record<IconName, ReactNode> = {
  activity: <><path d="M3 12h4l2.1-6 4.2 12L15.5 12H21" /><path d="M5 4v16" opacity=".4" /><path d="M19 4v16" opacity=".4" /></>,
  "arrow-up-right": <><path d="M7 17 17 7" /><path d="M8 7h9v9" /></>,
  backup: <><path d="M4 8.5h16v10H4z" /><path d="M7 8.5V5h10v3.5" /><path d="M8 13h8" /></>,
  "chevron-down": <path d="m6 9 6 6 6-6" />,
  "circle-check": <><circle cx="12" cy="12" r="9" /><path d="m8 12 2.5 2.5L16 9" /></>,
  "circle-info": <><circle cx="12" cy="12" r="9" /><path d="M12 11v5" /><path d="M12 8h.01" /></>,
  clock: <><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></>,
  database: <><ellipse cx="12" cy="5" rx="7" ry="3" /><path d="M5 5v7c0 1.7 3.1 3 7 3s7-1.3 7-3V5" /><path d="M5 12v7c0 1.7 3.1 3 7 3s7-1.3 7-3v-7" /></>,
  dns: <><rect x="4" y="4" width="16" height="5" rx="1" /><rect x="4" y="15" width="16" height="5" rx="1" /><path d="M8 6.5h.01M8 17.5h.01" /><path d="M12 9v6" /></>,
  globe: <><circle cx="12" cy="12" r="9" /><path d="M3.5 9h17M3.5 15h17M12 3c2.2 2.5 3.3 5.5 3.3 9S14.2 18.5 12 21c-2.2-2.5-3.3-5.5-3.3-9S9.8 5.5 12 3Z" /></>,
  grid: <><rect x="4" y="4" width="6" height="6" rx="1" /><rect x="14" y="4" width="6" height="6" rx="1" /><rect x="4" y="14" width="6" height="6" rx="1" /><rect x="14" y="14" width="6" height="6" rx="1" /></>,
  help: <><circle cx="12" cy="12" r="9" /><path d="M9.7 9a2.4 2.4 0 1 1 4.1 1.7c-1.1 1.1-1.8 1.3-1.8 2.8" /><path d="M12 16h.01" /></>,
  home: <><path d="m4 10 8-6 8 6" /><path d="M6 9v10h12V9" /><path d="M10 19v-6h4v6" /></>,
  lock: <><rect x="5" y="10" width="14" height="10" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></>,
  menu: <><path d="M4 7h16M4 12h16M4 17h16" /></>,
  package: <><path d="m12 3 8 4.5v9L12 21l-8-4.5v-9L12 3Z" /><path d="m4 7.5 8 4.5 8-4.5M12 12v9" /></>,
  search: <><circle cx="10.8" cy="10.8" r="6.5" /><path d="m16 16 4 4" /></>,
  server: <><rect x="4" y="4" width="16" height="6" rx="1" /><rect x="4" y="14" width="16" height="6" rx="1" /><path d="M8 7h.01M8 17h.01" /><path d="M12 7h5M12 17h5" /></>,
  settings: <><path d="M12 3v2M12 19v2M3 12h2M19 12h2M5.6 5.6 7 7M17 17l1.4 1.4M18.4 5.6 17 7M7 17l-1.4 1.4" /><circle cx="12" cy="12" r="4" /></>,
  shield: <><path d="M12 3 19 6v5c0 4.6-2.7 7.6-7 10-4.3-2.4-7-5.4-7-10V6l7-3Z" /><path d="m9 12 2 2 4-4" /></>,
  sliders: <><path d="M4 6h16M4 12h16M4 18h16" /><circle cx="8" cy="6" r="2" /><circle cx="15" cy="12" r="2" /><circle cx="11" cy="18" r="2" /></>,
  users: <><circle cx="9" cy="8" r="3" /><path d="M3.5 19a5.5 5.5 0 0 1 11 0" /><path d="M16 5.5a3 3 0 0 1 0 5.8M17 14a4.5 4.5 0 0 1 3.5 4" /></>,
  "x-circle": <><circle cx="12" cy="12" r="9" /><path d="m9 9 6 6M15 9l-6 6" /></>,
};

export function Icon({ name, size = 18, ...props }: IconProps) {
  return <svg aria-hidden="true" fill="none" height={size} viewBox="0 0 24 24" width={size} xmlns="http://www.w3.org/2000/svg" {...props}>{paths[name]}</svg>;
}
