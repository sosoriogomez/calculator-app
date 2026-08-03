import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: { port: 5173, proxy: { "/api": "http://localhost:8080", "/healthz": "http://localhost:8080" } },
  test: { environment: "jsdom", setupFiles: "./src/test/setup.ts", globals: true, coverage: { provider: "v8", reporter: ["text", "html", "lcov"], include: ["src/**/*.{ts,tsx}"], thresholds: { lines: 90, functions: 90, branches: 90, statements: 90 }, exclude: ["src/main.tsx", "src/test/**", "**/*.test.*", "**/*.d.ts"] } }
});
