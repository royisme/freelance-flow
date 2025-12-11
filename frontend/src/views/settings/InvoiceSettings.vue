<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useInvoiceEmailSettingsStore } from "@/stores/invoiceEmailSettings";
import { useMessage } from "naive-ui";
import type { InvoiceEmailSettings, EmailProvider } from "@/types";

const store = useInvoiceEmailSettingsStore();
const message = useMessage();

const form = ref<InvoiceEmailSettings>({
  provider: "mailto",
  subjectTemplate: "Invoice {{number}}",
  bodyTemplate: "Please find attached invoice {{number}}.",
  signature: "",
});

const saving = ref(false);

onMounted(async () => {
  await store.fetchSettings();
  if (store.settings) {
    form.value = { ...store.settings };
  }
});

const providerOptions: { label: string; value: EmailProvider }[] = [
  { label: "Mailto (default)", value: "mailto" },
  { label: "Resend", value: "resend" },
  { label: "SMTP", value: "smtp" },
];

async function handleSave() {
  saving.value = true;
  try {
    await store.saveSettings(form.value);
    message.success("Saved invoice email settings");
  } catch (e) {
    message.error(
      e instanceof Error ? e.message : "Failed to save invoice email settings"
    );
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div class="invoice-settings">
    <h2 class="section-title">Invoice Email Settings</h2>
    <n-form label-placement="top">
      <n-form-item label="Provider">
        <n-select
          v-model:value="form.provider"
          :options="providerOptions"
          :disabled="saving"
        />
      </n-form-item>

      <n-form-item label="From">
        <n-input v-model:value="form.from" :disabled="saving" />
      </n-form-item>

      <n-form-item label="Reply-To">
        <n-input v-model:value="form.replyTo" :disabled="saving" />
      </n-form-item>

      <n-form-item label="Subject Template">
        <n-input v-model:value="form.subjectTemplate" :disabled="saving" />
      </n-form-item>

      <n-form-item label="Body Template">
        <n-input
          type="textarea"
          v-model:value="form.bodyTemplate"
          :autosize="{ minRows: 3, maxRows: 6 }"
          :disabled="saving"
        />
      </n-form-item>

      <n-form-item label="Signature">
        <n-input
          type="textarea"
          v-model:value="form.signature"
          :autosize="{ minRows: 2, maxRows: 4 }"
          :disabled="saving"
        />
      </n-form-item>

      <template v-if="form.provider === 'resend'">
        <n-form-item label="Resend API Key">
          <n-input
            type="password"
            show-password-on="click"
            v-model:value="form.resendApiKey"
            :disabled="saving"
          />
        </n-form-item>
      </template>

      <template v-if="form.provider === 'smtp'">
        <n-form-item label="SMTP Host">
          <n-input v-model:value="form.smtpHost" :disabled="saving" />
        </n-form-item>
        <n-form-item label="SMTP Port">
          <n-input-number v-model:value="form.smtpPort" :disabled="saving" />
        </n-form-item>
        <n-form-item label="SMTP Username">
          <n-input v-model:value="form.smtpUsername" :disabled="saving" />
        </n-form-item>
        <n-form-item label="SMTP Password">
          <n-input
            type="password"
            show-password-on="click"
            v-model:value="form.smtpPassword"
            :disabled="saving"
          />
        </n-form-item>
        <n-form-item label="Use TLS">
          <n-switch v-model:value="form.smtpUseTLS" :disabled="saving" />
        </n-form-item>
      </template>

      <n-space justify="end">
        <n-button type="primary" :loading="saving" @click="handleSave">
          Save
        </n-button>
      </n-space>
    </n-form>
  </div>
</template>

<style scoped>
.invoice-settings {
  padding: 16px;
}

.section-title {
  font-size: 18px;
  margin-bottom: 12px;
}
</style>
