import "./globals.css";

export const metadata = {
  title: "HarborPanel Admin",
  description: "Hosting administration control plane",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
