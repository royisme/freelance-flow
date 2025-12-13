import { z } from "zod";
import { commonSchemas } from "./common";

export const timeEntrySchema = z.object({
  id: z.optional(commonSchemas.id),
  projectId: z.number().int().positive(),
  date: commonSchemas.dateString,
  startTime: z.optional(z.string()),
  endTime: z.optional(z.string()),
  durationSeconds: z.number().positive(), // Must be > 0
  description: commonSchemas.requiredString,
  invoiced: z.boolean(),
});
