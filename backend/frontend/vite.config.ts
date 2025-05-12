/// <reference types="vite/client" />
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const isDev = mode === "development";

  return {
    plugins: [react()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },

    define: {
      // This ensures React runs in development mode when in dev mode
      "process.env.NODE_ENV": JSON.stringify(
        isDev ? "development" : "production",
      ),
      __DEV__: isDev,
      "React.Env.DEV": isDev,
    },
    ...(isDev
      ? {
          server: {
            // Only use proxy in development
            proxy: {
              "/api": {
                target: "http://127.0.0.1:8080",
                changeOrigin: true,
                // secure: false,
              },
            },
          },
          build: {
            sourcemap: true,
            minify: false,
          },
        }
      : {
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
        }),
  };
});
