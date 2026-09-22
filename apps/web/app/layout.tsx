import type { ReactNode } from "react";

export const metadata = { title: "Cashtro fleet" };

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body style={{ margin: 0, background: "#0b0b0c", color: "#f4f1ea", fontFamily: "ui-sans-serif, system-ui" }}>
        <header style={{ position: "sticky", top: 0, padding: 16, borderBottom: "1px solid #2a2a2c" }}>
          <strong>Cashtro</strong>
          <nav style={{ display: "flex", gap: 12, marginTop: 8 }}>
            <a href="/" style={{ color: "inherit" }}>Fleet</a>
            <a href="/corp" style={{ color: "inherit" }}>Corp</a>
            <a href="/queue" style={{ color: "inherit" }}>Queue</a>
            <a href="/run" style={{ color: "inherit" }}>Run</a>
            <a href="/recon" style={{ color: "inherit" }}>Recon</a>
          </nav>
        </header>
        <main style={{ padding: 16, maxWidth: 720 }}>{children}</main>
      </body>
    </html>
  );
}
