<template>
  <div class="min-h-screen bg-[var(--color-background)] flex items-center justify-center px-4">
    <div class="bg-[var(--color-surface)] p-8 rounded-[var(--radius-card)] w-full max-w-md shadow-[var(--shadow-md)]" data-aos="fade-up">
      <h1 class="font-heading text-page-title text-[var(--color-primary)] mb-8 text-center">restoFnB</h1>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <AppInput
          v-model="email"
          label="Email"
          type="email"
          floating
          required
        />

        <AppInput
          v-model="password"
          label="Password"
          type="password"
          floating
          required
          :error="error"
        />

        <AppButton variant="primary" size="lg" type="submit" class="w-full">
          Login
        </AppButton>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  try {
    await auth.login(email.value, password.value)
    if (auth.isAdmin) {
      router.push('/admin')
    } else if (auth.isCashier) {
      router.push('/cashier')
    } else {
      router.push('/')
    }
  } catch {
    error.value = 'Login failed. Check your credentials.'
  }
}
</script>
