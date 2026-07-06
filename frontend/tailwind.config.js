/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{vue,ts}"],
  theme: {
    extend: {
      colors: {
        ink: "#202833",
        panel: "#f5f7fb",
        line: "#dfe6ef"
      }
    }
  },
  plugins: []
};
