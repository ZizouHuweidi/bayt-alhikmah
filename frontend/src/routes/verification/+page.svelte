<script lang="ts">
	import { CheckCircle, Library } from 'lucide-svelte';
	import { Button, Card, CardContent } from '$lib/components/ui';
</script>

<svelte:head>
	<title>Email Verification - Bayt al Hikmah</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 px-4">
	<div class="w-full max-w-md">
		<div class="mb-8 flex flex-col items-center">
			<div class="mb-4 rounded-xl bg-emerald-600 p-3"><Library class="h-8 w-8 text-white" /></div>
			<h1 class="text-3xl font-bold text-slate-900">Bayt al Hikmah</h1>
		</div>

		<Card>
			<CardContent class="pt-6">
				<div class="flex flex-col items-center text-center">
					<CheckCircle class="mb-4 h-16 w-16 text-emerald-500" />
					<h2 class="mb-2 text-xl font-semibold text-slate-900">Email Verification</h2>
					<p class="mb-6 text-slate-600">Email verification is not implemented yet.</p>
					<a href="/dashboard"><Button>Go to Dashboard</Button></a>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
