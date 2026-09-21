// cursor-canvas-title: Cashtro OS — builders vs products
import { Callout, Grid, H1, H2, Pill, Stack, Stat, Table, Text } from "cursor/canvas";

const builders = [
  ["Lautaro", "init", "kernel", "live"],
  ["Hedy", "router", "model", "resident"],
  ["Antoni", "delivery", "ship", "live"],
  ["Gabriela", "research", "library", "live"],
  ["Marco", "explorer", "search", "live"],
  ["Diego", "operator", "computer-use", "resident"],
  ["Juana", "reviewer", "qa", "resident"],
  ["Frida", "architect", "design", "resident"],
  ["Oscar", "deploy", "release", "resident"],
  ["Túpac", "security", "guard", "resident"],
  ["Luis", "memory", "recall", "live"],
  ["Violeta", "comms", "signal", "live"],
  ["Pablo", "planner", "backlog", "live"],
  ["Caupolicán", "investigator", "incident", "resident"],
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
        Products stay products. Processes wear builders and creators — Lautaro
        boots the house. IDs and capability verbs do not change. Casa Crypto
        stays a ship.
      </Callout>
      <Grid columns={3} gap={12}>
        <Stat value="14" label="Builders on the process table" />
        <Stat value="9" label="Products on the delivery line" />
        <Stat value="Lautaro" label="Kernel seat · init" />
      </Grid>
      <H2>Process table — builders</H2>
      <Text tone="secondary" size="small">
        Source: Cashtro kernel 0.3.1 · roster 2026-09-21
      </Text>
      <Table
        headers={["Builder", "Process ID", "Role", "Mode"]}
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
