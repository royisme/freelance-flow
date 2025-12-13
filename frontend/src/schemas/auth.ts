import { z } from "zod";
import { commonSchemas } from "./common";

export const loginSchema = z.object({
  password: commonSchemas.requiredString,
});

export const registerProfileSchema = z.object({
  username: z.string().min(3),
  email: z.optional(z.union([z.literal(""), z.string().email()])), // Optional empty or valid email
});

export const registerPasswordBaseSchema = z.object({
  password: commonSchemas.password,
  confirmPassword: z.string(),
});

export const registerPasswordSchema = registerPasswordBaseSchema.refine(
  (data) => data.password === data.confirmPassword,
  {
    message: "Passwords don't match",
    path: ["confirmPassword"], // path of error
  }
);

export const registerPreferencesSchema = z.object({
  language: z.string(),
  currency: z.string(),
  timezone: z.string(),
});
