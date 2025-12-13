<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton, NAvatar, NInput, NIcon, NText,
  NForm, NFormItem, NTooltip, NSelect, type FormInst
} from 'naive-ui'
import {
  UserOutlined, LockOutlined, MailOutlined,
  ReloadOutlined, ArrowLeftOutlined, ArrowRightOutlined, CheckOutlined
} from '@vicons/antd'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useI18n } from 'vue-i18n'
import { registerProfileSchema, registerPasswordBaseSchema, registerPreferencesSchema } from '@/schemas/auth'
import { useZodRule } from '@/utils/validation'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const { t, locale } = useI18n()

// Step state
const currentStep = ref(1)
const totalSteps = 3

// Form Refs
const step1FormRef = ref<FormInst | null>(null)
const step2FormRef = ref<FormInst | null>(null)
const step3FormRef = ref<FormInst | null>(null)

// Form state
const formValue = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  avatarSeed: Date.now().toString(),
  language: locale.value || 'zh-CN',
  currency: 'USD',
  timezone: (typeof Intl !== "undefined" && Intl.DateTimeFormat().resolvedOptions().timeZone) || "UTC"
})

const isRegistering = ref(false)
const registerError = ref<string | null>(null)

// Rules
const step1Rules = {
  username: useZodRule(registerProfileSchema.shape.username),
  email: useZodRule(registerProfileSchema.shape.email)
}

// Validation fix for confirm password with Zod logic
// We'll effectively duplicate the check or use a manual rule for this field to keep UI feedback simple
const confirmPasswordRule = {
  validator: (_rule: any, value: string) => {
    if (value !== formValue.password) {
      throw new Error(t('auth.passwordsNotMatch'))
    }
    return true
  },
  trigger: ['input', 'blur']
}// We can use Zod for the complexity check
const passwordRule = useZodRule(registerPasswordBaseSchema.shape.password)

const step3Rules = {
  // Preferences are usually select boxes with defaults, hard to fail, but good to have rules
  language: useZodRule(registerPreferencesSchema.shape.language),
  currency: useZodRule(registerPreferencesSchema.shape.currency),
  timezone: useZodRule(registerPreferencesSchema.shape.timezone)
}

// Options (Same as before)
const languageOptions = [
  { label: '中文 (简体)', value: 'zh-CN' },
  { label: 'English', value: 'en-US' }
]

const currencyOptions = [
  { label: 'CAD - Canadian Dollar', value: 'CAD' },
  { label: 'USD - US Dollar', value: 'USD' },
  { label: 'CNY - Chinese Yuan', value: 'CNY' },
  { label: 'EUR - Euro', value: 'EUR' },
]

const timezoneOptions = [
  { label: "UTC", value: "UTC" },
  { label: "Asia/Shanghai", value: "Asia/Shanghai" },
  { label: "America/Toronto", value: "America/Toronto" },
  { label: "America/New_York", value: "America/New_York" },
  { label: "Europe/London", value: "Europe/London" },
]

// Computed
const avatarUrl = computed(() => {
  return `https://api.dicebear.com/9.x/avataaars/svg?seed=${encodeURIComponent(formValue.avatarSeed)}`
})

// Functions
function regenerateAvatar() {
  formValue.avatarSeed = Date.now().toString() + Math.random().toString(36).substring(7)
}

watch(() => formValue.language, (newLang) => {
  appStore.setLocale(newLang as 'zh-CN' | 'en-US')
})

