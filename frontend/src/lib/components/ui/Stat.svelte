<script lang="ts">
	import { Card, CardContent } from '$lib/components/ui';

	let {
		icon: Icon = undefined as typeof import('lucide-svelte').Library | undefined,
		label = '',
		value = 0,
	} = $props();
</script>

<Card class="gap-2 py-4">
	<CardContent class="flex items-center gap-3 px-4">
		<div class="rounded-lg bg-emerald-50 p-2 text-emerald-700">
			{#if Icon}<Icon class="h-5 w-5" />{/if}
		</div>
		<div>
			<p class="text-2xl font-bold text-slate-900">{value}</p>
			<p class="text-xs text-slate-500">{label}</p>
		</div>
	</CardContent>
</Card>
