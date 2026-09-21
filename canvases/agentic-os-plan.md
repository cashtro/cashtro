# Agentic OS plan — builders vs products

Castro rule: **products stay products**. Processes wear names like
**Builder**, **Creator**, **Lautaro** — then the rest of the crew.

See `~/.cursor/admiral-fleet/ROSTER.md` for the live roster.

```mermaid
flowchart TB
  Castro[Castro]
  Lautaro[Lautaro · init]
  Builder[Builder · operator]
  Creator[Creator · architect]
  Maker[Maker · delivery]
  subgraph products [Products keep their names]
    CasaCrypto[Casa Crypto]
    Prolifik[Prolifik]
    ScanApp[ScanApp]
    Proximity[Proximity]
  end
  Castro --> Lautaro
  Lautaro --> Builder
  Lautaro --> Creator
  Lautaro --> Maker
  Maker --> CasaCrypto
  Maker --> Prolifik
  Maker --> ScanApp
  Maker --> Proximity
```
