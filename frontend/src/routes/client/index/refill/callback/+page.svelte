<script lang="ts">
	import { page } from '$app/stores';
	import Spinner from '$lib/components/spinner.svelte';
	import { accountsApi, refillsApi } from '$lib/requests/requests';
	import { store } from '$lib/store/store';
	import { onMount } from 'svelte';

	// derive orderId and code from the query parameters
	const checkoutIntentId = $page.url.searchParams.get('checkoutIntentId');
	const code = $page.url.searchParams.get('code');

    let loading: boolean = false;
    let success: boolean = false;
    let error: string | null = $page.url.searchParams.get('error');
    let show_disclaimer: boolean = false;

    // Try to validate the payment on load
    onMount(() => {

        if (error != null) {
            return;
        }

        if (code == undefined || checkoutIntentId == undefined || isNaN(parseInt(checkoutIntentId))) {
            error = "Bad Request";
            return;
        }

        if (code != "succeeded") {
            error = "Impossible de vérifier le paiement";
            return;
        }

        loading = true;
        refillsApi().selfValidateRemoteRefill(parseInt(checkoutIntentId), {withCredentials: true}).then(
            (res) => {
                if (res.status != 200) {
                    throw res;
                }
                success = true;
                
                // Refresh account data in layout
                accountsApi()
                    .getAccount({
                        withCredentials: true
                    })
                    .then((res) => {
                        store.set({ account: res.data.account });
                    })
                    .catch((_) => {})
            }
        ).catch((err) => {
            switch (err.status) {
                case 409:
                    error = "Ce rechargement a déjà été validé"
                    break;
                case 404:
                    error = "Rechargement introuvable"
                    break;
                case 402:
                    error = "Le paiment n'a pas encore été validé."
                    show_disclaimer = true
                    break;
                default:
                    error = "Une erreur innatendue est survenue"
                    break;
            }
        })
        .finally(() => {loading = false;})
    });
</script>

{#if loading}
<Spinner />
{/if}

<div class="grid justify-items-center mt-5 gap-5">

    {#if success}
        <div
            class="rounded p-5 border-2 border-green-500 bg-green-400 text-white text-lg font-bold text-center max-w-[75%]"
        >
            Rechargement effectué !
        </div>
    {/if}

    {#if error}
        <div
            class="rounded p-5 border-2 border-red-500 bg-red-400 text-white text-lg font-bold text-center max-w-[75%]"
        >
            <div>
                <span class="font-extrabold">Erreur : </span>{error}
            </div>
            {#if show_disclaimer}    
                <div class="font-normal italic mt-2">
                    Remarque : le paiement sera revérifié automatiquement.<br/>
                    HelloAsso peut mettre jusqu'à 30 minutes pour valider un paiement, au delà de ce délai, contactez un membre du bar.
                </div>
            {/if}
        </div>
    {/if}
    
	<a class="text-lg font-bold text-white underline" href="/client/index"> Retour </a>
</div>
