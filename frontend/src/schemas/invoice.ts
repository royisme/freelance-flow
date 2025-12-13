import { z } from "zod";
import { commonSchemas } from "./common";

export const invoiceItemSchema = z.object({
  id: z.number(), // Usually frontend temp id or backend id
  description: z.string().min(1),
  quantity: z.number().min(0),
  unitPrice: z.number().min(0),
  amount: z.number(),
});

export const invoiceSchema = z.object({
  id: z.optional(commonSchemas.id),
  clientId: z.number().int().positive(),
  number: commonSchemas.requiredString,
  issueDate: commonSchemas.dateString,
  dueDate: commonSchemas.dateString,
  items: z.array(invoiceItemSchema),
  subtotal: z.number(),
  taxRate: z.number().min(0).max(1),
  taxAmount: z.number(),
  total: z.number(),
  status: z.enum(["draft", "sent", "paid", "overdue"]),
});
