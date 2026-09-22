# Inventory graph

Cashtro (control plane) and Evolu-Jeunes (delivery org) mapped together.
Live client ships are catalog-named only. Do not touch production.
Canonical desk map: `GET /map` · JSON `GET /api/map` · `docs/FLEET_MAP.md`.

```mermaid
flowchart TB
  castro[Castro]
  subgraph cashtro_home[cashtro GitHub user]
    repo_cashtro_cashtro[cashtro/cashtro]
  end
  subgraph evolu_home[Evolu-Jeunes GitHub org]
    evolu_member[github.com/evoluJeunes]
    evolu_site[evolujeunes.ca / LIVE do not touch]
    evolu_dark[~48 private repos / access limited]
    subgraph live[LIVE CLIENT — catalog named only]
      cat_BTK_Avocats[BTK Avocats / production]
      cat_MD_Clinic[MD Clinic / production]
      cat_Solution_Hypoth_que_QC[Solution Hypothèque QC / production]
      cat__duconnexion[Éduconnexion / production]
      cat_Proximity[Proximity / production]
    end
    subgraph later[Safe to onboard later]
      cat_ScanApp[ScanApp / concept]
      cat_Cashtro_delivery_catalog[Cashtro delivery catalog / idea]
    end
  end
  castro -->|owns| cashtro_home
  castro -->|owns| evolu_home
  repo_cashtro_cashtro -->|control plane / recon| evolu_home
  repo_cashtro_cashtro -->|catalog seed| live
  repo_cashtro_cashtro -->|catalog seed| later
  evolu_home -.-> evolu_dark
```
