import { readFileSync } from "node:fs";
import path from "node:path";

export type Dept = {
  id: string;
  crew: string;
  role: string;
  layer: string;
  runtime?: string;
  mission: string;
  gap: string;
  ask: string[];
  selfImprove: string;
  maxPush: string;
};

export type Corporation = {
  schemaVersion: string;
  principle: string;
  kernelSeats: number;
  departments: number;
  depthLimit: number;
  liveClientsDoNotTouch: string[];
  seats: Dept[];
};

const GENERIC_ASK = [
  "Qu’est-ce qu’on ne sait pas encore, assez pour ne pas jouer ?",
  "Quel département est vraiment propriétaire ?",
  "Quelle option plus courte survit à une contradiction ?",
];

function readCorp(start = process.cwd()): Corporation | null {
  const files = [
    path.resolve(start, "state/corporation.json"),
    path.resolve(start, "../../state/corporation.json"),
    path.resolve(start, "../../../state/corporation.json"),
  ];
  for (const file of files) {
    try {
      return JSON.parse(readFileSync(file, "utf8")) as Corporation;
    } catch {
      // next
    }
  }
  return null;
}

export function loadCorporation(start = process.cwd()): Corporation {
  return (
    readCorp(start) ?? {
      schemaVersion: "1",
      principle: "Understand before positioning.",
      kernelSeats: 14,
      departments: 15,
      depthLimit: 3,
      liveClientsDoNotTouch: [],
      seats: [],
    }
  );
}

export function matchDepartment(blob: string, seats: Dept[] = loadCorporation().seats): Dept | undefined {
  const t = blob.toLowerCase();
  let best: Dept | undefined;
  let score = 0;
  for (const d of seats) {
    const keys = [d.id, d.crew, d.role, d.layer].filter(Boolean).map((k) => k.toLowerCase());
    for (const key of keys) {
      if (key.length >= 4 && t.includes(key) && key.length >= score) {
        best = d;
        score = key.length;
      }
    }
  }
  return best;
}

export function probesFor(blob: string, seats?: Dept[]): string[] {
  const dept = matchDepartment(blob, seats ?? loadCorporation().seats);
  const ask = dept?.ask?.length ? dept.ask : GENERIC_ASK;
  return ask.slice(0, 3);
}
