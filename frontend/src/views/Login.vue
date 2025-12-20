<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { Plus, Lock } from 'lucide-vue-next'
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
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { loginSchema } from '@/schemas/auth'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()

const selectedUserId = ref<string>('')
const isLoggingIn = ref(false)
const loginError = ref<string | null>(null)

// Selected user computed from ID
const selectedUser = computed(() => {
  if (!selectedUserId.value) return null
  return authStore.usersList.find(u => String(u.id) === selectedUserId.value) || null
})

// Avatar URL for selected user
const avatarUrl = computed(() => {
  if (!selectedUser.value) return ''
  return selectedUser.value.avatarUrl ||
    `https://api.dicebear.com/9.x/avataaars/svg?seed=${selectedUser.value.username}`
})

// Form validation schema - only validate password here
const formSchema = toTypedSchema(z.object({
  password: loginSchema.shape.password
}))

const { handleSubmit, resetForm } = useForm({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
  if (!selectedUser.value) return

  isLoggingIn.value = true
  loginError.value = null

  try {
    await authStore.login({
      username: selectedUser.value.username,
      password: values.password,
    })
    router.push('/dashboard')
  } catch (e) {
    loginError.value = t('auth.invalidPassword')
  } finally {
    isLoggingIn.value = false
  }
})

onMounted(() => {
  if (authStore.usersList.length === 0) {
    router.replace('/register')
  } else if (authStore.usersList.length === 1 && authStore.usersList[0]) {
    // Auto-select the only user
    selectedUserId.value = String(authStore.usersList[0].id)
  }
})

// Reset error when user changes
watch(selectedUserId, () => {
  loginError.value = null
  resetForm()
})

function goToRegister() {
  router.push('/register')
}
</script>

<template>
  <div class="h-full flex items-center justify-center p-4 min-h-0 overflow-auto">
    <div class="max-w-[400px] w-full p-8 glass-card">
      <h1 class="font-(family-name:--font-heading) text-3xl font-bold mb-2 text-center text-foreground">
        {{ t('auth.welcome') }}
      </h1>
      <p class="text-muted-foreground text-center mb-8">{{ t('auth.selectUser') }}</p>

      <!-- Login Form -->
      <form @submit="onSubmit" class="space-y-6">
        <!-- User Select with Avatar Preview -->
        <div class="flex flex-col items-center gap-4">
          <!-- Avatar Preview -->
          <Avatar class="h-20 w-20 border-2 border-muted">
            <AvatarImage v-if="selectedUser" :src="avatarUrl" :alt="selectedUser.username" />
            <AvatarFallback class="text-2xl text-muted-foreground">?</AvatarFallback>
          </Avatar>

          <!-- User Dropdown -->
          <div class="w-full">
            <Select v-model="selectedUserId">
              <SelectTrigger class="w-full h-12">
                <SelectValue :placeholder="t('auth.selectUser')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="user in authStore.usersList" :key="user.id" :value="String(user.id)">
                  <div class="flex items-center gap-2">
                    <Avatar class="h-6 w-6">
                      <AvatarImage
                        :src="user.avatarUrl || `https://api.dicebear.com/9.x/avataaars/svg?seed=${user.username}`"
                        :alt="user.username" />
                      <AvatarFallback class="text-xs">{{ user.username.charAt(0).toUpperCase() }}</AvatarFallback>
                    </Avatar>
                    <span>{{ user.username }}</span>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <!-- Password Field -->
        <div v-show="selectedUser">
          <FormField v-slot="{ componentField }" name="password">
            <FormItem>
              <FormLabel>{{ t('auth.password') }}</FormLabel>
              <FormControl>
                <div class="relative items-center">
                  <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input type="password" :placeholder="t('auth.passwordPlaceholder')" v-bind="componentField"
                    class="pl-10" />
                </div>
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <p v-if="loginError" class="text-destructive text-center text-sm font-medium mt-2">
            {{ loginError }}
          </p>

          <Button type="submit" class="w-full mt-4" size="lg" :disabled="isLoggingIn || !selectedUser">
            {{ isLoggingIn ? t('auth.loggingIn') : t('auth.login') }}
          </Button>
        </div>
      </form>

      <!-- Divider -->
      <div class="flex items-center gap-4 my-6">
        <div class="flex-1 h-px bg-border" />
        <span class="text-sm text-muted-foreground">{{ t('auth.or') }}</span>
        <div class="flex-1 h-px bg-border" />
      </div>

      <!-- Register Button -->
      <Button variant="outline" size="lg" class="w-full" @click="goToRegister">
        <Plus class="h-5 w-5 mr-2" />
        {{ t('auth.createAccount') }}
      </Button>
    </div>
  </div>
</template>
