import { z } from "zod";
import type { FormItemRule } from "naive-ui";
import { makeZodI18nMap } from "zod-i18n-map";

export function setupZodI18n(i18n: any) {
  z.setErrorMap(
    makeZodI18nMap({
      t: ((key: string, options?: any) => i18n.global.t(key, options)) as any,
    })
  );
}

/**
 * Creates a Naive UI validator for a Zod schema field.
 * @param schema The Zod schema to validate against
 */
export function createZodValidator(schema: z.ZodType<any, any>): FormItemRule {
  return {
    validator: async (_rule, value) => {
      try {
        await schema.parseAsync(value);
      } catch (error) {
        if (error instanceof z.ZodError) {
          // Return the first error message
          // ZodError issues can be multiple, we just take the first one relevant to the value
          throw new Error(error.errors[0].message);
        }
        throw error;
      }
    },
    trigger: ["blur", "input", "change"],
  };
}

// Helper to create a rule array for Naive UI
export function useZodRule(schema: z.ZodType<any, any>) {
  return [createZodValidator(schema)];
}
