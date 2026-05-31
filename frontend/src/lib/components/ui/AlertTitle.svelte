<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children,
		...rest
	} = $props();
</script>

<div data-slot="alert-title" class={cn('col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight', className)} {...rest}>
	{@render children?.()}
</div>
