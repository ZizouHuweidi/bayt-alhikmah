<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { ArrowRight, Library } from 'lucide-svelte';
	import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '$lib/components/ui';

	let email = $state('');
	let username = $state('');
	let password = $state('');
	let error = $state('');
	let isSubmitting = $state(false);

	async function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';
		isSubmitting = true;
		try {
			await auth.register(email, username, password);
			goto('/dashboard');
		} catch {
			error = 'Registration failed. Check your details or try another email/username.';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Create Account - Bayt al Hikmah</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 px-4">
	<div class="w-full max-w-md">
		<div class="mb-8 flex flex-col items-center">
			<div class="mb-4 rounded-xl bg-emerald-600 p-3"><Library class="h-8 w-8 text-white" /></div>
			<h1 class="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
			<p class="mt-2 text-slate-600">Create your library account</p>
		</div>

		<Card>
			<CardHeader class="text-center">
				<CardTitle>Create Account</CardTitle>
				<CardDescription>Email, username, and a password of at least 12 characters are required.</CardDescription>
			</CardHeader>
			<CardContent>
				<form class="space-y-4" onsubmit={handleSubmit}>
					<div class="space-y-2">
						<Label for="email">Email</Label>
						<Input id="email" type="email" bind:value={email} required />
					</div>
					<div class="space-y-2">
						<Label for="username">Username</Label>
						<Input id="username" bind:value={username} required minlength={3} maxlength={32} />
					</div>
					<div class="space-y-2">
						<Label for="password">Password</Label>
						<Input id="password" type="password" bind:value={password} required minlength={12} />
					</div>
					{#if error}<p class="text-sm text-red-600">{error}</p>{/if}
					<Button class="w-full" size="lg" disabled={isSubmitting}>
						Create Account
						<ArrowRight class="ml-2 h-4 w-4" />
					</Button>
				</form>
			</CardContent>
		</Card>
	</div>
</div>
