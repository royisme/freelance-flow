import { z } from "zod";
import { commonSchemas } from "./common";

export const timeEntrySchema = z.object({
  id: z.optional(commonSchemas.id),
  projectId: z.coerce.number().int().positive(),
  date: commonSchemas.dateString,
  startTime: z.optional(z.string()),
  endTime: z.optional(z.string()),
  durationSeconds: z.coerce.number().min(0),
  description: commonSchemas.requiredString,
  invoiced: z.boolean(),
  billingMode: z.enum(["hourly", "fixed"]).default("hourly"),
  manualAmount: z.number().nullable().optional(),
}).superRefine((data, ctx) => {
  if (data.billingMode === "hourly") {
    if (data.durationSeconds <= 0) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["durationSeconds"],
        message: "Duration must be greater than 0 for hourly billing",
      });
    }
    if (data.manualAmount != null) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["manualAmount"],
        message: "Manual amount is only available for fixed billing",
      });
    }
  }

  if (data.billingMode === "fixed") {
    if (data.manualAmount == null || data.manualAmount < 0) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["manualAmount"],
        message: "Manual amount must be 0 or greater for fixed billing",
      });
    }
  }
});
