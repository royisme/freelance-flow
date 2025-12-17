<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { toTypedSchema } from '@vee-validate/zod'
import { useI18n } from "vue-i18n";
import { useInvoiceEmailSettingsStore } from "@/stores/invoiceEmailSettings";
import { invoiceEmailSettingsSchema } from "@/schemas/settings";

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
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  SelectGroup
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { Info } from 'lucide-vue-next'
import { toast } from "vue-sonner";

const store = useInvoiceEmailSettingsStore();
const { t } = useI18n();

const saving = ref(false);

const formSchema = toTypedSchema(invoiceEmailSettingsSchema);

const initialValues = computed(() => {
  if (store.settings) {
    return {
      ...store.settings,
      smtpUseTls: store.settings.smtpUseTls || false,
    };
  }
  return {
    provider: "mailto",
    subjectTemplate: t("settings.email.defaults.subjectTemplate"),
    bodyTemplate: t("settings.email.defaults.bodyTemplate"),
    signature: "",
    smtpUseTls: false,
  };
});

onMounted(async () => {
  await store.fetchSettings();
});

const providerOptions = [
  { label: t("settings.email.options.provider.mailto"), value: "mailto" },
  { label: t("settings.email.options.provider.resend"), value: "resend" },
  { label: t("settings.email.options.provider.smtp"), value: "smtp" },
];

async function onSubmit(formValues: unknown) {
  const values = formValues as any;
  try {
    saving.value = true;
    await store.saveSettings(values);
    toast.success(t("settings.email.messages.saved"));
  } catch (e) {
    toast.error(e instanceof Error ? e.message : t("settings.email.messages.saveError"));
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div class="email-settings space-y-6">
    <Card>
      <CardHeader>
        <CardTitle>{{ t('settings.email.providerCardTitle') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <Form :validation-schema="formSchema" :initial-values="initialValues" @submit="onSubmit" class="space-y-6"
          :key="store.settings ? 'loaded' : 'loading'" v-slot="{ values, setFieldValue }">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Provider -->
            <FormField v-slot="{ componentField }" name="provider">
              <FormItem class="col-span-1 md:col-span-2">
                <FormLabel>{{ t('settings.email.fields.provider') }}</FormLabel>
                <Select v-bind="componentField">
                  <FormControl>
                    <SelectTrigger :disabled="saving">
                      <SelectValue :placeholder="t('settings.email.placeholders.provider')" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem v-for="opt in providerOptions" :key="opt.value" :value="opt.value">
                        {{ opt.label }}
                      </SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            </FormField>

            <!-- Resend Specific -->
            <template v-if="values.provider === 'resend'">
              <FormField v-slot="{ componentField }" name="resendApiKey">
                <FormItem>
                  <FormLabel>{{ t('settings.email.fields.resendApiKey') }}</FormLabel>
                  <FormControl>
                    <Input type="password" v-bind="componentField" :disabled="saving" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </template>

            <!-- Resend Tip -->
            <div v-if="values.provider === 'resend'" class="col-span-1 md:col-span-2">
              <Alert>
                <Info class="h-4 w-4" />
                <AlertTitle>{{ t('settings.email.resendCardTitle') }}</AlertTitle>
                <AlertDescription>
                  {{ t("settings.email.resendDomainTip") }}
                </AlertDescription>
              </Alert>
            </div>

            <!-- SMTP Settings -->
            <template v-if="values.provider === 'smtp'">
              <FormField v-slot="{ componentField }" name="smtpHost">
                <FormItem>
                  <FormLabel>{{ t('settings.email.fields.smtpHost') }}</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" :disabled="saving" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="smtpPort">
                <FormItem>
                  <FormLabel>{{ t('settings.email.fields.smtpPort') }}</FormLabel>
                  <FormControl>
                    <Input type="number" v-bind="componentField" :disabled="saving"
                      @input="(e: Event) => setFieldValue('smtpPort', parseInt((e.target as HTMLInputElement).value))" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="smtpUsername">
                <FormItem>
                  <FormLabel>{{ t('settings.email.fields.smtpUsername') }}</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" :disabled="saving" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="smtpPassword">
                <FormItem>
                  <FormLabel>{{ t('settings.email.fields.smtpPassword') }}</FormLabel>
                  <FormControl>
                    <Input type="password" v-bind="componentField" :disabled="saving" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ value, handleChange }" name="smtpUseTls">
                <FormItem
                  class="flex flex-row items-center justify-between rounded-lg border p-4 col-span-1 md:col-span-2">
                  <div class="space-y-0.5">
                    <FormLabel class="text-base">
                      {{ t('settings.email.fields.smtpUseTLS') }}
                    </FormLabel>
                  </div>
                  <FormControl>
                    <Switch :checked="value" @update:checked="handleChange" :disabled="saving" />
                  </FormControl>
                </FormItem>
              </FormField>
            </template>

            <!-- Common Fields -->
            <FormField v-slot="{ componentField }" name="from">
              <FormItem>
                <FormLabel>{{ t('settings.email.fields.from') }}</FormLabel>
                <FormControl>
                  <Input v-bind="componentField" :disabled="saving" :placeholder="t('settings.email.placeholders.from')" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="replyTo">
              <FormItem>
                <FormLabel>{{ t('settings.email.fields.replyTo') }}</FormLabel>
                <FormControl>
                  <Input v-bind="componentField" :disabled="saving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="subjectTemplate">
              <FormItem>
                <FormLabel>{{ t('settings.email.fields.subjectTemplate') }}</FormLabel>
                <FormControl>
                  <Input v-bind="componentField" :disabled="saving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="signature">
              <FormItem>
                <FormLabel>{{ t('settings.email.fields.signature') }}</FormLabel>
                <FormControl>
                  <Textarea v-bind="componentField" :rows="3" :disabled="saving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="bodyTemplate">
              <FormItem class="col-span-1 md:col-span-2">
                <FormLabel>{{ t('settings.email.fields.bodyTemplate') }}</FormLabel>
                <FormControl>
                  <Textarea v-bind="componentField" :rows="6" :disabled="saving" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>

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
/* No styles needed with Tailwind */
</style>
