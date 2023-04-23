/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx}',
  ],
  theme: {
    safelist: ['animate-[fade-in_1s_ease-in-out]',],
    extend: {
      width:{
        "almost-full":"99%",
      },
      height:{
        "almost-full":"99%",
      }
    },
  },
  plugins: [],
}

