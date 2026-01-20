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
		<div
			class="flex items-center gap-1 flex-wrap text-xs"
			role="button"
			tabindex="0"
			on:click={() => (expanded = true)}
			on:keydown={(e) => e.key === 'Enter' && (expanded = true)}
		>
			{#each priceRoles as role}
				<span class="bg-gray-100 dark:bg-gray-700 px-1 rounded">{getCompactPrice(role)}</span>
			{/each}
			<button
				class="ml-1 p-1 hover:bg-gray-200 dark:hover:bg-gray-600 rounded transition-colors"
				on:click|stopPropagation={() => (expanded = true)}
			>
				<iconify-icon icon="mdi:pencil" width="14" height="14"></iconify-icon>
			</button>
		</div>
	{:else}
		<div class="min-w-[180px] space-y-1 mt-8">
			<button
				class="absolute -top-8 right-2 px-2 py-1 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 rounded z-10 flex items-center gap-1 text-xs"
				on:click={() => (expanded = false)}
			>
				<iconify-icon icon="mdi:close" width="14" height="14"></iconify-icon>
				<span>Fermer</span>
			</button>
			{#each priceRoles as role}
				{@const price = item.prices[role]}
				{@const label = role === 'coutant' ? 'Coutant' : role === 'externe' ? 'Externe' : role === 'ceten' ? 'Ceten' : role === 'staff_bar' ? 'Staff' : role === 'privilegies' ? 'Privil.' : 'Menu'}
				<div class="flex items-center gap-2">
					<span class="text-xs w-12 shrink-0 dark:text-gray-300">{label}:</span>
					<input
						type="number"
						id="price-{item.id}-{role}"
						placeholder={formatPrice(price)}
						class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
						on:input={(e) => handlePriceChange(role, e)}
					/>
				</div>
			{/each}
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