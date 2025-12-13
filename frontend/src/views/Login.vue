<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NAvatar, NInput, NIcon, NSpace, NText, NForm, NFormItem, type FormInst } from 'naive-ui'
import { LockOutlined, PlusOutlined } from '@vicons/antd'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import type { UserListItem } from '@/types'
import { loginSchema } from '@/schemas/auth'
import { useZodRule } from '@/utils/validation'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()

const selectedUser = ref<UserListItem | null>(null)
const formRef = ref<FormInst | null>(null)
const formValue = reactive({
  password: ''
})

const rules = {
  password: useZodRule(loginSchema.shape.password)
}

const isLoggingIn = ref(false)
const loginError = ref<string | null>(null)

onMounted(() => {
  // If no users, redirect to register
  if (authStore.usersList.length === 0) {
    router.replace('/register')
  }
})

function selectUser(user: UserListItem) {
  selectedUser.value = user
  formValue.password = ''
  loginError.value = null
}

function cancelSelection() {
  selectedUser.value = null
  formValue.password = ''
  loginError.value = null
}

async function handleLogin() {
  if (!selectedUser.value) return

  try {
    await formRef.value?.validate()
  } catch {
    return // Validation failed
  }

  isLoggingIn.value = true
  loginError.value = null

  try {
    await authStore.login({
      username: selectedUser.value.username,
      password: formValue.password,
    })
    router.push('/dashboard')
  } catch (e) {
    loginError.value = t('auth.invalidPassword')
  } finally {
    isLoggingIn.value = false
  }
}

function goToRegister() {
  router.push('/register')
}
</script>

<template>
  <div class="login-container">
    <div class="login-card glass-card">
      <h1 class="auth-title">{{ t('auth.welcome') }}</h1>
      <p class="auth-subtitle">{{ t('auth.selectUser') }}</p>

      <!-- User Selection Grid -->
      <Transition name="fade" mode="out-in">
        <div v-if="!selectedUser" class="user-grid">
          <div v-for="user in authStore.usersList" :key="user.id" class="user-card" @click="selectUser(user)">
            <n-avatar :size="80" :src="user.avatarUrl"
              :fallback-src="`https://api.dicebear.com/9.x/avataaars/svg?seed=${user.username}`" />
            <span class="user-name">{{ user.username }}</span>
          </div>

          <!-- Add New User Card -->
          <div class="user-card add-user" @click="goToRegister">
            <div class="add-icon">
              <n-icon size="40">
                <PlusOutlined />
              </n-icon>
            </div>
            <span class="user-name">{{ t('auth.addUser') }}</span>
          </div>
        </div>

        <!-- Password Entry -->
        <div v-else class="password-section">
          <n-avatar :size="100" :src="selectedUser.avatarUrl"
            :fallback-src="`https://api.dicebear.com/9.x/avataaars/svg?seed=${selectedUser.username}`" />
          <h2 class="selected-username">{{ selectedUser.username }}</h2>

          <div class="password-form">
            <n-form ref="formRef" :model="formValue" :rules="rules" @submit.prevent="handleLogin">
              <n-form-item :show-label="false" path="password">
                <n-input v-model:value="formValue.password" type="password" :placeholder="t('auth.enterPassword')"
                  size="large" show-password-on="click" @keyup.enter="handleLogin">
                  <template #prefix>
                    <n-icon>
                      <LockOutlined />
                    </n-icon>
                  </template>
                </n-input>
              </n-form-item>

              <n-text v-if="loginError" type="error" class="login-error">{{ loginError }}</n-text>

              <n-space style="margin-top: 12px">
                <n-button type="primary" size="large" :loading="isLoggingIn" @click="handleLogin">
                  {{ t('auth.login') }}
                </n-button>
                <n-button size="large" @click="cancelSelection">
                  {{ t('common.cancel') }}
                </n-button>
              </n-space>
            </n-form>
          </div>
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4);
  min-height: 0;
  overflow: auto;
}

.login-card {
  max-width: 520px;
  width: 100%;
  padding: var(--space-8);
}

.user-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: var(--space-6);
  justify-items: center;
}

.user-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-5);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-normal);
  width: 120px;
}

.user-card:hover {
  background: rgba(0, 0, 0, 0.05);
  transform: translateY(-4px);
}

.user-card.add-user {
  border: 2px dashed var(--border-default);
}

.user-card.add-user:hover {
  border-color: var(--color-warm-orange);
}

.add-icon {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: var(--bg-base);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.user-card.add-user:hover .add-icon {
  background: var(--color-warm-orange);
  color: white;
}

.user-name {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--text-primary);
}

.password-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-5);
}

.selected-username {
  font-family: var(--font-heading);
  font-size: var(--text-2xl);
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.password-form {
  width: 100%;
  max-width: 300px;
}

.login-error {
  display: block;
  margin-top: 8px;
  text-align: center;
}
</style>
