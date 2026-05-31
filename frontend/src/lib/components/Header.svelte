<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { Button } from '$lib/components/ui';

	let {
		class: className = '',
		...rest
	} = $props();
</script>

<header class="sticky top-0 z-50 border-b border-slate-200 bg-white/80 px-4 backdrop-blur-lg" {...rest}>
	<nav class="mx-auto flex max-w-6xl flex-wrap items-center gap-x-3 gap-y-2 py-3 sm:py-4">
		<h2 class="m-0 flex-shrink-0 text-base font-semibold tracking-tight">
			<a
				href="/"
				class="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-slate-50 px-3 py-1.5 text-sm font-semibold text-slate-900 no-underline shadow-sm sm:px-4 sm:py-2"
			>
				<span class="h-2 w-2 rounded-full bg-gradient-to-r from-emerald-500 to-teal-400"></span>
				Bayt al Hikmah
			</a>
		</h2>

		<div class="ml-auto flex items-center gap-4">
			{#if auth.isAuthenticated}
				<a href="/dashboard"><Button variant="ghost" size="sm">Dashboard</Button></a>
				<a href="/settings"><Button variant="ghost" size="sm">Settings</Button></a>
			{:else}
				<a href="/login"><Button variant="ghost" size="sm">Sign In</Button></a>
				<a href="/registration"><Button size="sm" class="bg-emerald-600 hover:bg-emerald-700">Get Started</Button></a>
			{/if}
		</div>
	</nav>
</header>
