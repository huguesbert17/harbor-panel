import "./globals.css";

export const metadata = { title: "Dashboard · HarborPanel", description: "HarborPanel customer hosting dashboard" };

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>;
}
