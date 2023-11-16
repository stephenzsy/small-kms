/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      spacing: {
        ex: "1ex",
        em: "1em",
      },
      fontFamily: {
        sans: ["Mona sans", "ui-sans-serif", "system-ui", "sans-serif"],
      },
    },
  },
  corePlugins: {
    preflight: false,
  },
};
