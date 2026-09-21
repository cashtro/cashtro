// cursor-canvas-title: Cashtro OS — builders vs products
import { Callout, Grid, H1, H2, Pill, Stack, Stat, Table, Text } from "cursor/canvas";

const builders = [
  ["Lautaro", "init", "kernel", "live"],
  ["Builder", "operator", "computer-use", "resident"],
  ["Creator", "architect", "design", "resident"],
  ["Maker", "delivery", "ship", "live"],
  ["Forger", "deploy", "release", "resident"],
  ["Scribe", "research", "library", "live"],
  ["Pathfinder", "explorer", "search", "live"],
  ["Witness", "reviewer", "qa", "resident"],
  ["Sentinel", "security", "guard", "resident"],
  ["Keeper", "memory", "recall", "live"],
  ["Herald", "comms", "signal", "live"],
  ["Steward", "planner", "backlog", "live"],
  ["Seeker", "investigator", "incident", "resident"],
  ["Spark", "router", "model", "resident"],
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
        Products stay products. Processes wear names like Builder, Creator,
        Lautaro — then the rest of the crew. IDs and capability verbs do not
        change. Casa Crypto stays a ship.
      </Callout>
      <Grid columns={3} gap={12}>
        <Stat value="Lautaro" label="Kernel seat · init" />
        <Stat value="Builder" label="Hands · operator" />
        <Stat value="Creator" label="Form · architect" />
      </Grid>
      <H2>Process table — builders</H2>
      <Text tone="secondary" size="small">
        Source: Cashtro kernel 0.3.2 · roster 2026-09-21
      </Text>
      <Table
        headers={["Name", "Process ID", "Role", "Mode"]}
        striped
        stickyHeader
        rowTone={builders.map((row) => (row[3] === "live" ? "success" : "warning"))}
        rows={builders.map((row) => [
          row[0],
          row[1],
          row[2],
          <Pill key={row[1]} size="sm">
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
