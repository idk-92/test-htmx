import presetQuick from "franken-ui/shadcn-ui/preset-quick";

/** @type {import('tailwindcss').Config} */
module.exports = {
  presets: [presetQuick()],

  content: ["./app/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [],
};