async function nextStep() {
  if (currentStep.value === 1) {
    try {
      await step1FormRef.value?.validate()
      currentStep.value++
    } catch {
      // Ignore
    }
  } else if (currentStep.value === 2) {
    try {
      if (formValue.password !== formValue.confirmPassword) throw new Error() // Double check
      await step2FormRef.value?.validate()
      currentStep.value++
    } catch {
      // Ignore
    }
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

async function handleRegister() {
  isRegistering.value = true
  registerError.value = null

  try {
    // Final validate just in case
    await step3FormRef.value?.validate()

    const settings = {
      language: formValue.language,
      currency: formValue.currency,
      timezone: formValue.timezone,
      theme: appStore.theme,
    }

    await authStore.register({
      username: formValue.username,
      password: formValue.password,
      email: formValue.email || "",
      avatarUrl: avatarUrl.value,
      settingsJson: JSON.stringify(settings),
    })

    router.push('/dashboard')
  } catch (e) {
    registerError.value = e instanceof Error ? e.message : t('auth.registerFailed')
  } finally {
    isRegistering.value = false
  }
}


</script>

<template>
  <div class="register-container">
    <div class="register-card glass-card">
      <!-- Compact Header: Step Tabs Only -->
      <div class="compact-header center-content">
        <!-- Inline Step Tabs -->
        <div class="step-tabs">
          <div v-for="step in 3" :key="step" class="step-tab" :class="{
            'tab-active': currentStep === step,
            'tab-completed': currentStep > step
          }">
            <span class="tab-dot" />
            <span class="tab-label">{{
              step === 1 ? t('auth.stepProfile') :
                step === 2 ? t('auth.stepSecurity') :
                  t('auth.stepPreferences')
            }}</span>
          </div>
        </div>
      </div>

      <!-- Step Content -->
      <div class="step-content">
        <!-- Step 1: Profile -->
        <Transition name="slide" mode="out-in">
          <div v-if="currentStep === 1" key="step1" class="step-panel">
            <div class="avatar-section">
              <n-avatar :size="80" :src="avatarUrl" class="avatar-preview" />
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button circle class="refresh-avatar" size="small" @click="regenerateAvatar">
                    <template #icon>
                      <n-icon>
                        <ReloadOutlined />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                {{ t('auth.regenerateAvatar') }}
              </n-tooltip>
            </div>

            <n-form ref="step1FormRef" :model="formValue" :rules="step1Rules" class="step-form">
              <n-form-item :label="t('auth.username')" path="username">
                <n-input v-model:value="formValue.username" :placeholder="t('auth.usernamePlaceholder')" size="large">
                  <template #prefix>
                    <n-icon>
                      <UserOutlined />
                    </n-icon>
                  </template>
                </n-input>
              </n-form-item>

              <n-form-item :label="t('auth.email')" path="email" :show-require-mark="false">
                <n-input v-model:value="formValue.email" :placeholder="t('auth.emailPlaceholder')" size="large">
                  <template #prefix>
                    <n-icon>
                      <MailOutlined />
                    </n-icon>
                  </template>
                </n-input>
              </n-form-item>
            </n-form>
          </div>

          <!-- Step 2: Password -->
          <div v-else-if="currentStep === 2" key="step2" class="step-panel">
            <h2 class="step-title">{{ t('auth.setPassword') }}</h2>

            <n-form ref="step2FormRef" :model="formValue"
              :rules="{ password: passwordRule, confirmPassword: confirmPasswordRule }" class="step-form">
              <n-form-item :label="t('auth.password')" path="password">
                <n-input v-model:value="formValue.password" type="password" :placeholder="t('auth.passwordPlaceholder')"
                  size="large" show-password-on="click">
                  <template #prefix>
                    <n-icon>
                      <LockOutlined />
                    </n-icon>
                  </template>
                </n-input>
              </n-form-item>

              <n-form-item :label="t('auth.confirmPassword')" path="confirmPassword">
                <n-input v-model:value="formValue.confirmPassword" type="password"
                  :placeholder="t('auth.confirmPasswordPlaceholder')" size="large" show-password-on="click">
                  <template #prefix>
                    <n-icon>
                      <LockOutlined />
                    </n-icon>
                  </template>
                </n-input>
              </n-form-item>
            </n-form>
          </div>

          <!-- Step 3: Preferences -->
          <div v-else-if="currentStep === 3" key="step3" class="step-panel">
            <h2 class="step-title">{{ t('auth.financialPreferences') }}</h2>

            <n-form ref="step3FormRef" :model="formValue" :rules="step3Rules" class="step-form">
              <n-form-item :label="t('auth.language')" path="language">
                <n-select v-model:value="formValue.language" :options="languageOptions" size="large" />
              </n-form-item>

              <n-form-item :label="t('auth.currency')" path="currency">
                <n-select v-model:value="formValue.currency" :options="currencyOptions" size="large" />
              </n-form-item>

              <n-form-item :label="t('auth.timezone')" path="timezone">
                <n-select v-model:value="formValue.timezone" :options="timezoneOptions" size="large" filterable />
              </n-form-item>
            </n-form>

            <n-text v-if="registerError" type="error" class="error-text">
              {{ registerError }}
            </n-text>
          </div>
        </Transition>
      </div>

      <!-- Navigation Buttons -->
      <div class="step-actions">
        <!-- Previous Step Button -->
        <n-button v-if="currentStep > 1" size="large" @click="prevStep">
          <template #icon>
            <n-icon>
              <ArrowLeftOutlined />
            </n-icon>
          </template>
          {{ t('common.prev') }}
        </n-button>

        <!-- Next Step Button -->
        <n-button v-if="currentStep < totalSteps" type="primary" size="large" @click="nextStep">
          {{ t('common.next') }}
          <template #icon>
            <n-icon>
              <ArrowRightOutlined />
            </n-icon>
          </template>
        </n-button>

        <n-button v-else type="primary" size="large" :loading="isRegistering" @click="handleRegister">
          <template #icon>
            <n-icon>
              <CheckOutlined />
            </n-icon>
          </template>
          {{ t('auth.createProfile') }}
        </n-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Container handles scrolling */
.register-container {
  height: 100%;
  width: 100%;
  display: flex;
  overflow-y: auto;
  padding: var(--space-4);
}

.register-card {
  max-width: 480px;
  width: 100%;
  padding: var(--space-6);
  /* Safe centering: pushes element to center if space allows, 
     but avoids top-clipping if content overflows */
  margin: auto;
}

/* Compact Header: Center Step Tabs */
.compact-header.center-content {
  display: flex;
  justify-content: center;
  margin-bottom: var(--space-4);
}

.back-button {
  flex-shrink: 0;
}

.header-spacer {
  width: 32px;
  /* Balance with back button */
  flex-shrink: 0;
}

/* Inline Step Tabs */
.step-tabs {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.step-tab {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: default;
}

.tab-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--border-default);
  transition: all var(--transition-normal);
}

