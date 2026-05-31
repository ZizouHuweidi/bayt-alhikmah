<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		variant = 'default' as 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link',
		size = 'default' as 'default' | 'xs' | 'sm' | 'lg' | 'icon',
		children,
		...rest
	} = $props();

	const variants: Record<string, string> = {
		default: 'bg-emerald-600 text-white hover:bg-emerald-700',
		destructive: 'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-500',
		outline: 'border border-slate-300 bg-white hover:bg-slate-100 hover:text-slate-900',
		secondary: 'bg-slate-100 text-slate-900 hover:bg-slate-200',
		ghost: 'hover:bg-slate-100 hover:text-slate-900',
		link: 'text-emerald-600 underline-offset-4 hover:underline',
	};

	const sizes: Record<string, string> = {
		default: 'h-9 px-4 py-2',
		xs: 'h-7 rounded-md px-2 text-xs',
		sm: 'h-8 rounded-md px-3',
		lg: 'h-10 rounded-md px-6',
		icon: 'h-9 w-9',
	};
</script>

<button
	class={cn(
		'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		variants[variant],
		sizes[size],
		className
	)}
	{...rest}
>
	{@render children?.()}
</button>
