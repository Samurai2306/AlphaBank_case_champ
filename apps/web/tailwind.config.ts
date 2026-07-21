import type { Config } from "tailwindcss";

export default {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: { DEFAULT: "#EF3124", dark: "#C4291E" },
        ink: "#0B1117",
        canvas: "#F5F5F5",
        platinum: "#D1D5D8",
        mint: "#E8F6EE",
        lavender: "#EDE7F6",
        peach: "#FFF0E6",
        rose: "#FEE1E1",
        cyanSoft: "#E3F4F8",
      },
      borderRadius: {
        card: "28px",
        pill: "999px",
      },
      fontFamily: {
        sans: [
          "Alfa Interface Sans",
          "Manrope",
          "system-ui",
          "sans-serif",
        ],
      },
      boxShadow: {
        soft: "0 8px 30px rgba(11, 17, 23, 0.06)",
      },
    },
  },
  plugins: [],
} satisfies Config;
