import { z } from "zod";
import { commonSchemas } from "./common";

export const clientSchema = z.object({
  id: z.optional(commonSchemas.id),
  name: commonSchemas.requiredString,
  email: commonSchemas.email,
  website: z.optional(z.string()),
  avatar: z.optional(z.string()),
  contactPerson: z.optional(z.string()),
  address: z.optional(z.string()),
  currency: commonSchemas.requiredString,
  status: z.enum(["active", "inactive"]),
  notes: z.optional(z.string()),
});
