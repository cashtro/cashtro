import { PrismaClient } from "@prisma/client";
import { buildApp } from "./app.js";

const prisma = new PrismaClient();
const port = Number(process.env.PORT || 8787);

const app = await buildApp({ prisma });
await app.listen({ port, host: "127.0.0.1" });
console.log(`control plane listening on http://127.0.0.1:${port}  docs /docs`);
