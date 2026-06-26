/**
 * Register UI components globally so they can be used without importing.
 * Nuxt auto-imports from components/ dir, but the /ui/ subdirectory
 * prefixes them with 'Ui' by default. This plugin removes the prefix
 * so <AppButton> works rather than <UiAppButton>.
 */
export default defineNuxtPlugin(() => {})
