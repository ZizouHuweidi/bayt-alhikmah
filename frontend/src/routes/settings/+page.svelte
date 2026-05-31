<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Library, Loader2, Mail, Shield, User } from 'lucide-svelte';
	import { Button, Card, CardContent, Input } from '$lib/components/ui';
	import type { Profile } from '$lib/api';

	let profile = $state<Profile | null>(null);
	let displayName = $state('');
	let bio = $state('');
	let publicProfile = $state(false);
	let loadingProfile = $state(true);
	let saving = $state(false);
	let message = $state<string | null>(null);
	let error = $state<string | null>(null);

	onMount(() => {
		if (!auth.isAuthenticated) {
			goto('/login');
		}
	});

	$effect(() => {
		if (!auth.isAuthenticated || !auth.accessToken) return;
		let cancelled = false;
		async function loadProfile() {
			loadingProfile = true;
			error = null;
			try {
				const { getProfile } = await import('$lib/api');
				const result = await getProfile(auth.accessToken!);
				if (cancelled) return;
				profile = result;
				displayName = result.display_name || '';
				bio = result.bio || '';
				publicProfile = result.public_profile;
			} catch (err) {
				if (!cancelled) error = err instanceof Error ? err.message : 'Failed to load profile';
			} finally {
				if (!cancelled) loadingProfile = false;
			}
		}
		loadProfile();
		return () => { cancelled = true; };
	});

	async function handleSave(event: Event) {
		event.preventDefault();
		if (!auth.accessToken) return;
		saving = true;
		error = null;
		message = null;
		try {
			const { updateProfile } = await import('$lib/api');
			const updated = await updateProfile(auth.accessToken, {
				display_name: displayName.trim() || undefined,
				bio: bio.trim() || undefined,
				public_profile: publicProfile,
			});
			profile = updated;
			message = 'Profile saved';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save profile';
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Settings - Bayt al Hikmah</title>
</svelte:head>

<div class="min-h-screen bg-slate-50 py-8">
	<div class="mx-auto max-w-2xl px-4">
		<div class="mb-8 flex items-center justify-between gap-3">
			<div class="flex items-center gap-3">
				<div class="rounded-xl bg-emerald-600 p-2"><Library class="h-6 w-6 text-white" /></div>
				<h1 class="text-2xl font-bold text-slate-900">Account Settings</h1>
			</div>
			<a href="/dashboard"><Button variant="outline">Dashboard</Button></a>
		</div>

		{#if error}
			<div class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
		{/if}
		{#if message}
			<div class="mb-6 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">{message}</div>
		{/if}

		<Card>
			<CardContent class="space-y-6 pt-6">
				<div class="flex items-center gap-4">
					<div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100"><User class="h-6 w-6 text-emerald-600" /></div>
					<div><p class="font-medium text-slate-900">{auth.user.username}</p><p class="text-sm text-slate-500">Username</p></div>
				</div>
				<div class="flex items-center gap-4">
					<div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100"><Mail class="h-6 w-6 text-emerald-600" /></div>
					<div><p class="font-medium text-slate-900">{auth.user.email || 'No email'}</p><p class="text-sm text-slate-500">Email address</p></div>
				</div>
				<div class="flex items-center gap-4">
					<div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100"><Shield class="h-6 w-6 text-emerald-600" /></div>
					<div><p class="font-medium text-slate-900">Password auth enabled</p><p class="text-sm text-slate-500">OAuth support can be added later.</p></div>
				</div>

				<form class="space-y-4 border-t border-slate-200 pt-6" onsubmit={handleSave}>
					<div>
						<label for="display-name" class="mb-2 block text-sm font-medium text-slate-700">Display name</label>
						<Input id="display-name" bind:value={displayName} placeholder="How your public profile should appear" disabled={loadingProfile} />
					</div>
					<div>
						<label for="bio" class="mb-2 block text-sm font-medium text-slate-700">Bio</label>
						<textarea class="min-h-28 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm" bind:value={bio} placeholder="What are you reading, researching, or collecting?" disabled={loadingProfile}></textarea>
					</div>
					<label class="flex items-start gap-3 rounded-lg border border-slate-200 bg-slate-50 p-4 text-sm text-slate-700">
						<input type="checkbox" class="mt-1" bind:checked={publicProfile} disabled={loadingProfile} />
						<span>
							<span class="block font-medium text-slate-900">Make my profile public</span>
							Public profiles can be opened at <code>/users/{auth.user.username}/profile</code> and can anchor public notes, reviews, collections, and library items.
						</span>
					</label>
					{#if profile?.public_profile && auth.user.username}
						<p class="text-sm text-slate-500">
							Public profile: <a href="/users/{auth.user.username}/profile" class="font-medium text-emerald-700 underline-offset-4 hover:underline">/users/{auth.user.username}/profile</a>
						</p>
					{/if}
					<Button type="submit" disabled={saving || loadingProfile}>{saving ? 'Saving...' : 'Save profile'}</Button>
				</form>
			</CardContent>
		</Card>
	</div>
</div>
