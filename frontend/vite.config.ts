import { sveltekit } from '@sveltejs/kit/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		SvelteKitPWA({
			registerType: 'autoUpdate',
			devOptions: { enabled: true },
			includeAssets: ['favicon.ico', 'logo192.png', 'logo512.png'],
			manifest: {
				name: 'Bayt al Hikmah',
				short_name: 'Bayt al Hikmah',
				description:
					'A modern, AI-powered knowledge management platform.',
				theme_color: '#173a40',
				background_color: '#ffffff',
				display: 'standalone',
				scope: '/',
				start_url: '/',
				icons: [
					{
						src: '/logo192.png',
						sizes: '192x192',
						type: 'image/png',
					},
					{
						src: '/logo512.png',
						sizes: '512x512',
						type: 'image/png',
					},
				],
			},
			workbox: {
				globPatterns: [
					'**/*.{js,css,html,ico,png,svg,webmanifest}',
				],
				cleanupOutdatedCaches: true,
				clientsClaim: true,
			},
		}),
	],
});
