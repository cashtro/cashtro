import assert from "node:assert/strict";
import { test } from "node:test";
import { coupComplete, missingBeats, parseCoup, parseUrlEncoded } from "./coup.js";

test("urlencoded parser keeps the four beats", () => {
  const body = parseUrlEncoded(
    "question=quoi+%3F&position=ici&shortest=un+lien&opponent=un+clic+de+moins",
  );
  const coup = parseCoup(body);
  assert.equal(coup.question, "quoi ?");
  assert.equal(coup.position, "ici");
  assert.equal(coup.shortest, "un lien");
  assert.equal(coup.opponent, "un clic de moins");
  assert.equal(coupComplete(coup), true);
  assert.deepEqual(missingBeats(coup), []);
});

test("empty or partial coup is not playable", () => {
  assert.equal(coupComplete(parseCoup({})), false);
  assert.deepEqual(missingBeats(parseCoup({ question: "x" })), ["question", "position", "shortest", "opponent"]);
  assert.equal(coupComplete(parseCoup({ question: "  ", position: "p", shortest: "s", opponent: "o" })), false);
  assert.equal(coupComplete(parseCoup({ question: "pas une question", position: "p", shortest: "s", opponent: "o" })), false);
});
