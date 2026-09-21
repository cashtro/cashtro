#!/usr/bin/env node
import { execFileSync } from "node:child_process";

const PATTERNS = [
  /AKIA[0-9A-Z]{16}/,
  /ghp_[A-Za-z0-9]{20,}/,
  /ghs_[A-Za-z0-9]{20,}/,
  /sk-[A-Za-z0-9]{20,}/,
  /-----BEGIN (RSA |OPENSSH )?PRIVATE KEY-----/,
  /xox[baprs]-[A-Za-z0-9-]{10,}/,
];

function files(): string[] {
  const out = execFileSync("git", ["ls-files"], { encoding: "utf8" });
  return out.split("\n").filter((f) => f && !f.startsWith("inventory/repos.json"));
}

export function scanText(name: string, text: string): string[] {
  const hits: string[] = [];
  for (const re of PATTERNS) {
    if (re.test(text)) hits.push(`${name}: matched ${re}`);
  }
  return hits;
}

function main() {
  const hits: string[] = [];
  for (const file of files()) {
    if (file.includes("node_modules") || file.endsWith(".lock") || file.includes("secret-scan")) continue;
    let text = "";
    try {
      text = execFileSync("git", ["show", `:${file}`], { encoding: "utf8" });
    } catch {
      continue;
    }
    hits.push(...scanText(file, text));
  }
  if (hits.length) {
    process.stderr.write(hits.join("\n") + "\n");
    process.exit(2);
  }
  process.stdout.write("secret-scan ok\n");
}

if (import.meta.url === `file://${process.argv[1]}` || process.argv[1]?.endsWith("cli.ts")) {
  main();
}
