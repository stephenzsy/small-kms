import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          react: ["react", "react-dom"],
          "react-router-dom": ["react-router-dom"],
          "msal": ["@azure/msal-browser", "@azure/msal-react"],
          //antd: ["antd"],
        },
      },
    },
  },
});
