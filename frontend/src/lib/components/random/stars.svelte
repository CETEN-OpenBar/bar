<script lang="ts">
	import { formatPrice } from '$lib/utils';

	export let stars: number = 0;
	export let icon_size: number = 7;
</script>

<div class={$$restProps.class || ''}>
	<!-- filled yellow star svg -->
	<svg
		xmlns="http://www.w3.org/2000/svg"
		class="h-{icon_size} w-{icon_size} inline-block text-yellow-500"
		viewBox="0 0 20 20"
		fill="currentColor"
	>
		<path
			fill-rule="evenodd"
			d="M10 12.585l-4.243 2.415 1.02-4.744L2.93 6.456l4.752-.41L10 2l1.318 4.046 4.752.41-3.847 3.8 1.02 4.744z"
			clip-rule="evenodd"
		/>
	</svg>
	{stars} ({formatPrice(stars)})
</div>
