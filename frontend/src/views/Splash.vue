<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Rocket } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/auth'
import { useBootstrapStore } from '@/stores/bootstrap'
import SplashProgress from '@/components/SplashProgress.vue'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const authStore = useAuthStore()
const bootstrapStore = useBootstrapStore()
const { t } = useI18n()

const isBackendReady = ref(false)
const isAutoRedirecting = ref(false)
const showProgress = computed(() => !isBackendReady.value || isAutoRedirecting.value)

// Typewriter effect state
const taglines = [
  'Empower your freelance journey.',
  '掌控自由职业，从这里开始。',
  'Track time. Manage clients. Invoice smart.',
  '时间追踪、客户管理、智能开票。',
  'Your work. Your rules.',
  '你的工作，你的规则。',
]

const currentTaglineIndex = ref(0)
const displayedText = ref('')
const isTyping = ref(true)
let typewriterInterval: ReturnType<typeof setInterval> | null = null
let pauseTimeout: ReturnType<typeof setTimeout> | null = null

const currentFullText = computed(() => taglines[currentTaglineIndex.value])

function startTypewriter() {
  let charIndex = 0
  isTyping.value = true
  displayedText.value = ''

  const fullText = currentFullText.value ?? ''

  typewriterInterval = setInterval(() => {
    if (charIndex < fullText.length) {
      displayedText.value += fullText[charIndex]
      charIndex++
    } else {
      // Finished typing this line
      if (typewriterInterval) clearInterval(typewriterInterval)
      typewriterInterval = null
      isTyping.value = false

      // Pause then start erasing
      pauseTimeout = setTimeout(() => {
        startEraser()
      }, 2500)
    }
  }, 80)
}

function startEraser() {
  isTyping.value = true

  typewriterInterval = setInterval(() => {
    if (displayedText.value.length > 0) {
      displayedText.value = displayedText.value.slice(0, -1)
    } else {
      // Finished erasing
      if (typewriterInterval) clearInterval(typewriterInterval)
      typewriterInterval = null
      isTyping.value = false

      // Move to next tagline
      currentTaglineIndex.value = (currentTaglineIndex.value + 1) % taglines.length

      // Small pause then start typing next
      pauseTimeout = setTimeout(() => {
        startTypewriter()
      }, 500)
    }
  }, 40)
}

onMounted(async () => {
  bootstrapStore.mark('splashMountedMs')
  // Start typewriter effect
  startTypewriter()

  // Initialize auth and check for existing session
  const authStart = typeof performance !== 'undefined' ? performance.now() : 0
  await authStore.initialize()
  if (typeof performance !== 'undefined') {
    bootstrapStore.mark('authInitMs', performance.now() - authStart)
  }

  // Backend is ready
  isBackendReady.value = true
  bootstrapStore.mark('totalToReadyMs')

  // If user is already authenticated, auto-redirect to dashboard
  if (authStore.isAuthenticated) {
    isAutoRedirecting.value = true
    setTimeout(() => {
      const last = localStorage.getItem("lastAuthedRoute");
      if (
        last &&
        last.startsWith("/") &&
        !["/splash", "/login", "/register"].includes(last) &&
        router.resolve(last).matched.length > 0
      ) {
        router.replace(last);
      } else {
        router.replace('/dashboard');
      }
    }, 600)
    return
  }

  // No prefetching needed - using static imports for desktop app
})

onUnmounted(() => {
  if (typewriterInterval) clearInterval(typewriterInterval)
  if (pauseTimeout) clearTimeout(pauseTimeout)
})

function handleStart() {
  if (authStore.usersList.length > 0) {
    // Users exist, go to login selection
    router.push('/login')
  } else {
    // No users, go to onboarding/registration
    router.push('/register')
  }
}
</script>

<template>
  <div class="splash-container" :class="{ 'auto-redirecting': isAutoRedirecting }">
    <!-- Background Image with Ken Burns effect -->
    <div class="splash-background" />

    <!-- Overlay for better text contrast -->
    <div class="splash-overlay" />

    <!-- Content -->
    <div class="splash-content">
      <!-- Logo / Brand -->
      <div class="brand-section">
        <h1 class="brand-title">FreelanceFlow</h1>
        <p class="brand-tagline">
          <span class="typewriter-text">{{ displayedText }}</span>
          <span class="cursor" :class="{ typing: isTyping }">|</span>
        </p>
      </div>

      <!-- Action Button / Auto-enter -->
      <div class="action-section">
        <SplashProgress v-if="showProgress" :is-ready="isBackendReady" :is-auto-redirecting="isAutoRedirecting" />
        <Button v-if="!showProgress && isBackendReady" size="lg" class="start-button rounded-full text-lg h-14 px-12"
          @click="handleStart">
          <Rocket class="mr-2 h-5 w-5" />
          <span>{{ t('splash.start') }}</span>
        </Button>
      </div>

      <!-- Version / Footer -->
      <div class="splash-footer">
        <span>v1.0.0</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.splash-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100%;
  overflow: hidden;
}

.splash-container.auto-redirecting {
  animation: splashFadeOut 0.6s ease forwards;
}

@keyframes splashFadeOut {
  to {
    opacity: 0;
    transform: scale(1.02);
  }
}

.splash-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url('/splash_bg.jpg');
  background-size: cover;
  background-position: center;
  animation: kenBurns 20s ease-in-out infinite alternate;
}

.splash-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: var(--gradient-overlay, linear-gradient(135deg,
        rgba(0, 0, 0, 0.6) 0%,
        rgba(0, 0, 0, 0.3) 50%,
        rgba(0, 0, 0, 0.6) 100%));
}

.splash-content {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: var(--space-10);
  text-align: center;
  color: var(--text-on-gradient);
}

.brand-section {
  margin-bottom: var(--space-12);
}

.brand-title {
  font-family: var(--font-display);
  font-size: var(--text-5xl);
  font-weight: 800;
  letter-spacing: -0.02em;
  margin: 0 0 var(--space-4) 0;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  animation: fadeInUp 1s ease-out;
}

.brand-tagline {
  font-size: var(--text-xl);
  opacity: 0.9;
  margin: 0;
  min-height: 1.6em;
  animation: fadeInUp 1s ease-out 0.2s both;
}

.typewriter-text {
  display: inline;
}

.cursor {
  display: inline-block;
  margin-left: 2px;
  animation: blink 1s step-end infinite;
}

.cursor.typing {
  animation: none;
  opacity: 1;
}

.action-section {
  min-height: 80px;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  align-items: center;
  justify-content: center;
}

.start-button {
  padding: 0 var(--space-12);
  height: 52px;
  font-size: var(--text-lg);
  font-weight: 600;
  min-width: 200px;
  box-shadow: var(--shadow-xl);
  transition: all var(--transition-slow);
}

.start-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4);
}

.start-button:disabled {
  opacity: 0.9;
  cursor: wait;
}

.splash-footer {
  position: absolute;
  bottom: var(--space-6);
  left: 0;
  right: 0;
  text-align: center;
  font-size: var(--text-sm);
  opacity: 0.6;
}
</style>
