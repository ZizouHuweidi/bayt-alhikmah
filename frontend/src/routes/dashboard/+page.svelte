<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Stat } from '$lib/components/ui';
	import { Book, CheckCircle2, Layers3, Library, Loader2, LogOut, Plus, Star, StickyNote } from 'lucide-svelte';

	type DashboardData = {
		sources: import('$lib/api').Source[];
		library: import('$lib/api').LibraryItemWithSource[];
		notes: import('$lib/api').Note[];
		reviews: import('$lib/api').Review[];
		collections: import('$lib/api').Collection[];
	};

	let data = $state<DashboardData>({ sources: [], library: [], notes: [], reviews: [], collections: [] });
	let loadingData = $state(true);
	let error = $state<string | null>(null);
	let bookTitle = $state('');
	let bookAuthor = $state('');
	let bookISBN = $state('');
	let savingBook = $state(false);
	let selectedSourceID = $state('');
	let noteContent = $state('');
	let notePublic = $state(false);
	let reviewRating = $state('5');
	let reviewContent = $state('');
	let collectionName = $state('');
	let collectionPublic = $state(false);
	let savingActivity = $state(false);

	onMount(() => {
		if (!auth.isAuthenticated) {
			goto('/login');
		}
	});

	async function loadDashboard() {
		if (!auth.accessToken) return;
		loadingData = true;
		error = null;
		try {
			const [
				{ listSources },
				{ listLibrary },
				{ listNotes },
				{ listReviews },
				{ listCollections },
			] = await Promise.all([
				import('$lib/api'),
				import('$lib/api'),
				import('$lib/api'),
				import('$lib/api'),
				import('$lib/api'),
			]);
			const [sources, library, notes, reviews, collections] = await Promise.all([
				listSources(),
				listLibrary(auth.accessToken!),
				listNotes(auth.accessToken!),
				listReviews(auth.accessToken!),
				listCollections(auth.accessToken!),
			]);
			data = {
				sources: Array.isArray(sources) ? sources : [],
				library: Array.isArray(library) ? library : [],
				notes: Array.isArray(notes) ? notes : [],
				reviews: Array.isArray(reviews) ? reviews : [],
				collections: Array.isArray(collections) ? collections : [],
			};
			if (!selectedSourceID && data.library.length > 0) {
				selectedSourceID = data.library[0].source_id;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load dashboard';
		} finally {
			loadingData = false;
		}
	}

	$effect(() => {
		if (auth.isAuthenticated && auth.accessToken) {
			loadDashboard();
		}
	});

	async function handleLogout() {
		await auth.logout();
		goto('/');
	}

	async function handleCreateBook(event: Event) {
		event.preventDefault();
		if (!auth.accessToken || !bookTitle.trim()) return;
		savingBook = true;
		error = null;
		try {
			const { createBook, addLibraryItem } = await import('$lib/api');
			const book = await createBook(auth.accessToken, {
				title: bookTitle.trim(),
				isbn_13: bookISBN.trim() || undefined,
				contributors: bookAuthor.trim() ? [{ name: bookAuthor.trim(), role: 'author' }] : [],
			});
			await addLibraryItem(auth.accessToken, book.source.id);
			bookTitle = '';
			bookAuthor = '';
			bookISBN = '';
			await loadDashboard();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create book';
		} finally {
			savingBook = false;
		}
	}

	async function handleCreateNote(event: Event) {
		event.preventDefault();
		if (!auth.accessToken || !selectedSourceID || !noteContent.trim()) return;
		savingActivity = true;
		error = null;
		try {
			const { createNote } = await import('$lib/api');
			await createNote(auth.accessToken, {
				source_id: selectedSourceID,
				content: noteContent.trim(),
				content_type: 'note',
				is_public: notePublic,
			});
			noteContent = '';
			notePublic = false;
			await loadDashboard();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create note';
		} finally {
			savingActivity = false;
		}
	}

	async function handleCreateReview(event: Event) {
		event.preventDefault();
		if (!auth.accessToken || !selectedSourceID) return;
		savingActivity = true;
		error = null;
		try {
			const { createReview } = await import('$lib/api');
			await createReview(auth.accessToken, {
				source_id: selectedSourceID,
				rating: Number(reviewRating),
				content: reviewContent.trim() || undefined,
				is_public: true,
			});
			reviewRating = '5';
			reviewContent = '';
			await loadDashboard();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create review';
		} finally {
			savingActivity = false;
		}
	}

	async function handleCreateCollection(event: Event) {
		event.preventDefault();
		if (!auth.accessToken || !selectedSourceID || !collectionName.trim()) return;
		savingActivity = true;
		error = null;
		try {
			const { createCollection } = await import('$lib/api');
			await createCollection(auth.accessToken, {
				name: collectionName.trim(),
				is_public: collectionPublic,
				source_ids: [selectedSourceID],
			});
			collectionName = '';
			collectionPublic = false;
			await loadDashboard();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create collection';
		} finally {
			savingActivity = false;
		}
	}

	async function runAction(action: () => Promise<unknown>, fallback: string) {
		error = null;
		try {
			await action();
			await loadDashboard();
		} catch (err) {
			error = err instanceof Error ? err.message : fallback;
		}
	}

	const librarySourceIDs = $derived(new Set(data.library.map(item => item.source_id)));
	const librarySources = $derived(data.library.map(item => ({ item, source: item.source })));
</script>

<svelte:head>
	<title>Dashboard - Bayt al Hikmah</title>
</svelte:head>

{#if auth.isLoading}
	<div class="flex min-h-screen items-center justify-center bg-slate-50">
		<Loader2 class="h-8 w-8 animate-spin text-emerald-600" />
	</div>
{/if}

{#if auth.isAuthenticated}
	<div class="min-h-screen bg-slate-50">
		<header class="border-b border-slate-200 bg-white">
			<div class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
				<div class="flex items-center gap-3">
					<div class="rounded-lg bg-emerald-600 p-2"><Library class="h-5 w-5 text-white" /></div>
					<span class="text-xl font-bold text-slate-900">Bayt al Hikmah</span>
				</div>
				<div class="flex items-center gap-4">
					<span class="hidden text-sm text-slate-600 sm:inline">Welcome, {auth.user.username || auth.user.email}</span>
					<Button variant="ghost" size="sm" onclick={handleLogout} class="text-slate-600 hover:text-red-600"><LogOut class="mr-2 h-4 w-4" /> Logout</Button>
				</div>
			</div>
		</header>

		<main class="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
			<div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
				<div>
					<h1 class="text-3xl font-bold text-slate-900">Dashboard</h1>
					<p class="mt-2 text-slate-600">Build and test your knowledge library with real platform data.</p>
				</div>
				<a href="/settings"><Button variant="outline">Profile settings</Button></a>
			</div>

			<div class="mb-8 grid grid-cols-2 gap-4 lg:grid-cols-5">
				<Stat icon={Book} label="Book sources" value={data.sources.length} />
				<Stat icon={Library} label="Library" value={data.library.length} />
				<Stat icon={StickyNote} label="Notes" value={data.notes.length} />
				<Stat icon={Star} label="Reviews" value={data.reviews.length} />
				<Stat icon={Layers3} label="Collections" value={data.collections.length} />
			</div>

			{#if error}
				<div class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
			{/if}

			<div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_380px]">
				<div class="space-y-6">
					<Card>
						<CardHeader>
							<CardTitle class="flex items-center gap-2"><Library class="h-5 w-5 text-emerald-600" /> My Library</CardTitle>
							<CardDescription>Your tracked sources. New books created here are added to your library automatically.</CardDescription>
						</CardHeader>
						<CardContent>
							{#if loadingData}
								<div class="flex items-center gap-2 py-8 text-sm text-slate-500"><Loader2 class="h-4 w-4 animate-spin" /> Loading platform data...</div>
							{:else if data.library.length === 0}
								<div class="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">No library items yet. Create a book to start testing the platform.</div>
							{:else}
								<div class="space-y-3">
									{#each data.library.slice(0, 8) as item (item.id)}
										<div class="rounded-lg border border-slate-200 bg-white p-4">
											<div class="flex items-start justify-between gap-3">
												<div>
													<a href="/sources/books/{item.source_id}" class="font-medium text-slate-900 hover:text-emerald-700">{item.source?.title || item.source_id}</a>
													<p class="mt-1 text-sm text-slate-500">{item.status.replace('_', ' ')} · {item.visibility}</p>
												</div>
												<div class="flex flex-col items-end gap-2">
													<CheckCircle2 class="h-5 w-5 text-emerald-600" />
													<select class="h-8 rounded-md border border-slate-300 bg-white px-2 text-xs" value={item.status}
														onchange={async (e) => {
															if (!auth.accessToken) return;
															const { updateLibraryItem } = await import('$lib/api');
															runAction(() => updateLibraryItem(auth.accessToken!, item.id, { status: (e.target as HTMLSelectElement).value }).then(() => {}), 'Failed to update status');
														}}
													>
														<option value="to_consume">To consume</option>
														<option value="in_progress">In progress</option>
														<option value="completed">Completed</option>
														<option value="paused">Paused</option>
														<option value="abandoned">Abandoned</option>
													</select>
												</div>
											</div>
										</div>
									{/each}
								</div>
							{/if}
						</CardContent>
					</Card>

					<Card>
						<CardHeader>
							<CardTitle>Available Book Sources</CardTitle>
							<CardDescription>Public book records currently in the platform.</CardDescription>
						</CardHeader>
						<CardContent>
							{#if loadingData}
								<div class="flex items-center gap-2 py-8 text-sm text-slate-500"><Loader2 class="h-4 w-4 animate-spin" /> Loading platform data...</div>
							{:else if data.sources.length === 0}
								<div class="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">No book sources yet.</div>
							{:else}
								<div class="space-y-3">
									{#each data.sources.slice(0, 8) as source (source.id)}
										<div class="flex flex-col gap-3 rounded-lg border border-slate-200 bg-white p-4 sm:flex-row sm:items-center sm:justify-between">
											<div><p class="font-medium text-slate-900">{source.title}</p><p class="mt-1 text-sm text-slate-500">{source.publisher || source.type}</p></div>
											{#if !librarySourceIDs.has(source.id)}
												<Button variant="outline" size="sm" onclick={async () => {
													if (!auth.accessToken) return;
													try {
														const { addLibraryItem } = await import('$lib/api');
														await addLibraryItem(auth.accessToken!, source.id);
														await loadDashboard();
													} catch (err) {
														error = err instanceof Error ? err.message : 'Failed to add source';
													}
												}}>
													Add to library
												</Button>
											{/if}
										</div>
									{/each}
								</div>
							{/if}
						</CardContent>
					</Card>
				</div>

				<aside class="space-y-6">
					<Card>
						<CardHeader>
							<CardTitle class="flex items-center gap-2"><Plus class="h-5 w-5 text-emerald-600" /> Add a Book</CardTitle>
							<CardDescription>Minimal book creation for MVP testing.</CardDescription>
						</CardHeader>
						<CardContent>
							<form class="space-y-4" onsubmit={handleCreateBook}>
								<Input bind:value={bookTitle} placeholder="Book title" required />
								<Input bind:value={bookAuthor} placeholder="Author" />
								<Input bind:value={bookISBN} placeholder="ISBN-13" />
								<Button type="submit" class="w-full" disabled={savingBook || !bookTitle.trim()}>{savingBook ? 'Saving...' : 'Create and add'}</Button>
							</form>
						</CardContent>
					</Card>

					<Card>
						<CardHeader>
							<CardTitle>Work With a Source</CardTitle>
							<CardDescription>Select a library item, then add notes, reviews, and collections.</CardDescription>
						</CardHeader>
						<CardContent class="space-y-4">
							<select class="h-9 w-full rounded-md border border-slate-300 bg-white px-3 text-sm" bind:value={selectedSourceID} disabled={librarySources.length === 0}>
								{#if librarySources.length === 0}
									<option value="">No library sources yet</option>
								{:else}
									{#each librarySources as { source } (source.id)}
										<option value={source.id}>{source?.title || source.id}</option>
									{/each}
								{/if}
							</select>

							<form class="space-y-3" onsubmit={handleCreateNote}>
								<textarea class="min-h-24 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm" bind:value={noteContent} placeholder="Write a note..."></textarea>
								<label class="flex items-center gap-2 text-sm text-slate-600">
									<input type="checkbox" bind:checked={notePublic} /> Public note
								</label>
								<Button type="submit" variant="outline" class="w-full" disabled={savingActivity || !selectedSourceID || !noteContent.trim()}>Add note</Button>
							</form>

							<form class="space-y-3 border-t border-slate-200 pt-4" onsubmit={handleCreateReview}>
								<div class="grid grid-cols-[90px_1fr] gap-3">
									<select class="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm" bind:value={reviewRating}>
										{#each [5, 4, 3, 2, 1] as rating}
											<option value={rating}>{rating} star</option>
										{/each}
									</select>
									<Input bind:value={reviewContent} placeholder="Short review" />
								</div>
								<Button type="submit" variant="outline" class="w-full" disabled={savingActivity || !selectedSourceID}>Add public review</Button>
							</form>

							<form class="space-y-3 border-t border-slate-200 pt-4" onsubmit={handleCreateCollection}>
								<Input bind:value={collectionName} placeholder="Collection name" />
								<label class="flex items-center gap-2 text-sm text-slate-600">
									<input type="checkbox" bind:checked={collectionPublic} /> Public collection
								</label>
								<Button type="submit" variant="outline" class="w-full" disabled={savingActivity || !selectedSourceID || !collectionName.trim()}>Create collection</Button>
							</form>
						</CardContent>
					</Card>

					<Card>
						<CardHeader><CardTitle>Recent Notes</CardTitle></CardHeader>
						<CardContent>
							{#if data.notes.length === 0}
								<div class="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">No notes yet.</div>
							{:else}
								<div class="space-y-3">
									{#each data.notes.slice(0, 5) as note (note.id)}
										<div class="rounded-lg border border-slate-200 bg-white p-3 text-sm text-slate-700">
											<p>{note.content}</p>
											<Button variant="ghost" size="xs" class="mt-2 text-red-600" onclick={async () => {
												if (!auth.accessToken) return;
												const { deleteNote } = await import('$lib/api');
												runAction(() => deleteNote(auth.accessToken!, note.id).then(() => {}), 'Failed to delete note');
											}}>Delete</Button>
										</div>
									{/each}
								</div>
							{/if}
						</CardContent>
					</Card>

					<Card>
						<CardHeader><CardTitle>Recent Reviews</CardTitle></CardHeader>
						<CardContent>
							{#if data.reviews.length === 0}
								<div class="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">No reviews yet.</div>
							{:else}
								<div class="space-y-3">
									{#each data.reviews.slice(0, 5) as review (review.id)}
										<div class="rounded-lg border border-slate-200 bg-white p-3 text-sm text-slate-700">
											<p class="font-medium">{review.rating}/5 stars</p>
											{#if review.content}<p class="mt-1">{review.content}</p>{/if}
											<Button variant="ghost" size="xs" class="mt-2 text-red-600" onclick={async () => {
												if (!auth.accessToken) return;
												const { deleteReview } = await import('$lib/api');
												runAction(() => deleteReview(auth.accessToken!, review.id).then(() => {}), 'Failed to delete review');
											}}>Delete</Button>
										</div>
									{/each}
								</div>
							{/if}
						</CardContent>
					</Card>

					<Card>
						<CardHeader><CardTitle>Collections</CardTitle></CardHeader>
						<CardContent>
							{#if data.collections.length === 0}
								<div class="rounded-lg border border-dashed border-slate-300 bg-slate-50 p-6 text-center text-sm text-slate-500">No collections yet.</div>
							{:else}
								<div class="space-y-3">
									{#each data.collections.slice(0, 5) as collection (collection.id)}
										<div class="rounded-lg border border-slate-200 bg-white p-3 text-sm text-slate-700">
											<p class="font-medium">{collection.name}</p>
											<p class="mt-1 text-slate-500">{collection.is_public ? 'Public' : 'Private'} · {collection.source_ids?.length || 0} sources</p>
											<Button variant="ghost" size="xs" class="mt-2 text-red-600" onclick={async () => {
												if (!auth.accessToken) return;
												const { deleteCollection } = await import('$lib/api');
												runAction(() => deleteCollection(auth.accessToken!, collection.id).then(() => {}), 'Failed to delete collection');
											}}>Delete</Button>
										</div>
									{/each}
								</div>
							{/if}
						</CardContent>
					</Card>
				</aside>
			</div>
		</main>
	</div>
{/if}
