<script lang="ts">
	import type { ItemPrices, AccountPriceRole } from '$lib/api';
	import { formatPrice, parsePrice } from '$lib/utils';

	export let item: { id: string; prices: ItemPrices; category_id: string };
	export let editItem: (id: string, item: any, category_id: string) => void;

	let expanded = false;
	let rebounceTimeout: number | null = null;

	const priceLabels: Record<string, string> = {
		'coutant': 'Cout',
		'externe': 'Ext',
		'ceten': 'Ceten',
		'staff_bar': 'Staff',
		'privilegies': 'Priv',
		'menu': 'Menu'
	};

	const priceRoles: AccountPriceRole[] = ['coutant', 'externe', 'ceten', 'staff_bar', 'privilegies', 'menu'];

	function getCompactPrice(role: AccountPriceRole): string {
		return `${priceLabels[role]}:${formatPrice(item.prices[role])}`;
	}

	function handlePriceChange(role: AccountPriceRole, e: Event) {
		let prices = item.prices;
		const target = e.target as HTMLInputElement;
		prices[role] = parsePrice(target.value);
		editItem(item.id, { prices: prices }, item.category_id);

		if (rebounceTimeout) clearTimeout(rebounceTimeout);
		rebounceTimeout = setTimeout(() => {
			let elems = document.querySelectorAll(`[id^="price-${item.id}-"]`);
			elems.forEach((elem) => {
				(elem as HTMLInputElement).value = '';
			});
		}, 1000);
	}
</script>

<div class="relative">
	{#if !expanded}
		<button
			class="w-full text-left p-2 rounded-md border border-gray-200 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/30 transition-all group"
			on:click={() => (expanded = true)}
		>
			<div class="flex items-center justify-between gap-2">
				<div class="grid grid-cols-2 gap-x-3 gap-y-1 text-xs flex-1 min-w-0">
					{#each priceRoles as role}
						{@const label = priceLabels[role]}
						{@const price = formatPrice(item.prices[role])}
						<div class="flex items-center gap-1.5">
							<span class="text-gray-500 dark:text-gray-400 font-medium">{label}:</span>
							<span class="font-semibold dark:text-gray-200">{price}</span>
						</div>
					{/each}
				</div>
				<iconify-icon icon="mdi:chevron-down" width="16" height="16" class="text-gray-400 group-hover:text-blue-500 transition-colors flex-shrink-0"></iconify-icon>
			</div>
		</button>
	{:else}
		<div class="border border-blue-500 dark:border-blue-400 rounded-md p-3 bg-blue-50/50 dark:bg-blue-950/20">
			<div class="flex items-center justify-between mb-3">
				<span class="text-xs font-semibold text-blue-700 dark:text-blue-300">Modifier Prix</span>
				<button
					class="px-2 py-1 text-xs bg-white dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-300 dark:border-gray-500 rounded flex items-center gap-1 transition-colors"
					on:click={() => (expanded = false)}
				>
					<iconify-icon icon="mdi:check" width="14" height="14"></iconify-icon>
					<span>Fermer</span>
				</button>
			</div>
			<div class="space-y-2">
				{#each priceRoles as role}
					{@const price = item.prices[role]}
					{@const label = role === 'coutant' ? 'Coutant' : role === 'externe' ? 'Externe' : role === 'ceten' ? 'Ceten' : role === 'staff_bar' ? 'Staff' : role === 'privilegies' ? 'Privil.' : 'Menu'}
					<div class="flex items-center gap-2">
						<label class="text-xs w-16 shrink-0 font-medium dark:text-gray-300">{label}:</label>
						<div class="flex-1 relative">
							<input
								type="number"
								id="price-{item.id}-{role}"
								placeholder={formatPrice(price)}
								class="w-full text-xs px-2 py-1.5 pr-6 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
								on:input={(e) => handlePriceChange(role, e)}
							/>
							<span class="absolute right-2 top-1/2 -translate-y-1/2 text-xs text-gray-400">€</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	input[type='number'] {
		-moz-appearance: textfield;
	}
	input[type='number']::-webkit-outer-spin-button,
	input[type='number']::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}
</style>