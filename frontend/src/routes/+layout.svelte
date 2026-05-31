<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { auth } from '$lib/auth.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();

	onMount(() => {
		auth.refreshSession();
	});
</script>

<svelte:head>
	<title>Bayt al Hikmah - House of Wisdom</title>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<Header />
{@render children()}
<Footer />
