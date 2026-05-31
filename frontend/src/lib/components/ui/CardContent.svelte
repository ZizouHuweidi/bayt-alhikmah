<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children,
		...rest
	} = $props();
</script>

<div data-slot="card-content" class={cn('px-6', className)} {...rest}>
	{@render children?.()}
</div>
