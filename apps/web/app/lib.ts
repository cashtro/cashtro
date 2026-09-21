const API = process.env.CONTROL_PLANE_URL || "http://127.0.0.1:8787";
const KEY = process.env.CONTROL_PLANE_API_KEY || "change-me";

export async function api<T>(path: string): Promise<T | null> {
  try {
    const res = await fetch(`${API}${path}`, {
      headers: { authorization: `Bearer ${KEY}` },
      cache: "no-store",
    });
    if (!res.ok) return null;
    return (await res.json()) as T;
  } catch {
    return null;
  }
}

export async function uiHtml(screen: string): Promise<string> {
  try {
    const res = await fetch(`${API}/ui/${screen}`, { cache: "no-store" });
    return await res.text();
  } catch {
    return "<p>Control plane API is not up. Start with <code>pnpm api</code>.</p>";
  }
}
