# ACCESS_REQUIRED

Written 2026-09-21. Do not guess missing repos.

## Token in this Cloud Agent

`gh auth status` is logged in as GitHub account **`cursor`** (Cursor
integration token, `ghs_…`). Not Castro's user. Not an org install.

| Check | Result |
| --- | --- |
| `gh api user` | 403 Resource not accessible by integration |
| `gh repo list cashtro --limit 200` | 1 public repo: `cashtro/cashtro` |
| `gh api users/cashtro/repos` | same: `cashtro` public |
| `gh api orgs/Evolu-Jeunes` | 200 — org exists |
| `gh api orgs/Evolu-Jeunes/repos?type=all` | `[]` — zero visible repos |
| `gh search repos org:Evolu-Jeunes` | empty |
| `gh repo list Evolu-Jeunes` | empty |

README says Evolu-Jeunes holds ~48 private client repos. This token
cannot see them. Recon records the org as `access: denied` and continues.

## Orgs / users affected

- **Evolu-Jeunes** — all repositories (expected private client fleet)
- **cashtro** — any private repos besides public `cashtro/cashtro`
- Any other org Castro owns that was not listed in the brief

## What I need — Option A first (fine-grained PAT)

GitHub → Settings → Developer settings → Personal access tokens →
Fine-grained → new token.

1. Resource owner = **Evolu-Jeunes** (repeat for **cashtro** if needed)
2. Repository access = **All repositories**
3. Permissions:
   - Contents: **read**
   - Metadata: **read**
   - Pull requests: **read and write** (needed later for onboard PRs)
   - Issues: **read and write**
   - Actions: **read**

Then in this environment:

```bash
gh auth login --with-token
pnpm recon
```

Option B (GitHub App, long-term) and Option C (machine account) are
documented in the mission brief. Start with A.

I am not requesting a token with admin, org:write, or delete scopes.
Do not paste the token into chat — put it in the Cloud Agent secret
store or `gh auth login`.
