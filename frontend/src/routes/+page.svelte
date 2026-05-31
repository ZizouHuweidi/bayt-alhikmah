<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { Button } from '$lib/components/ui';
	import { Book, BookOpen, Brain, Library, Quote, Scroll, Sparkles, Star, Users } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Bayt al Hikmah - House of Wisdom</title>
</svelte:head>

<main class="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100">
	<section class="relative overflow-hidden px-6 py-24">
		<div class="absolute inset-0 bg-gradient-to-br from-emerald-50 via-teal-50 to-cyan-50 opacity-70"></div>
		<div class="absolute -left-20 -top-24 h-96 w-96 rounded-full bg-emerald-200 opacity-30 mix-blend-multiply blur-3xl animate-pulse"></div>
		<div class="absolute -bottom-20 -right-20 h-96 w-96 rounded-full bg-teal-200 opacity-30 mix-blend-multiply blur-3xl animate-pulse"></div>

		<div class="relative mx-auto max-w-6xl text-center">
			<div class="mb-8 inline-flex items-center gap-2 rounded-full bg-emerald-100 px-4 py-2 text-sm font-medium text-emerald-800">
				<Sparkles class="h-4 w-4" />
				<span>Inspired by the House of Wisdom</span>
			</div>

			<h1 class="mb-6 text-5xl font-bold tracking-tight text-slate-900 md:text-7xl">
				<span class="block">Bayt al</span>
				<span class="block bg-gradient-to-r from-emerald-600 to-teal-600 bg-clip-text text-transparent">Hikmah</span>
			</h1>

			<p class="mx-auto mb-4 max-w-2xl text-xl text-slate-600">
				A modern, AI-powered platform for organizing, engaging with, and tracking knowledge sources across all media.
			</p>

			<p class="mx-auto mb-10 max-w-xl text-slate-500">
				Inspired by the historical center of learning from the Islamic Golden Age, reimagined for the 21st century.
			</p>

			<div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
				{#if auth.isAuthenticated}
					<a href="/dashboard">
						<Button size="lg" class="rounded-full bg-emerald-600 px-8 py-6 text-lg font-semibold text-white shadow-lg transition-all hover:bg-emerald-700 hover:shadow-xl">
							Go to Dashboard
						</Button>
					</a>
				{:else}
					<a href="/registration">
						<Button size="lg" class="rounded-full bg-emerald-600 px-8 py-6 text-lg font-semibold text-white shadow-lg transition-all hover:bg-emerald-700 hover:shadow-xl">
							Start Your Journey
						</Button>
					</a>
					<a href="/login">
						<Button size="lg" variant="outline" class="rounded-full border-2 border-slate-300 px-8 py-6 text-lg font-semibold hover:border-emerald-500 hover:text-emerald-700">
							Sign In
						</Button>
					</a>
				{/if}
			</div>

			{#if auth.isAuthenticated && auth.user.email}
				<p class="mt-4 font-medium text-emerald-600">Welcome back, {auth.user.firstName || auth.user.email}!</p>
			{/if}
		</div>
	</section>

	<section class="border-y border-slate-200 bg-white/50 px-6 py-12">
		<div class="mx-auto max-w-6xl">
			<p class="mb-8 text-center text-sm font-semibold uppercase tracking-wider text-slate-500">Track Everything You Learn From</p>
			<div class="flex flex-wrap items-center justify-center gap-8">
				<div class="flex items-center gap-3 rounded-full bg-white px-6 py-3 shadow-sm"><Book class="h-6 w-6 text-emerald-600" /><span class="font-medium text-slate-700">Books</span></div>
				<div class="flex items-center gap-3 rounded-full bg-white px-6 py-3 shadow-sm"><Scroll class="h-6 w-6 text-emerald-600" /><span class="font-medium text-slate-700">Academic Papers</span></div>
				<div class="flex items-center gap-3 rounded-full bg-white px-6 py-3 shadow-sm"><Quote class="h-6 w-6 text-emerald-600" /><span class="font-medium text-slate-700">Articles</span></div>
				<div class="flex items-center gap-3 rounded-full bg-white px-6 py-3 shadow-sm"><Sparkles class="h-6 w-6 text-emerald-600" /><span class="font-medium text-slate-700">Podcasts & Videos</span></div>
			</div>
		</div>
	</section>

	<section class="px-6 py-24">
		<div class="mx-auto max-w-6xl">
			<div class="mb-16 text-center">
				<h2 class="mb-4 text-3xl font-bold text-slate-900">Your Knowledge, Organized</h2>
				<p class="mx-auto max-w-2xl text-lg text-slate-600">Build a personal library that grows with you. Track what you read, watch, and listen to — and connect it all together.</p>
			</div>
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
				<div class="rounded-2xl border border-slate-200 bg-white p-8 transition-all duration-300 hover:border-emerald-200 hover:shadow-lg">
					<div class="mb-6 w-fit rounded-xl bg-slate-50 p-3"><Library class="h-8 w-8 text-emerald-600" /></div>
					<h3 class="mb-3 text-xl font-semibold text-slate-900">Organize Your Knowledge</h3>
					<p class="leading-relaxed text-slate-600">Track books, papers, podcasts, videos, and articles in one unified library. Build your personal collection with rich metadata.</p>
				</div>
				<div class="rounded-2xl border border-slate-200 bg-white p-8 transition-all duration-300 hover:border-emerald-200 hover:shadow-lg">
					<div class="mb-6 w-fit rounded-xl bg-slate-50 p-3"><BookOpen class="h-8 w-8 text-emerald-600" /></div>
					<h3 class="mb-3 text-xl font-semibold text-slate-900">Smart Annotations</h3>
					<p class="leading-relaxed text-slate-600">Create structured notes and annotations tied to specific pages, timestamps, or sections. Link your thoughts to source materials.</p>
				</div>
				<div class="rounded-2xl border border-slate-200 bg-white p-8 transition-all duration-300 hover:border-emerald-200 hover:shadow-lg">
					<div class="mb-6 w-fit rounded-xl bg-slate-50 p-3"><Brain class="h-8 w-8 text-emerald-600" /></div>
					<h3 class="mb-3 text-xl font-semibold text-slate-900">Connect Ideas</h3>
					<p class="leading-relaxed text-slate-600">Use tags, topics, and taxonomies to create meaningful connections between different sources and your personal insights.</p>
				</div>
				<div class="rounded-2xl border border-slate-200 bg-white p-8 transition-all duration-300 hover:border-emerald-200 hover:shadow-lg">
					<div class="mb-6 w-fit rounded-xl bg-slate-50 p-3"><Users class="h-8 w-8 text-emerald-600" /></div>
					<h3 class="mb-3 text-xl font-semibold text-slate-900">Build Your Identity</h3>
					<p class="leading-relaxed text-slate-600">Create a public or private knowledge profile. Share reviews, ratings, and curated lists with the community.</p>
				</div>
			</div>
		</div>
	</section>

	<section class="px-6 py-24">
		<div class="mx-auto max-w-4xl text-center">
			<div class="rounded-3xl border border-emerald-100 bg-white p-12 shadow-sm">
				<h2 class="mb-4 text-3xl font-bold text-slate-900">Ready to Build Your Library?</h2>
				<p class="mb-8 text-lg text-slate-600">Join Bayt al Hikmah and start tracking your knowledge journey today.</p>
				{#if auth.isAuthenticated}
					<a href="/dashboard"><Button size="lg" variant="secondary" class="rounded-full bg-white px-8 py-6 text-lg font-semibold text-emerald-700 hover:bg-emerald-50">Go to Dashboard</Button></a>
				{:else}
					<a href="/registration"><Button size="lg" variant="secondary" class="rounded-full bg-white px-8 py-6 text-lg font-semibold text-emerald-700 hover:bg-emerald-50">Create Free Account</Button></a>
				{/if}
			</div>
		</div>
	</section>
</main>
