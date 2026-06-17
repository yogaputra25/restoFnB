<template>
  <div
    class="relative overflow-hidden select-none rounded-xl h-[180px] sm:h-[220px]"
    @mouseenter="isPaused = true"
    @mouseleave="isPaused = false"
    @touchstart.passive="onTouchStart"
    @touchend.passive="onTouchEnd"
  >
    <Transition name="slide-left">
      <div :key="currentSlide" class="w-full h-full">
        <div
          v-if="promos[currentSlide]?.image"
          class="relative w-full h-full flex items-center justify-center px-6"
          :style="{ backgroundImage: `url(${promos[currentSlide].image})`, backgroundSize: 'cover', backgroundPosition: 'center' }"
        >
          <div class="absolute inset-0 bg-gradient-to-b from-black/40 to-black/60" />
          <div class="relative text-center text-white">
            <p class="text-xl sm:text-2xl font-bold">{{ promos[currentSlide]?.title }}</p>
            <p class="text-sm sm:text-base opacity-90 mt-1">{{ promos[currentSlide]?.description }}</p>
          </div>
        </div>
        <div
          v-else
          :class="['w-full h-full flex items-center justify-center px-6 text-center text-white', promos[currentSlide]?.bgColor || 'bg-primary-500']"
        >
          <div>
            <p class="text-xl sm:text-2xl font-bold">{{ promos[currentSlide]?.title }}</p>
            <p class="text-sm sm:text-base opacity-90 mt-1">{{ promos[currentSlide]?.description }}</p>
          </div>
        </div>
      </div>
    </Transition>

    <div v-if="promos.length > 1" class="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
      <button
        v-for="(promo, i) in promos"
        :key="i"
        @click="goTo(i)"
        :class="[
          'w-2 h-2 rounded-full transition-all duration-300',
          i === currentSlide ? 'bg-white w-4' : 'bg-white/50 hover:bg-white/70'
        ]"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Promo {
  title: string
  description: string
  image?: string
  bgColor?: string
}

const props = defineProps<{ promos: Promo[] }>()

const currentSlide = ref(0)
const isPaused = ref(false)

function advance() {
  if (props.promos.length === 0) return
  currentSlide.value = (currentSlide.value + 1) % props.promos.length
}

function goTo(i: number) {
  if (i === currentSlide.value) return
  currentSlide.value = i
}

let touchStartX = 0
function onTouchStart(e: TouchEvent) {
  touchStartX = e.touches[0].clientX
  isPaused.value = true
}
function onTouchEnd(e: TouchEvent) {
  const diff = e.changedTouches[0].clientX - touchStartX
  if (Math.abs(diff) > 50) {
    advance()
  }
  isPaused.value = false
}

onMounted(() => {
  const interval = setInterval(() => {
    if (!isPaused.value) advance()
  }, 6000)
  onUnmounted(() => clearInterval(interval))
})
</script>

<style scoped>
.slide-left-enter-active,
.slide-left-leave-active {
  transition: all 0.4s ease;
  position: absolute;
  inset: 0;
}
.slide-left-enter-from {
  transform: translateX(100%);
}
.slide-left-leave-to {
  transform: translateX(-100%);
}
</style>
