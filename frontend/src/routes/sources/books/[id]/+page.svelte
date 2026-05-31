<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { BookOpen, Library, Loader2, UserRound } from 'lucide-svelte';
	import { Button, Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui';
	import type { Book } from '$lib/api';

	let { id } = $props();

	let book = $state<Book | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let added = $state(false);

	$effect(() => {
		let cancelled = false;
		async function loadBook() {
			loading = true;
			error = null;
			try {
				const { getBook } = await import('$lib/api');
				const result = await getBook(id);
				if (!cancelled) book = result;
			} catch (err) {
				if (!cancelled) error = err instanceof Error ? err.message : 'Failed to load book';
			} finally {
				if (!cancelled) loading = false;
			}
		}
		loadBook();
		return () => { cancelled = true; };
	});

	async function handleAdd() {
		if (!auth.accessToken) return;
		error = null;
		try {
			const { addLibraryItem } = await import('$lib/api');
			await addLibraryItem(auth.accessToken, id);
			added = true;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to add to library';
		}
	}
</script>

<svelte:head>
	<title>{book?.source.title || 'Book'} - Bayt al Hikmah</title>
</svelte:head>

{#if loading}
	<main class="flex min-h-screen items-center justify-center bg-slate-50">
		<Loader2 class="h-8 w-8 animate-spin text-emerald-600" />
	</main>
{:else if !book}
	<main class="min-h-screen bg-slate-50 px-4 py-16">
		<div class="mx-auto max-w-2xl rounded-2xl border bg-white p-8 text-center">
			<h1 class="text-2xl font-bold text-slate-900">Book not found</h1>
			<a href="/dashboard"><Button class="mt-6">Back to dashboard</Button></a>
		</div>
	</main>
{:else}
	<main class="min-h-screen bg-slate-50 px-4 py-10">
		<div class="mx-auto max-w-5xl">
			<a href="/dashboard" class="text-sm font-medium text-emerald-700 hover:underline">Back to dashboard</a>

			{#if error}
				<div class="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div>
			{/if}

			<section class="mt-6 rounded-3xl bg-white p-8 shadow-sm">
				<div class="flex flex-col gap-6 md:flex-row md:items-start md:justify-between">
					<div>
						<div class="mb-4 inline-flex items-center gap-2 rounded-full bg-emerald-50 px-3 py-1 text-sm font-medium text-emerald-700">
							<BookOpen class="h-4 w-4" /> Book
						</div>
						<h1 class="text-4xl font-bold text-slate-900">{book.source.title}</h1>
						{#if book.source.subtitle}
							<p class="mt-3 text-xl text-slate-600">{book.source.subtitle}</p>
						{/if}
						{#if book.contributors && book.contributors.length > 0}
							<p class="mt-4 flex items-center gap-2 text-slate-600">
								<UserRound class="h-4 w-4" />
								{book.contributors.map(c => c.name).join(', ')}
							</p>
						{/if}
					</div>
					{#if auth.isAuthenticated}
						<Button onclick={handleAdd} disabled={added}>
							<Library class="h-4 w-4" /> {added ? 'Added' : 'Add to library'}
						</Button>
					{/if}
				</div>
			</section>

			<div class="mt-6 grid gap-6 md:grid-cols-2">
				<Card>
					<CardHeader><CardTitle>Metadata</CardTitle></CardHeader>
					<CardContent class="space-y-3 text-sm text-slate-700">
						<div class="flex justify-between gap-4 border-b border-slate-100 pb-2">
							<span class="text-slate-500">ISBN-13</span>
							<span class="font-medium text-slate-900">{book.metadata?.isbn_13 || book.source.isbn || 'Not set'}</span>
						</div>
						<div class="flex justify-between gap-4 border-b border-slate-100 pb-2">
							<span class="text-slate-500">Publisher</span>
							<span class="font-medium text-slate-900">{book.metadata?.publisher || book.source.publisher || 'Not set'}</span>
						</div>
						<div class="flex justify-between gap-4 border-b border-slate-100 pb-2">
							<span class="text-slate-500">Language</span>
							<span class="font-medium text-slate-900">{book.metadata?.language || 'Not set'}</span>
						</div>
						<div class="flex justify-between gap-4 border-b border-slate-100 pb-2">
							<span class="text-slate-500">Pages</span>
							<span class="font-medium text-slate-900">{book.metadata?.page_count?.toString() || 'Not set'}</span>
						</div>
					</CardContent>
				</Card>
				<Card>
					<CardHeader><CardTitle>Description</CardTitle></CardHeader>
					<CardContent>
						<p class="text-sm leading-6 text-slate-700">{book.source.description || 'No description yet.'}</p>
					</CardContent>
				</Card>
			</div>
		</div>
	</main>
{/if}
