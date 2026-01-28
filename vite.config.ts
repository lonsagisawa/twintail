import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
	plugins: [tailwindcss()],
	build: {
		outDir: "static/dist",
		emptyOutDir: true,
		manifest: true,
		rollupOptions: {
			input: {
				app: "assets/main.ts",
			},
		},
	},
	server: {
		strictPort: true,
		port: 5173,
		cors: true,
		origin: "http://localhost:5173",
		hmr: { host: "localhost", protocol: "ws", port: 5173 },
	},
});
