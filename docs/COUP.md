# Coup — questions d’abord, puis on joue

Avant chaque action : des **questions spécifiques** pour combler le
contexte, la position, le coup le plus court, **Steel**. Ensuite seulement on joue.

Le pane Fastify ne dispatch plus d’un clic. L’API JSON `/tasks/:id/dispatch` reste pour les tests et le CI.

```mermaid
flowchart TB
  Q[1 Questions spécifiques]
  P[2 Position]
  S[3 Coup le plus court]
  R[4 Steel — réponse adverse]
  J[Jouer]
  Q --> P --> S --> R --> J
  J --> CP[orchestrate]
  CP --> FAB[58 specialists]
  FAB --> V[14 kernel :8080]
  FAB --> ST[Steel local]
```

## Pane

1. Queue → lien **Coup** → `GET /ui/coup/:id`
2. Le formulaire propose les questions du département matché (`state/corporation.json`)
3. Le premier champ doit contenir un `?`
4. **Jouer** poste `POST /ui/act/dispatch/:id`
5. Incomplet → `?missing=1`
6. Complet → événement `coup.before`, puis orchestrate
7. Si le siège est Steel → adapter local `$0`, refuse live clients

## Ce que ça n’est pas

- Pas une partie d’échecs
- Pas un second produit
- Pas un gate sur l’API JSON
- Kill switch Pause reste un clic
