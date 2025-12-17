<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { toTypedSchema } from '@vee-validate/zod'
import { api } from "@/api";
import { useI18n } from "vue-i18n";
import type { FinanceAccount, FinanceSettings } from "@/types/finance";
import { financeSettingsSchema } from "@/schemas/settings";

import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  SelectGroup
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { toast } from "vue-sonner";

const { t } = useI18n();

const accounts = ref<FinanceAccount[]>([]);
const settings = ref<FinanceSettings | null>(null);
const saving = ref(false);

const formSchema = toTypedSchema(financeSettingsSchema);

const initialValues = computed(() => {
  if (settings.value) {
    return {
      defaultAccountId: settings.value.defaultAccountId,
      autoCategorize: settings.value.autoCategorize,
      autoReconcile: settings.value.autoReconcile,
      userId: settings.value.userId,
    };
  }
  return {
    defaultAccountId: undefined,
    autoCategorize: true,
    autoReconcile: false,
    userId: 0,
  };
});

const accountOptions = computed(() =>
  accounts.value.map((a) => ({
    label: `${a.name} (${a.type})`,
    value: a.id,
  }))
);

onMounted(async () => {
  try {
    // Use API calls for testing
    accounts.value = await api.finance.accounts.list();
    settings.value = await api.finance.settings.get();
  } catch (e) {
    toast.error(t("settings.finance.messages.loadError"));
  }
});

async function onSubmit(formValues: unknown) {
  const values = formValues as any;
  saving.value = true;
  try {
    await api.finance.settings.update(values);
    toast.success(t("settings.finance.messages.saved"));
  } catch (e) {
    toast.error(e instanceof Error ? e.message : t("settings.finance.messages.saveError"));
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div class="finance-settings space-y-6">
    <Card>
      <CardHeader>
        <CardTitle>{{ t('settings.finance.title') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <Form :validation-schema="formSchema" :initial-values="initialValues" @submit="onSubmit" class="space-y-6"
          :key="settings ? 'loaded' : 'loading'">
          <!-- Default Account -->
          <FormField v-slot="{ componentField }" name="defaultAccountId">
            <FormItem>
              <FormLabel>{{ t('settings.finance.fields.defaultAccount') }}</FormLabel>
              <Select v-bind="componentField" :model-value="componentField.modelValue?.toString()">
                <FormControl>
                  <SelectTrigger :disabled="saving">
                    <SelectValue :placeholder="t('settings.finance.placeholders.defaultAccount')" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem v-for="opt in accountOptions" :key="opt.value" :value="opt.value.toString()">
                      {{ opt.label }}
                    </SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Auto Categorize -->
          <FormField v-slot="{ value, handleChange }" name="autoCategorize">
            <FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
              <div class="space-y-0.5">
                <FormLabel class="text-base">
                  {{ t('settings.finance.fields.autoCategorize') }}
                </FormLabel>
                <FormDescription>
                  {{ t("settings.finance.hints.autoCategorize") }}
                </FormDescription>
              </div>
              <FormControl>
                <Switch :checked="value" @update:checked="handleChange" :disabled="saving" />
              </FormControl>
            </FormItem>
          </FormField>

          <!-- Auto Reconcile -->
          <FormField v-slot="{ value, handleChange }" name="autoReconcile">
            <FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
              <div class="space-y-0.5">
                <FormLabel class="text-base">
                  {{ t('settings.finance.fields.autoReconcile') }}
                </FormLabel>
                <FormDescription>
                  {{ t("settings.finance.hints.autoReconcile") }}
                </FormDescription>
              </div>
              <FormControl>
                <Switch :checked="value" @update:checked="handleChange" :disabled="saving" />
              </FormControl>
            </FormItem>
          </FormField>

          <div class="flex justify-end">
            <Button type="submit" :disabled="saving">
              <span v-if="saving">{{ t("common.saving") }}</span>
              <span v-else>{{ t("common.save") }}</span>
            </Button>
          </div>
        </Form>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>
/* No scoped styles needed */
</style>
