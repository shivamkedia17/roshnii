import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const isProd = mode === "production";

  return {
    plugins: [react()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      // Only use proxy in development
      proxy: {
        "/api": {
          target: "http://127.0.0.1:8080",
          changeOrigin: true,
          secure: false,
        },
      },
    },
    // Production optimizations
    build: {
      sourcemap: false,
      minify: "terser",
      terserOptions: {
        compress: {
          drop_console: true,
        },
      },
    },
  };
});
