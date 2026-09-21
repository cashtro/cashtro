import { setTimeout as sleep } from "node:timers/promises";

export type GhRepo = {
  name: string;
  full_name: string;
  private: boolean;
  archived: boolean;
  default_branch: string;
  updated_at: string;
  pushed_at: string;
  language: string | null;
  owner: { login: string; type: string };
  html_url: string;
  description: string | null;
  stargazers_count?: number;
};

export class GitHub {
  constructor(
    private token: string,
    private fetchImpl: typeof fetch = fetch,
  ) {}

  async request<T>(path: string, init: RequestInit = {}): Promise<{ status: number; data: T; headers: Headers }> {
    const url = path.startsWith("http") ? path : `https://api.github.com${path}`;
    for (let attempt = 0; attempt < 6; attempt++) {
      const res = await this.fetchImpl(url, {
        ...init,
        headers: {
          Accept: "application/vnd.github+json",
          "X-GitHub-Api-Version": "2022-11-28",
          "User-Agent": "cashtro-recon",
          ...(this.token ? { Authorization: `Bearer ${this.token}` } : {}),
          ...(init.headers || {}),
        },
      });
      if (res.status === 403 && res.headers.get("x-ratelimit-remaining") === "0") {
        const reset = Number(res.headers.get("x-ratelimit-reset") || "0") * 1000;
        const wait = Math.max(reset - Date.now(), 1000 * 2 ** attempt);
        await sleep(Math.min(wait, 60_000));
        continue;
      }
      if (res.status >= 500) {
        await sleep(250 * 2 ** attempt);
        continue;
      }
      const text = await res.text();
      const data = text ? (JSON.parse(text) as T) : (undefined as T);
      return { status: res.status, data, headers: res.headers };
    }
    throw new Error(`github retry exhausted: ${path}`);
  }

  async listOwnerRepos(owner: string, kind: "user" | "org"): Promise<{ access: "ok" | "denied" | "limited"; repos: GhRepo[] }> {
    const path =
      kind === "org"
        ? `/orgs/${owner}/repos?per_page=100&type=all`
        : `/users/${owner}/repos?per_page=100&type=all`;
    const { status, data } = await this.request<GhRepo[] | { message?: string }>(path);
    if (status === 404 || status === 403) return { access: "denied", repos: [] };
    if (status !== 200 || !Array.isArray(data)) return { access: "denied", repos: [] };
    if (kind === "org" && data.length === 0) {
      const org = await this.request(`/orgs/${owner}`);
      if (org.status === 200) return { access: "limited", repos: [] };
    }
    return { access: "ok", repos: data };
  }

  async getRepo(fullName: string): Promise<GhRepo | null> {
    const { status, data } = await this.request<GhRepo>(`/repos/${fullName}`);
    if (status === 404 || status === 403) return null;
    if (status !== 200) return null;
    return data;
  }

  async getFile(fullName: string, path: string): Promise<string | null> {
    const { status, data } = await this.request<{ content?: string; encoding?: string }>(
      `/repos/${fullName}/contents/${path}`,
    );
    if (status === 404 || status === 403) return null;
    if (status !== 200 || !data.content) return null;
    if (data.encoding === "base64") return Buffer.from(data.content.replace(/\n/g, ""), "base64").toString("utf8");
    return data.content;
  }

  async listDir(fullName: string, path: string): Promise<string[]> {
    const { status, data } = await this.request<Array<{ name: string; type: string }>>(
      `/repos/${fullName}/contents/${path}`,
    );
    if (status !== 200 || !Array.isArray(data)) return [];
    return data.map((e) => e.name);
  }

  async latestWorkflow(fullName: string): Promise<string | null> {
    const { status, data } = await this.request<{ workflow_runs?: Array<{ conclusion: string | null }> }>(
      `/repos/${fullName}/actions/runs?per_page=1`,
    );
    if (status !== 200 || !data.workflow_runs?.length) return null;
    return data.workflow_runs[0].conclusion ?? "unknown";
  }

  async openCounts(fullName: string): Promise<{ prs: number; issues: number }> {
    const prs = await this.request<unknown[]>(`/repos/${fullName}/pulls?state=open&per_page=1`);
    const issues = await this.request<unknown[]>(`/repos/${fullName}/issues?state=open&per_page=1`);
    const count = (headers: Headers) => {
      const link = headers.get("link") || "";
      const m = link.match(/page=(\d+)>; rel="last"/);
      if (m) return Number(m[1]);
      return 0;
    };
    const prN = prs.status === 200 ? (count(prs.headers) || (Array.isArray(prs.data) ? prs.data.length : 0)) : 0;
    const issueN = issues.status === 200 ? (count(issues.headers) || (Array.isArray(issues.data) ? issues.data.length : 0)) : 0;
    return { prs: prN, issues: Math.max(0, issueN - prN) };
  }
}

export async function resolveToken(): Promise<string> {
  if (process.env.GITHUB_TOKEN) return process.env.GITHUB_TOKEN;
  try {
    const { execFile } = await import("node:child_process");
    const { promisify } = await import("node:util");
    const exec = promisify(execFile);
    const { stdout } = await exec("gh", ["auth", "token"]);
    return stdout.trim();
  } catch {
    return "";
  }
}
