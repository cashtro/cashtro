import { z } from "zod";

export const ProjectKind = z.enum(["control-plane", "product", "client", "internal"]);
export const ProjectStatus = z.enum(["active", "concept", "idea", "archived", "denied"]);
export const TaskStatus = z.enum(["todo", "doing", "done", "blocked", "cancelled"]);
export const RunStatus = z.enum(["queued", "running", "succeeded", "failed", "cancelled", "blocked-budget"]);
export const AgentRuntime = z.enum(["cursor", "claude-code", "bedrock", "n8n", "ghl", "http"]);

export const Project = z.object({
  id: z.string(),
  slug: z.string(),
  repoUrl: z.string().nullable(),
  org: z.string(),
  kind: ProjectKind,
  status: ProjectStatus,
  deployTarget: z.string().nullable(),
  healthUrl: z.string().nullable(),
  tags: z.array(z.string()),
});
export type Project = z.infer<typeof Project>;

export const Agent = z.object({
  id: z.string(),
  projectId: z.string().nullable(),
  name: z.string(),
  role: z.string(),
  runtime: AgentRuntime,
  modelRoute: z.string(),
  systemPromptRef: z.string().nullable(),
  maxCostPerRun: z.number(),
  enabled: z.boolean(),
});
export type Agent = z.infer<typeof Agent>;

export const Task = z.object({
  id: z.string(),
  projectId: z.string().nullable(),
  title: z.string(),
  body: z.string(),
  source: z.string(),
  priority: z.string(),
  status: TaskStatus,
  dependsOn: z.array(z.string()),
  assigneeAgentId: z.string().nullable(),
  idempotencyKey: z.string().nullable(),
});
export type Task = z.infer<typeof Task>;

export const Run = z.object({
  id: z.string(),
  taskId: z.string(),
  agentId: z.string(),
  startedAt: z.string(),
  endedAt: z.string().nullable(),
  status: RunStatus,
  inputTokens: z.number(),
  outputTokens: z.number(),
  costUsd: z.number(),
  exitReason: z.string().nullable(),
});
export type Run = z.infer<typeof Run>;

export const CreateTask = z.object({
  title: z.string().min(1),
  body: z.string().default(""),
  source: z.string().default("api"),
  priority: z.string().default("p2"),
  projectId: z.string().optional(),
  assigneeAgentId: z.string().optional(),
  dependsOn: z.array(z.string()).default([]),
  idempotencyKey: z.string().min(8),
});
export type CreateTask = z.infer<typeof CreateTask>;
