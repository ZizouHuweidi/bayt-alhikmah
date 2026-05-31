<script lang="ts">
	let {
		class: className = '',
		...rest
	} = $props();

	const year = new Date().getFullYear();
</script>

<footer class="mt-20 border-t border-slate-200 px-4 pb-14 pt-10" {...rest}>
	<div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 text-center sm:flex-row sm:text-left">
		<p class="m-0 text-sm text-slate-500">&copy; {year} Bayt al Hikmah. All rights reserved.</p>
		<p class="m-0 text-sm text-slate-500">House of Wisdom — Reimagined for the Modern Age</p>
	</div>
</footer>
