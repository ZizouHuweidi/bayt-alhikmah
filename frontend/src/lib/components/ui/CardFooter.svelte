<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children,
		...rest
	} = $props();
</script>

<div data-slot="card-footer" class={cn('flex items-center px-6 [.border-t]:pt-6', className)} {...rest}>
	{@render children?.()}
</div>
