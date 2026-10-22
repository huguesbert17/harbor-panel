export const navigation = {
  admin: [
    { label: "Home", href: "/" },
    { label: "Accounts", href: "/accounts" },
    { label: "Packages", href: "/packages" },
    { label: "Domains", href: "/domains" },
    { label: "DNS Zones", href: "/dns" },
    { label: "Servers", href: "/servers" },
    { label: "Backups", href: "/backups" },
    { label: "Security", href: "/security" },
  ],
  customer: [
    { label: "Dashboard", href: "/" },
    { label: "Domains", href: "/domains" },
    { label: "File Manager", href: "/files" },
    { label: "Databases", href: "/databases" },
    { label: "Email", href: "/email" },
    { label: "DNS Zone Editor", href: "/dns" },
    { label: "SSL/TLS", href: "/ssl" },
    { label: "Backups", href: "/backups" },
  ],
} as const;
