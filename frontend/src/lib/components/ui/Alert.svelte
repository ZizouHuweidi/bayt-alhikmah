<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		variant = 'default' as 'default' | 'destructive',
		children,
		...rest
	} = $props();

	const variants: Record<string, string> = {
		default: 'bg-card text-card-foreground',
		destructive: 'text-destructive bg-card [&>svg]:text-current *:data-[slot=alert-description]:text-destructive/90',
	};
</script>

<div
	data-slot="alert"
	role="alert"
	class={cn(
		'relative w-full rounded-lg border px-4 py-3 text-sm grid has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] grid-cols-[0_1fr] has-[>svg]:gap-x-3 gap-y-0.5 items-start [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current',
		variants[variant],
		className
	)}
	{...rest}
>
	{@render children?.()}
</div>
