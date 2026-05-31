<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children,
		...rest
	} = $props();
</script>

<div data-slot="alert-description" class={cn('text-muted-foreground col-start-2 grid justify-items-start gap-1 text-sm [&_p]:leading-relaxed', className)} {...rest}>
	{@render children?.()}
</div>
