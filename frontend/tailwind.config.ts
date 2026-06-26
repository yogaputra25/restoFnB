import type { Config } from 'tailwindcss'

export default {
  content: [
    './components/**/*.{vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './app.vue',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#fdf4ef',  /* --color-primary at 5% */
          100: '#f6d9cf',
          200: '#efbeaf',
          300: '#e7a38f',
          400: '#d7886e',
          500: '#C65D3B',
          600: '#a84e31',
          700: '#8a3f27',
          800: '#6c301d',
          900: '#4e2113',
          DEFAULT: '#C65D3B',
        },
        secondary: {
          50: '#fef7ee',
          500: '#D97706',
          DEFAULT: '#D97706',
        },
        accent: {
          500: '#F4B400',
          DEFAULT: '#F4B400',
        },
        surface: {
          DEFAULT: '#FFFFFF',
          secondary: '#F5F3EF',
        },
        background: {
          DEFAULT: '#FAF8F5',
        },
        border: {
          DEFAULT: '#E8E3DD',
        },
        text: {
          primary: '#2F2F2F',
          secondary: '#6B6B6B',
        },
        success: {
          500: '#2E7D32',
          DEFAULT: '#2E7D32',
        },
        danger: {
          500: '#C0392B',
          DEFAULT: '#C0392B',
        },
        dark: {
          50: '#f5f5f5',
          100: '#eeeeee',
          200: '#e0e0e0',
          300: '#bdbdbd',
          400: '#9e9e9e',
          500: '#424242',
          600: '#303030',
          700: '#212121',
          800: '#1a1a1a',
          900: '#0d0d0d',
        },
      },
      fontFamily: {
        heading: ['Outfit', 'sans-serif'],
        body: ['Inter', 'sans-serif'],
      },
      fontSize: {
        hero: ['56px', { lineHeight: '1.1', fontWeight: '700' }],
        'page-title': ['36px', { lineHeight: '1.2', fontWeight: '700' }],
        section: ['28px', { lineHeight: '1.3', fontWeight: '600' }],
        'card-title': ['22px', { lineHeight: '1.3', fontWeight: '600' }],
      },
      borderRadius: {
        card: '20px',
        button: '14px',
        input: '14px',
        image: '18px',
      },
      boxShadow: {
        'soft-sm': '0 1px 3px 0 rgba(47, 47, 47, 0.06)',
        'soft-md': '0 4px 12px 0 rgba(47, 47, 47, 0.08)',
        'soft-lg': '0 8px 24px 0 rgba(47, 47, 47, 0.10)',
      },
      spacing: {
        18: '72px',
      },
    },
  },
  plugins: [],
} satisfies Config
