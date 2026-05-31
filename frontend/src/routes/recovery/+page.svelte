<script lang="ts">
	import { Library } from 'lucide-svelte';
	import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui';
</script>

<svelte:head>
	<title>Account Recovery - Bayt al Hikmah</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 px-4">
	<div class="w-full max-w-md">
		<div class="mb-8 flex flex-col items-center">
			<div class="mb-4 rounded-xl bg-emerald-600 p-3"><Library class="h-8 w-8 text-white" /></div>
			<h1 class="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
		</div>

		<Card>
			<CardHeader class="text-center">
				<CardTitle>Account Recovery</CardTitle>
				<CardDescription>Password recovery is not implemented yet.</CardDescription>
			</CardHeader>
			<CardContent>
				<a href="/login"><Button class="w-full">Back to Sign In</Button></a>
			</CardContent>
		</Card>
	</div>
</div>
