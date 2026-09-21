# Agentic OS plan — builders vs products

Castro rule: **products stay products**. Processes wear builder / creator
names like **Lautaro**.

See `~/.cursor/admiral-fleet/ROSTER.md` for the live roster.

```mermaid
flowchart TB
  Castro[Castro]
  Lautaro[Lautaro · init]
  Antoni[Antoni · delivery]
  subgraph products [Products keep their names]
    CasaCrypto[Casa Crypto]
    Prolifik[Prolifik]
    ScanApp[ScanApp]
    Proximity[Proximity]
  end
  Castro --> Lautaro
  Lautaro --> Antoni
  Antoni --> CasaCrypto
  Antoni --> Prolifik
  Antoni --> ScanApp
  Antoni --> Proximity
```
