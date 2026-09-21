# ACCESS_REQUIRED

Checked again 2026-09-21T22:34Z after Castro connected Cursor to
cashtro + Evolu-Jeunes. I did check Evolu-Jeunes. Still dark.

Connecting Cursor on your machine is **not** this Cloud Agent's token.

## What I just ran

```
gh auth status
→ logged in as GitHub account **cursor** (integration token ghs_…)
   not cashtro, not an Evolu-Jeunes install

gh api orgs/Evolu-Jeunes
→ 200  login=Evolu-Jeunes  id=201155686  name="Evolu Jeunes"
   public_repos=0

gh api graphql organization(login:"Evolu-Jeunes").repositories
→ totalCount: 0  nodes: []

gh api orgs/Evolu-Jeunes/repos?type=all
→ []

gh api /installation/repositories
→ total_count: 1
   repository_selection: **selected**
   only repo: **cashtro/cashtro**

https://github.com/orgs/Evolu-Jeunes/repositories
→ "This organization has no public repositories." / 0 repositories
```

So: the org is real. Every repo in it is private (or the org is empty).
This agent’s Cursor GitHub App is installed on **one selected repo**
(`cashtro/cashtro`). It was never granted Evolu-Jeunes.

## What I need — Cursor GitHub App on the org (do this, not a PAT)

You already connected Cursor. Now give **this Cloud Agent’s GitHub App**
the org.

1. Open https://github.com/organizations/Evolu-Jeunes/settings/installations
2. Find **Cursor** (or **Cursor Cloud Agents**). If it is missing: install it.
3. Repository access = **All repositories** (not “Only select”).
4. Confirm cashtro user install is not the only one:
   https://github.com/settings/installations → Cursor → must not be
   limited to `cashtro/cashtro` if you want private user repos too.

Then tell me “Cursor App is on Evolu-Jeunes / all repos”. I re-run
`pnpm recon`. I will not invent repo names.

Do not paste a token in chat.

## Fallback — Option A PAT (only if the App cannot be installed on the org)

Fine-grained PAT, resource owner **Evolu-Jeunes**, all repos,
Contents/Metadata/Actions read, Issues/PRs read-write. Put it in
`GH_TOKEN` on this Cloud Agent. Never in chat.

## Orgs affected

- **Evolu-Jeunes** — 0 public, all private invisible from here
- **cashtro** — only public `cashtro/cashtro` is on the installation
