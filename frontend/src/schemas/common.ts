import { z } from "zod";

export const commonSchemas = {
  id: z.number().int().positive(),
  slug: z
    .string()
    .min(1)
    .regex(/^[a-z0-9-]+$/),
  email: z.string().email(),
  password: z.string().min(6), // Minimum password length
  requiredString: z.string().min(1),
  dateString: z.string().refine((val) => !isNaN(Date.parse(val)), {
    message: "Invalid date format",
  }),
};
