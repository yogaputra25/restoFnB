<template>
  <div class="min-h-screen bg-dark-900 flex items-center justify-center">
    <div class="bg-dark-700 p-8 rounded-lg w-full max-w-md">
      <h1 class="text-2xl font-bold text-primary-500 mb-6 text-center">restoFnB</h1>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm text-dark-200 mb-1">Chain ID</label>
          <input
            v-model="chainId"
            type="text"
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-400 focus:border-primary-500 outline-none"
            required
          />
        </div>

        <div>
          <label class="block text-sm text-dark-200 mb-1">Email</label>
          <input
            v-model="email"
            type="email"
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-400 focus:border-primary-500 outline-none"
            required
          />
        </div>

        <div>
          <label class="block text-sm text-dark-200 mb-1">Password</label>
          <input
            v-model="password"
            type="password"
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-400 focus:border-primary-500 outline-none"
            required
          />
        </div>

        <p v-if="error" class="text-red-400 text-sm">{{ error }}</p>

        <button
          type="submit"
          class="w-full py-2 rounded bg-primary-500 text-white font-semibold hover:bg-primary-600 transition"
        >
          Login
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()

const chainId = ref('')
const email = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  try {
    await auth.login(chainId.value, email.value, password.value)
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
