# Inventory graph

```mermaid
flowchart LR
  manager[cashtro/cashtro]
  cashtro[cashtro/cashtro]
  manager --> cashtro
  Evolu_Jeunes_dark[Evolu-Jeunes / denied]
  manager -.-> Evolu_Jeunes_dark
  subgraph catalog[Kernel catalog — not scanned]
    cat_btk_avocats[BTK Avocats / production]
    cat_md_clinic[MD Clinic / production]
    cat_solution_hypotheque_qc[Solution Hypothèque QC / production]
    cat_educonnexion[Éduconnexion / production]
    cat_proximity[Proximity / production]
    cat_scanapp[ScanApp / concept]
    cat_cashtro_catalog[Cashtro delivery catalog / idea]
  end
  manager -.-> catalog
```