.step-tab.tab-active .tab-dot {
  width: 8px;
  height: 8px;
  background: var(--color-warm-amber);
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.2);
}

.step-tab.tab-completed .tab-dot {
  background: var(--color-success);
}

.tab-label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--text-muted);
  transition: color var(--transition-normal);
}

.step-tab.tab-active .tab-label {
  color: var(--text-primary);
  font-weight: 600;
}

.step-tab.tab-completed .tab-label {
  color: var(--text-secondary);
}

.step-content {
  min-height: 240px;
}

.step-panel {
  text-align: center;
}

.step-title {
  font-family: var(--font-heading);
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--space-5) 0;
}

.avatar-section {
  position: relative;
  display: inline-block;
  margin-bottom: var(--space-5);
}

.avatar-preview {
  border: 4px solid var(--bg-card-solid);
  box-shadow: var(--shadow-lg);
}

.refresh-avatar {
  position: absolute;
  bottom: 0;
  right: -8px;
  background: var(--bg-card-solid);
  box-shadow: var(--shadow-md);
}

.step-form {
  text-align: left;
  max-width: 320px;
  margin: 0 auto;
}

.error-text {
  display: block;
  text-align: center;
  margin-top: var(--space-2);
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: var(--space-4);
  margin-top: var(--space-5);
}

.step-actions .n-button {
  min-width: 120px;
}

/* Step transition */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.25s ease-out;
}

.slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
