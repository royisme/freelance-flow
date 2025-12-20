<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import {
  User, Lock, Mail, RefreshCw,
  ArrowLeft, ArrowRight, Check
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useI18n } from 'vue-i18n'
import { registerProfileSchema, registerPasswordBaseSchema, registerPreferencesSchema } from '@/schemas/auth'
import { toast } from 'vue-sonner'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const { t, locale } = useI18n()

// Step state
const currentStep = ref(1)
const totalSteps = 3

// Merged Schema for the whole wizard
const formSchema = toTypedSchema(
  registerProfileSchema
    .merge(registerPasswordBaseSchema)
    .merge(registerPreferencesSchema)
    .extend({
      confirmPassword: z.string(),
      avatarSeed: z.string()
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: t('auth.passwordsNotMatch'),
      path: ["confirmPassword"],
    })
)

// Form setup
const { handleSubmit, values, setFieldValue, validateField, setFieldError } = useForm({
  validationSchema: formSchema,
  initialValues: {
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    avatarSeed: Date.now().toString(),
    language: locale.value || 'zh-CN',
    currency: 'USD',
    timezone: (typeof Intl !== "undefined" && Intl.DateTimeFormat().resolvedOptions().timeZone) || "UTC"
  }
})

const isRegistering = ref(false)
const registerError = ref<string | null>(null)

// Options
const languageOptions = computed(() => [
  { label: t('settings.general.options.language.zhCN'), value: 'zh-CN' },
  { label: t('settings.general.options.language.enUS'), value: 'en-US' },
])

const currencyOptions = computed(() => [
  { label: `CAD - ${t('settings.general.options.currency.cad')}`, value: 'CAD' },
  { label: `USD - ${t('settings.general.options.currency.usd')}`, value: 'USD' },
  { label: `CNY - ${t('settings.general.options.currency.cny')}`, value: 'CNY' },
  { label: `EUR - ${t('settings.general.options.currency.eur')}`, value: 'EUR' },
])

const timezoneOptions = [
  { label: "UTC", value: "UTC" },
  { label: "Asia/Shanghai", value: "Asia/Shanghai" },
  { label: "America/Toronto", value: "America/Toronto" },
  { label: "America/New_York", value: "America/New_York" },
  { label: "Europe/London", value: "Europe/London" },
]

// Computed
const avatarUrl = computed(() => {
  return `https://api.dicebear.com/9.x/avataaars/svg?seed=${encodeURIComponent(values.avatarSeed || '')}`
})

// Functions
function regenerateAvatar() {
  setFieldValue('avatarSeed', Date.now().toString() + Math.random().toString(36).substring(7))
}

watch(() => values.language, (newLang) => {
  if (newLang) appStore.setLocale(newLang as 'zh-CN' | 'en-US')
})

