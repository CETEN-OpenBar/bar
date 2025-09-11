<script lang="ts">
	import { onMount } from 'svelte';
	import type { RemoteRefill } from '$lib/api';
    import { RemoteRefillState } from '$lib/api';
	import { refillsApi } from '$lib/requests/requests';
	import 'iconify-icon';
	import Spinner from '$lib/components/spinner.svelte';
	import Success from '$lib/components/success.svelte';
	import Error from '$lib/components/error.svelte';
	import { goto } from '$app/navigation';

    const state_names: Record<RemoteRefillState, string> = {
        [RemoteRefillState.RemoteRefillStarted]: "En attente",
        [RemoteRefillState.RemoteRefillProcessed]: "Validé",
        [RemoteRefillState.RemoteRefillAbandoned]: "Abandonné",
    };

	let refills: RemoteRefill[] = [];
	let start_date = '';
	let end_date = '';
	let nameFilter = '';
    let state: RemoteRefillState | undefined = undefined;
	let page = 0;
	let maxPage = 0;
	let refills_per_page = 15;
	const nextPage = () => {
		if (page < maxPage) {
			page++;
			fetchRefills();
		}
	};
	const prevPage = () => {
		if (page > 0) {
			page--;
			fetchRefills();
		}
	};

	async function fetchRefills() {
		const response = await refillsApi().getRemoteRefills(
			page,
			refills_per_page,
			start_date,
			end_date,
            state,
			nameFilter.trim(),
			{ withCredentials: true }
		);
		const data = response.data;
		refills = Array.isArray(data.remote_refills) ? data.remote_refills : [];
		page = data.page;
		maxPage = data.max_page;
		refills_per_page = data.limit;
	}

    let loading = false;
    let error: string | undefined = undefined
    let success: boolean = false

    const verifyRefill = (id: string) => {
        loading = true;
        refillsApi().verifyRemoteRefill(id, {withCredentials: true})
        .then((resp) => {
            if (resp.status != 200) {
                throw resp;
            }
            success = true
            fetchRefills();
        })
        .catch((err) => {
            switch (err.status) {
                case 401:
                    goto("/auth");
                    break;
                case 403:
                    error = "Vous n'êtes pas autorisé à faire ceci";
                    break;
                case 404:
                    error = "Rechargement introuvable";
                    fetchRefills();
                    break;
                case 402:
                    error = "Paiement non validé par HelloAsso";
                    break;
                case 409:
                    error = "Ce rechargement a déjà été validé";
                    fetchRefills();
                    break;
                case 500:
                default:
                    error = "Une erreur est survenue, réessayez plus tard"
            }
        })
        .finally(() => {
            loading = false;
            setTimeout(() => {error = undefined; success = false}, 2000);
        })
    }

    let remote_refills_available: boolean = false;
    onMount(() => {
        fetchRefills();
        refillsApi().getRemoteRefillStatus({withCredentials: true})
        .then((resp) => {
            remote_refills_available = resp.status == 200;
        })
        .catch((_) => {
            remote_refills_available = false;
        });
    })

</script>

{#if success}
<Success message = "Rechargement validé !" close={() => {success = false}}/>
{/if}

{#if error}
<Error {error} close={() => {error = undefined}}/>
{/if}

{#if !remote_refills_available}
<div class="w-full">
    <div class="bg-red-500 p-3">
        <div class="font-bold text-center text-lg">Les rechargements en ligne sont actuellement indisponibles</div>
    </div>
</div>
{/if}


<div class="flex-grow grid grid-cols-1 grid-rows-[auto_1fr_auto] bg-gray-50 dark:bg-gray-900 min-w-full">
	<div class="m-3 p-2">
		<div class="flex flex-wrap items-center gap-6">
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Du:</span>
				<input
					type="date"
					bind:value={start_date}
					on:change={fetchRefills}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
				/>
			</div>
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Au:</span>
				<input
					type="date"
					bind:value={end_date}
					on:change={fetchRefills}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
				/>
			</div>
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Compte:</span>
				<input
					type="text"
					placeholder="Rechercher par nom..."
					bind:value={nameFilter}
					on:input={() => {
						page = 0;
						fetchRefills();
					}}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
				/>
			</div>
            <div class="flex items-center gap-3">
                <label for="state" class="text-sm font-medium text-gray-700 dark:text-gray-300">État:</label>
                <select 
                    class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm dark:text-white"
                    name="state"
                    bind:value={state}
                    on:change={() => {
                        page = 0;
                        fetchRefills();
                    }}
                >
                    <option value={undefined}></option>
                    {#each Object.values(RemoteRefillState) as state}
                        <option value={state}>{state_names[state]}</option>
                    {/each}
                </select>
            </div>
		</div>
	</div>

    <div class="flex-grow w-full overflow-x-auto min-w-full">
        <table class="w-full min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm">
            <thead class="bg-gray-50 dark:bg-gray-700">
            <tr class="divide-x divide-gray-200 dark:divide-gray-700">
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                Date
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                Compte
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                Montant
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider opacity-60 w-1">
                État
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider opacity-60 w-1">
                Action
                </th>
            </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            {#each refills as refill}
                <tr class="divide-x divide-gray-200 dark:divide-gray-700">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                    {new Date(refill.created_at * 1000).toLocaleDateString('fr-FR', {
                    day: 'numeric',
                    month: 'long',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                    })}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                    {refill.account_name}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-300">
                    {refill.amount / 100 } €
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                    {state_names[refill.state]}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                    {#if refill.state == RemoteRefillState.RemoteRefillAbandoned || refill.state == RemoteRefillState.RemoteRefillStarted}
                    <button 
                        class="relative rounded-sm bg-green-500 hover:bg-green-600 disabled:hover:bg-green-500 text-white dark:text-white p-2 text-left flex items-center gap-2"
                        on:click={() => {verifyRefill(refill.id)}}
                        disabled={loading || !remote_refills_available}
                    >
                        <iconify-icon icon="mdi:refresh" width="20" height="20" />
                        <div class="align-middle">Vérifier</div>
                        {#if !remote_refills_available}
                        <div class="absolute w-full h-full top-0 left-0 flex justify-center items-center bg-gray-600 bg-opacity-50 rounded-sm"></div>
                        {/if}
                        {#if loading}
                        <div class="absolute w-full h-full top-0 left-0 flex justify-center items-center bg-gray-600 bg-opacity-50 rounded-sm">
                            <Spinner size="23px"/>
                        </div>
                        {/if}
                    </button>
                    {/if}
                </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

	<!-- Pagination -->
	<div class="sticky bottom-0 left-0 right-0 mt-4 px-6 py-4 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 flex flex-col sm:flex-row justify-between items-center gap-4">
		<div>
			<p class="text-sm text-gray-600 dark:text-gray-400">
				<span class="font-semibold text-gray-800 dark:text-gray-200">{refills.length}</span> résultats
			</p>
		</div>
		<div class="flex items-center gap-x-4">
			<button
				type="button"
				class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-offset-gray-800"
				on:click={prevPage}
				disabled={page === 0}
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				Précédent
			</button>
			<p class="text-sm font-medium text-gray-700 dark:text-gray-300">
				Page <span class="font-bold">{page}</span> sur <span class="font-bold">{maxPage}</span>
			</p>
			<button
				type="button"
				class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-offset-gray-800"
				on:click={nextPage}
				disabled={page === maxPage}
			>
				Suivant
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		</div>
	</div>
</div>
