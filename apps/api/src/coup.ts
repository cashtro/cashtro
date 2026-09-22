export const COUP_FIELDS = ["question", "position", "shortest", "opponent"] as const;
export type CoupBeat = (typeof COUP_FIELDS)[number];
export type Coup = Record<CoupBeat, string>;

export const COUP_LABELS: Record<CoupBeat, string> = {
  question: "Questions spécifiques",
  position: "Position",
  shortest: "Coup le plus court",
  opponent: "Réponse adverse (Steel)",
};

export function parseUrlEncoded(raw: string): Record<string, string> {
  const out: Record<string, string> = {};
  for (const [k, v] of new URLSearchParams(raw)) out[k] = v;
  return out;
}

export function parseCoup(body: unknown): Coup {
  const src = body && typeof body === "object" ? (body as Record<string, unknown>) : {};
  const pick = (k: CoupBeat) => String(src[k] ?? "").trim();
  return {
    question: pick("question"),
    position: pick("position"),
    shortest: pick("shortest"),
    opponent: pick("opponent"),
  };
}

/** Four beats, and the first must actually be questions that fill context. */
export function coupComplete(coup: Coup): boolean {
  if (!COUP_FIELDS.every((k) => coup[k].length > 0)) return false;
  return coup.question.includes("?");
}

export function missingBeats(coup: Coup): CoupBeat[] {
  const missing = COUP_FIELDS.filter((k) => !coup[k]);
  if (!missing.includes("question") && coup.question.length > 0 && !coup.question.includes("?")) {
    return ["question", ...missing];
  }
  return missing;
}
