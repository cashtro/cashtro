import assert from "node:assert/strict";
import { test } from "node:test";
import { loadCorporation, matchDepartment, probesFor } from "./corporation.js";

test("corporation has 14 kernel seats plus Steel", () => {
  const corp = loadCorporation();
  assert.equal(corp.kernelSeats, 14);
  assert.equal(corp.departments, 15);
  assert.equal(corp.seats.length, 15);
  const steel = corp.seats.find((s) => s.id === "contrarian");
  assert.ok(steel);
  assert.equal(steel?.crew, "Steel");
  assert.equal(steel?.layer, "control-plane");
});

test("probes are specific to the matched department", () => {
  assert.equal(matchDepartment("improve delivery ship")?.id, "delivery");
  const probes = probesFor("steel contradict the plan");
  assert.ok(probes.some((q) => q.includes("?")));
  assert.ok(probes.length >= 2);
});
