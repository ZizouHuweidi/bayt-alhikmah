<script lang="ts">
	import { BookOpen, Layers3, Library, Loader2, Star, StickyNote } from 'lucide-svelte';
	import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Stat } from '$lib/components/ui';
	import type { Collection, LibraryItemWithSource, Note, Profile, Review } from '$lib/api';

	let { username } = $props();

	type PublicData = {
		profile: Profile | null;
		library: LibraryItemWithSource[];
		notes: Note[];
		reviews: Review[];
		collections: Collection[];
	};

	let data = $state<PublicData>({ profile: null, library: [], notes: [], reviews: [], collections: [] });
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		let cancelled = false;
		async function loadProfile() {
			loading = true;
			error = null;
			try {
				const apis = await import('$lib/api');
				const profile = await apis.getPublicProfile(username);
				const [library, notes, reviews, collections] = await Promise.all([
					apis.listPublicLibrary(username),
					apis.listPublicNotesByUser(profile.user_id),
					apis.listPublicReviewsByUser(profile.user_id),
					apis.listPublicCollectionsByUser(profile.user_id),
				]);
				if (!cancelled) data = { profile, library, notes, reviews, collections };
			} catch (err) {
				if (!cancelled) error = err instanceof Error ? err.message : 'Failed to load profile';
			} finally {
				if (!cancelled) loading = false;
			}
		}
		loadProfile();
		return () => { cancelled = true; };
	});
</script>

<svelte:head>
	<title>{data.profile?.display_name || data.profile?.username || username} - Bayt al Hikmah</title>
</svelte:head>

{#if loading}
	<main class="flex min-h-screen items-center justify-center bg-slate-50">
		<Loader2 class="h-8 w-8 animate-spin text-emerald-600" />
	</main>
{:else if error || !data.profile}
	<main class="min-h-screen bg-slate-50 px-4 py-16">
		<div class="mx-auto max-w-2xl rounded-2xl border border-slate-200 bg-white p-8 text-center shadow-sm">
			<h1 class="text-2xl font-bold text-slate-900">Profile unavailable</h1>
			<p class="mt-3 text-slate-600">This profile is private or does not exist.</p>
			<a href="/"><Button class="mt-6">Go home</Button></a>
		</div>
	</main>
{:else}
	<main class="min-h-screen bg-slate-50 px-4 py-10">
		<div class="mx-auto max-w-6xl">
			<section class="mb-8 rounded-3xl border border-emerald-100 bg-white p-8 text-slate-900 shadow-sm">
				<div class="flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
					<div>
						<p class="text-sm font-medium uppercase tracking-wide text-emerald-700">Public Knowledge Profile</p>
						<h1 class="mt-3 text-4xl font-bold">{data.profile.display_name || data.profile.username || username}</h1>
						{#if data.profile.bio}<p class="mt-4 max-w-2xl text-slate-600">{data.profile.bio}</p>{/if}
					</div>
					<div class="rounded-2xl border border-emerald-100 bg-emerald-50 px-4 py-3 text-sm text-emerald-800">@{data.profile.username || username}</div>
				</div>
			</section>

			<div class="mb-8 grid grid-cols-2 gap-4 md:grid-cols-4">
				<Stat icon={Library} label="Library" value={data.library.length} />
				<Stat icon={StickyNote} label="Notes" value={data.notes.length} />
				<Stat icon={Star} label="Reviews" value={data.reviews.length} />
				<Stat icon={Layers3} label="Collections" value={data.collections.length} />
			</div>

			<div class="grid gap-6 lg:grid-cols-2">
				{@render PublicSection('Public Library', 'Sources this reader has chosen to share.', 'No public library items yet.', data.library.map(item => ({ id: item.id, title: item.source?.title || item.source_id, meta: `${item.status.replace('_', ' ')} · ${item.visibility}` })))}
				{@render PublicSection('Public Notes', 'Notes and reflections shared by this reader.', 'No public notes yet.', data.notes.map(note => ({ id: note.id, title: note.content, meta: note.content_type })))}
				{@render PublicSection('Public Reviews', 'Ratings and reviews shared publicly.', 'No public reviews yet.', data.reviews.map(review => ({ id: review.id, title: review.content || `${review.rating}/5 stars`, meta: `${review.rating}/5 stars` })))}
				{@render PublicSection('Public Collections', 'Curated lists from this profile.', 'No public collections yet.', data.collections.map(c => ({ id: c.id, title: c.name, meta: `${c.source_ids?.length || 0} sources` })))}
			</div>
		</div>
	</main>
{/if}

{#snippet PublicSection(title: string, description: string, emptyMsg: string, items: Array<{ id: string; title: string; meta: string }>)}
	<Card>
		<CardHeader>
			<CardTitle class="flex items-center gap-2"><BookOpen class="h-5 w-5 text-emerald-600" /> {title}</CardTitle>
			<CardDescription>{description}</CardDescription>
		</CardHeader>
		<CardContent>
			{#if items.length === 0}
				<div class="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">{emptyMsg}</div>
			{:else}
				<div class="space-y-3">
					{#each items.slice(0, 8) as item (item.id)}
						<div class="rounded-lg border border-slate-200 bg-white p-4">
							<p class="font-medium text-slate-900">{item.title}</p>
							<p class="mt-1 text-sm text-slate-500">{item.meta}</p>
						</div>
					{/each}
				</div>
			{/if}
		</CardContent>
	</Card>
{/snippet}
