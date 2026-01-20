<script lang="ts">
	import type { Fournisseur } from '$lib/api';
	import { fournisseurIterator } from '$lib/utils';

	export let item: {
		id: string;
		fournisseur?: Fournisseur;
		amount_per_bundle?: number;
		ref_bundle?: string;
		category_id: string;
	};
	export let editItem: (id: string, item: any, category_id: string) => void;

	let expanded = false;

	function getFournisseurName(fournisseur?: Fournisseur): string {
		if (!fournisseur) return '-';
		for (const [val, name] of fournisseurIterator) {
			if (val === fournisseur) return name;
		}
		return String(fournisseur);
	}
</script>

<div class="relative">
	{#if !expanded}
		<div
			class="flex items-center gap-2 text-xs"
			role="button"
			tabindex="0"
			on:click={() => (expanded = true)}
			on:keydown={(e) => e.key === 'Enter' && (expanded = true)}
		>
			<span class="dark:text-gray-300">{getFournisseurName(item.fournisseur)}</span>
			{#if item.amount_per_bundle || item.ref_bundle}
				<span class="text-gray-500 dark:text-gray-400">
					{item.amount_per_bundle ? `| B:${item.amount_per_bundle}` : ''}
					{item.ref_bundle ? ` | R:${item.ref_bundle}` : ''}
				</span>
			{/if}
			<button
				class="p-1 hover:bg-gray-200 dark:hover:bg-gray-600 rounded transition-colors"
				on:click|stopPropagation={() => (expanded = true)}
			>
				<iconify-icon icon="mdi:pencil" width="14" height="14"></iconify-icon>
			</button>
		</div>
	{:else}
		<div class="min-w-[160px] space-y-2 mt-8">
			<button
				class="absolute -top-8 right-2 px-2 py-1 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 rounded z-10 flex items-center gap-1 text-xs"
				on:click={() => (expanded = false)}
			>
				<iconify-icon icon="mdi:close" width="14" height="14"></iconify-icon>
				<span>Fermer</span>
			</button>
			<div>
				<label class="text-xs block dark:text-gray-400 mb-1">Fournisseur</label>
				<select
					class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
					value={item.fournisseur ?? ''}
					on:change={(e) => {
						// @ts-ignore
						editItem(item.id, { fournisseur: e.target?.value || undefined }, item.category_id);
					}}
				>
					<option value="">-</option>
					{#each fournisseurIterator as [val, name]}
						<option value={val}>{name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="text-xs block dark:text-gray-400 mb-1">Bundle/Qté</label>
				<input
					type="number"
					class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
					value={item.amount_per_bundle ?? ''}
					placeholder="0"
					on:input={(e) => {
						// @ts-ignore
						let val = e.target?.value;
						editItem(item.id, { amount_per_bundle: val ? parseInt(val) : undefined }, item.category_id);
					}}
				/>
			</div>
			<div>
				<label class="text-xs block dark:text-gray-400 mb-1">Ref Code</label>
				<input
					type="text"
					class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
					value={item.ref_bundle ?? ''}
					placeholder="-"
					on:input={(e) => {
						// @ts-ignore
						editItem(item.id, { ref_bundle: e.target?.value || undefined }, item.category_id);
					}}
				/>
			</div>
		</div>
	{/if}
</div>