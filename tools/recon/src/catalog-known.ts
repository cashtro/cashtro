import shipsFile from "../../../state/catalog-ships.json" with { type: "json" };

export type CatalogShip = {
  slug: string;
  name: string;
  org: string;
  kind: string;
  status: string;
  sector: string;
  stack: string[];
  liveClient: boolean;
  notes: string;
};

export const CATALOG_SHIPS: CatalogShip[] = shipsFile.ships;
