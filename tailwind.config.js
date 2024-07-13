/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,templ,go}"],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light"]
  }
};