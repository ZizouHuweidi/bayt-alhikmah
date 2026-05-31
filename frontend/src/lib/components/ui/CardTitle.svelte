<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children,
		...rest
	} = $props();
</script>

<div data-slot="card-title" class={cn('leading-none font-semibold', className)} {...rest}>
	{@render children?.()}
</div>
