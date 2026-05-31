<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children,
		...rest
	} = $props();
</script>

<div data-slot="card-description" class={cn('text-muted-foreground text-sm', className)} {...rest}>
	{@render children?.()}
</div>
