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
		<button
			class="w-full text-left p-2 rounded-md border border-gray-200 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/30 transition-all group"
			on:click={() => (expanded = true)}
		>
			<div class="flex items-center justify-between gap-2">
				<div class="flex flex-col gap-1 text-xs flex-1 min-w-0">
					<div class="flex items-center gap-2">
						<span class="text-gray-500 dark:text-gray-400 font-medium">Fourn:</span>
						<span class="font-semibold dark:text-gray-200 truncate">{getFournisseurName(item.fournisseur)}</span>
					</div>
					{#if item.amount_per_bundle || item.ref_bundle}
						<div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
							{#if item.amount_per_bundle}
								<span>Bundle: {item.amount_per_bundle}</span>
							{/if}
							{#if item.ref_bundle}
								<span>Ref: {item.ref_bundle}</span>
							{/if}
						</div>
					{/if}
				</div>
				<iconify-icon icon="mdi:chevron-down" width="16" height="16" class="text-gray-400 group-hover:text-blue-500 transition-colors flex-shrink-0"></iconify-icon>
			</div>
		</button>
	{:else}
		<div class="border border-blue-500 dark:border-blue-400 rounded-md p-3 bg-blue-50/50 dark:bg-blue-950/20">
			<div class="flex items-center justify-between mb-3">
				<span class="text-xs font-semibold text-blue-700 dark:text-blue-300">Modifier Fournisseur</span>
				<button
					class="px-2 py-1 text-xs bg-white dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-300 dark:border-gray-500 rounded flex items-center gap-1 transition-colors"
					on:click={() => (expanded = false)}
				>
					<iconify-icon icon="mdi:check" width="14" height="14"></iconify-icon>
					<span>Fermer</span>
				</button>
			</div>
			<div class="space-y-2">
				<div>
					<label class="text-xs block font-medium dark:text-gray-300 mb-1">Fournisseur</label>
					<select
						class="w-full text-xs px-2 py-1.5 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
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
					<label class="text-xs block font-medium dark:text-gray-300 mb-1">Bundle/Quantité</label>
					<input
						type="number"
						class="w-full text-xs px-2 py-1.5 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
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
					<label class="text-xs block font-medium dark:text-gray-300 mb-1">Code Référence</label>
					<input
						type="text"
						class="w-full text-xs px-2 py-1.5 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
						value={item.ref_bundle ?? ''}
						placeholder="-"
						on:input={(e) => {
							// @ts-ignore
							editItem(item.id, { ref_bundle: e.target?.value || undefined }, item.category_id);
						}}
					/>
				</div>
			</div>
		</div>
	{/if}
</div>