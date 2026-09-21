/** Ships named in the Go kernel catalog. Not GitHub-scanned. */
export const CATALOG_SHIPS = [
  { name: "BTK Avocats", stage: "production", stack: "WordPress", sector: "law", liveClient: true },
  { name: "MD Clinic", stage: "production", stack: "WordPress, ACF", sector: "health", liveClient: true },
  { name: "Solution Hypothèque QC", stage: "production", stack: "WordPress", sector: "finance", liveClient: true },
  { name: "Éduconnexion", stage: "production", stack: "WordPress", sector: "education", liveClient: true },
  { name: "Proximity", stage: "production", stack: "Next.js, Azure", sector: "agency", liveClient: true },
  { name: "ScanApp", stage: "concept", stack: "TypeScript, Python", sector: "ops", liveClient: false },
  { name: "Cashtro delivery catalog", stage: "idea", stack: "Go", sector: "internal", liveClient: false },
] as const;
