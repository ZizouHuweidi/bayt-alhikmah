<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { ArrowRight, Library } from 'lucide-svelte';
	import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '$lib/components/ui';

	let loginValue = $state('');
	let password = $state('');
	let error = $state('');
	let isSubmitting = $state(false);

	async function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';
		isSubmitting = true;
		try {
			await auth.login(loginValue, password);
			goto('/dashboard');
		} catch {
			error = 'Invalid email/username or password.';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Sign In - Bayt al Hikmah</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 px-4">
	<div class="w-full max-w-md">
		<div class="mb-8 flex flex-col items-center">
			<div class="mb-4 rounded-xl bg-emerald-600 p-3"><Library class="h-8 w-8 text-white" /></div>
			<h1 class="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
			<p class="mt-2 text-slate-600">Welcome back to your library</p>
		</div>

		<Card>
			<CardHeader class="text-center">
				<CardTitle>Sign In</CardTitle>
				<CardDescription>Sign in with your email or username.</CardDescription>
			</CardHeader>
			<CardContent>
				<form class="space-y-4" onsubmit={handleSubmit}>
					<div class="space-y-2">
						<Label for="login">Email or username</Label>
						<Input id="login" bind:value={loginValue} required />
					</div>
					<div class="space-y-2">
						<Label for="password">Password</Label>
						<Input id="password" type="password" bind:value={password} required />
					</div>
					{#if error}<p class="text-sm text-red-600">{error}</p>{/if}
					<Button class="w-full" size="lg" disabled={isSubmitting}>
						Sign In
						<ArrowRight class="ml-2 h-4 w-4" />
					</Button>
				</form>
			</CardContent>
		</Card>
	</div>
</div>
