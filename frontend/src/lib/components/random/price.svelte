<script lang="ts">
	import { formatPrice } from '$lib/utils';

	export let amount: number = 0;
</script>

<div class={$$restProps.class || ''}>
	{formatPrice(amount)}
</div>
