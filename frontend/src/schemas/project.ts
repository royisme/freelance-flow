import { z } from "zod";
import { commonSchemas } from "./common";

export const projectSchema = z.object({
  id: z.optional(commonSchemas.id),
  clientId: z.number().int().positive(),
  name: commonSchemas.requiredString,
  description: z.optional(z.string()),
  hourlyRate: z.number().min(0),
  currency: commonSchemas.requiredString,
  status: z.enum(["active", "archived", "completed"]),
  deadline: z.optional(z.string()), // Or date validator
  tags: z.array(z.string()),
});