async function nextStep() {
  let validGroup = false
  if (currentStep.value === 1) {
    const r1 = await validateField('username')
    const r2 = await validateField('email')
    if (r1.valid && r2.valid) validGroup = true
  } else if (currentStep.value === 2) {
    const r1 = await validateField('password')
    const r2 = await validateField('confirmPassword')
    // Manual password match check (refine rules don't run on field-level validation)
    if (r1.valid && r2.valid) {
      if (values.password === values.confirmPassword) {
        validGroup = true
      } else {
        setFieldError('confirmPassword', t('auth.passwordsNotMatch'))
      }
    }
  }

  if (validGroup) currentStep.value++
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

const handleRegister = handleSubmit(async (formValues) => {
  isRegistering.value = true
  registerError.value = null

  // Debug: Log all form data
  console.log('=== Registration Form Data ===')
  console.log('Username:', formValues.username)
  console.log('Email:', formValues.email)
  console.log('Password:', formValues.password ? '***' : '(empty)')
  console.log('Avatar URL:', avatarUrl.value)
  console.log('Language:', formValues.language)
  console.log('Currency:', formValues.currency)
  console.log('Timezone:', formValues.timezone)
  console.log('Avatar Seed:', formValues.avatarSeed)
  console.log('==============================')

  try {
    const settings = {
      language: formValues.language,
      currency: formValues.currency,
      timezone: formValues.timezone,
      theme: appStore.theme,
    }

    const registerPayload = {
      username: formValues.username,
      password: formValues.password,
      email: formValues.email || "",
      avatarUrl: avatarUrl.value,
      settingsJson: JSON.stringify(settings),
    }
    
    console.log('Register Payload:', { ...registerPayload, password: '***' })
    console.log('Settings JSON:', settings)

    const user = await authStore.register(registerPayload)
    
    console.log('Registration successful! User:', user)
    console.log('Is authenticated:', authStore.isAuthenticated)

    toast.success(t('auth.registerSuccess') || 'Registration successful!')
    router.push('/dashboard')
  } catch (e) {
    console.error('=== Registration Error ===')
    console.error('Error object:', e)
    console.error('Error message:', e instanceof Error ? e.message : String(e))
    console.error('==========================')
    
    const errorMessage = e instanceof Error ? e.message : t('auth.registerFailed')
    registerError.value = errorMessage
    toast.error(errorMessage)
    // NOTE: Do NOT jump to step 1, stay on current step and let user see the error
  } finally {
    isRegistering.value = false
  }
}, (submitErrors) => {
  // Validation errors - just log them, don't jump steps
  console.error('=== Validation Errors ===')
  console.error('Errors:', submitErrors.errors)
  console.error('Values:', submitErrors.values)
  console.error('=========================')
  
  const errors = submitErrors.errors
  const errorMessages = Object.values(errors || {}).filter(Boolean)
  
  if (errorMessages.length > 0) {
    toast.error(errorMessages[0] as string)
  } else {
    toast.error(t('auth.validationError') || 'Please check the form and try again')
  }
  // NOTE: Do NOT jump to any step, stay on current step
})
</script>

<template>
  <div class="h-full w-full flex overflow-y-auto p-4">
    <div class="glass-card max-w-[480px] w-full p-6 m-auto">
      <!-- Compact Header: Step Tabs Only -->
      <div class="flex justify-center mb-6">
        <!-- Inline Step Tabs -->
        <div class="flex items-center gap-4">
          <div v-for="step in 3" :key="step"
            class="flex items-center gap-2 cursor-default transition-colors duration-200" :class="{
              'text-primary font-semibold': currentStep === step,
              'text-muted-foreground font-medium': currentStep !== step,
              'text-muted-foreground/70': currentStep > step
            }">
            <div class="w-2 h-2 rounded-full transition-all duration-300" :class="{
              'bg-primary ring-2 ring-primary/20 w-2.5 h-2.5': currentStep === step,
              'bg-muted-foreground': currentStep !== step,
              'bg-primary': currentStep > step
            }" />
            <span class="text-sm">{{
              step === 1 ? t('auth.stepProfile') :
                step === 2 ? t('auth.stepSecurity') :
                  t('auth.stepPreferences')
            }}</span>
          </div>
        </div>
      </div>

      <!-- Step Content -->
      <div class="min-h-[240px]">

        <!-- Step 1: Profile -->
        <div v-show="currentStep === 1" class="text-center">
            <div class="relative inline-block mb-6 pt-2">
              <Avatar class="w-24 h-24 border-4 border-card shadow-lg">
                <AvatarImage :src="avatarUrl" />
                  <AvatarFallback>{{ t('common.user') }}</AvatarFallback>
              </Avatar>

              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button variant="outline" size="icon"
                      class="absolute bottom-0 -right-2 h-8 w-8 rounded-full shadow-md bg-card hover:bg-accent"
                      @click="regenerateAvatar">
                      <RefreshCw class="w-4 h-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {{ t('auth.regenerateAvatar') }}
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>

            <form @submit.prevent="nextStep" class="text-left space-y-4 max-w-sm mx-auto">
              <FormField v-slot="{ componentField }" name="username">
                <FormItem>
                  <FormLabel>{{ t('auth.username') }}</FormLabel>
                  <FormControl>
                    <div class="relative items-center">
                      <User class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input v-bind="componentField" :placeholder="t('auth.usernamePlaceholder')" class="pl-9" />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="email">
                <FormItem>
                  <FormLabel>{{ t('auth.email') }}</FormLabel>
                  <FormControl>
                    <div class="relative items-center">
                      <Mail class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input v-bind="componentField" :placeholder="t('auth.emailPlaceholder')" class="pl-9" />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </form>
        </div>

        <!-- Step 2: Password -->
        <div v-show="currentStep === 2" class="text-center">
            <h2 class="text-xl font-semibold mb-6 tracking-tight">{{ t('auth.setPassword') }}</h2>

            <form @submit.prevent="nextStep" class="text-left space-y-4 max-w-sm mx-auto">
              <FormField v-slot="{ componentField }" name="password">
                <FormItem>
                  <FormLabel>{{ t('auth.password') }}</FormLabel>
                  <FormControl>
                    <div class="relative items-center">
                      <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input type="password" v-bind="componentField" :placeholder="t('auth.passwordPlaceholder')"
                        class="pl-9" />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="confirmPassword">
                <FormItem>
                  <FormLabel>{{ t('auth.confirmPassword') }}</FormLabel>
                  <FormControl>
                    <div class="relative items-center">
                      <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                      <Input type="password" v-bind="componentField" :placeholder="t('auth.confirmPasswordPlaceholder')"
                        class="pl-9" />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </form>
        </div>

        <!-- Step 3: Preferences -->
        <div v-show="currentStep === 3" class="text-center">
            <h2 class="text-xl font-semibold mb-6 tracking-tight">{{ t('auth.financialPreferences') }}</h2>

            <form @submit.prevent="handleRegister" class="space-y-4">
              <FormField v-slot="{ componentField }" name="language">
                <FormItem>
                  <FormLabel>{{ t('auth.language') }}</FormLabel>
                  <Select v-bind="componentField">
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem v-for="opt in languageOptions" :key="opt.value" :value="opt.value">
                        {{ opt.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="currency">
                <FormItem>
                  <FormLabel>{{ t('auth.currency') }}</FormLabel>
                  <Select v-bind="componentField">
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem v-for="opt in currencyOptions" :key="opt.value" :value="opt.value">
                        {{ opt.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="timezone">
                <FormItem>
                  <FormLabel>{{ t('auth.timezone') }}</FormLabel>
                  <Select v-bind="componentField">
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue :placeholder="t('auth.timezonePlaceholder')" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem v-for="opt in timezoneOptions" :key="opt.value" :value="opt.value">
                        {{ opt.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              </FormField>
            </form>

            <p v-if="registerError" class="text-destructive text-sm mt-2 font-medium">
              {{ registerError }}
            </p>
        </div>

      </div>

      <!-- Navigation Buttons -->
      <div class="flex justify-center gap-4 mt-8">
        <!-- Previous Step Button -->
        <Button v-if="currentStep > 1" variant="outline" size="lg" class="min-w-[120px]" @click="prevStep">
          <ArrowLeft class="w-4 h-4 mr-2" />
          {{ t('common.prev') }}
        </Button>

        <!-- Next Step Button -->
        <Button v-if="currentStep < totalSteps" size="lg" class="min-w-[120px]" @click="nextStep">
          {{ t('common.next') }}
          <ArrowRight class="w-4 h-4 ml-2" />
        </Button>

        <Button v-else size="lg" class="min-w-[120px]" :disabled="isRegistering" @click="handleRegister">
          <RefreshCw v-if="isRegistering" class="w-4 h-4 mr-2 animate-spin" />
          <Check v-else class="w-4 h-4 mr-2" />
          {{ t('auth.createProfile') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Scoped styles removed in favor of Tailwind classes, keeping empty block or removing entirely. 
   Removing entirely as everything is utility classes now. 
*/
</style>
