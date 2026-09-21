// cursor-canvas-title: Cashtro OS — thinkers vs products
import { Callout, Grid, H1, H2, Pill, Stack, Stat, Table, Text } from "cursor/canvas";

const thinkers = [
  ["Lovelace", "init", "kernel", "live"],
  ["Turing", "router", "model", "resident"],
  ["Deming", "delivery", "ship", "live"],
  ["Hypatia", "research", "library", "live"],
  ["Copernicus", "explorer", "search", "live"],
  ["Galileo", "operator", "computer-use", "resident"],
  ["Kant", "reviewer", "qa", "resident"],
  ["Leonardo", "architect", "design", "resident"],
  ["Hopper", "deploy", "release", "resident"],
  ["Machiavelli", "security", "guard", "resident"],
  ["Locke", "memory", "recall", "live"],
  ["Voltaire", "comms", "signal", "live"],
  ["Confucius", "planner", "backlog", "live"],
  ["Socrates", "investigator", "incident", "resident"],
] as const;

const products = [
  ["Casa Crypto", "concept", "finance"],
  ["Prolifik", "concept", "brand"],
  ["ScanApp", "concept", "ops"],
  ["Proximity", "production", "agency"],
  ["BTK Avocats", "production", "law"],
  ["MD Clinic", "production", "health"],
  ["Solution Hypothèque QC", "production", "finance"],
  ["Éduconnexion", "production", "education"],
  ["Cashtro delivery catalog", "idea", "internal"],
] as const;

export default function AgenticOsPlan() {
  return (
    <Stack gap={20}>
      <H1>Cashtro OS naming</H1>
      <Callout tone="info">
        Products stay products. Processes wear thinkers. IDs and capability verbs
        do not change — Lovelace is still process init, Casa Crypto is still a ship.
      </Callout>
      <Grid columns={3} gap={12}>
        <Stat value="14" label="Thinkers on the process table" />
        <Stat value="9" label="Products on the delivery line" />
        <Stat value="0" label="Products worn as process titles" tone="success" />
      </Grid>
      <H2>Process table — thinkers</H2>
      <Text tone="secondary" size="small">
        Source: Cashtro kernel 0.3.0 · roster 2026-09-21
      </Text>
      <Table
        headers={["Thinker", "Process ID", "Role", "Mode"]}
        striped
        stickyHeader
        rowTone={thinkers.map((row) => (row[3] === "live" ? "success" : "warning"))}
        rows={thinkers.map((row) => [
          row[0],
          row[1],
          row[2],
          <Pill key={row[1]} tone={row[3] === "live" ? "success" : "warning"} size="sm">
            {row[3]}
          </Pill>,
        ])}
      />
      <H2>Delivery line — products keep their names</H2>
      <Text tone="secondary" size="small">
        Casa Crypto and Prolifik stay ships. They are never process titles.
      </Text>
      <Table
        headers={["Product", "Stage", "Sector"]}
        striped
        rowTone={products.map((row) =>
          row[1] === "production" ? "success" : row[1] === "concept" ? "warning" : "info",
        )}
        rows={products.map((row) => [row[0], row[1], row[2]])}
      />
    </Stack>
  );
}
