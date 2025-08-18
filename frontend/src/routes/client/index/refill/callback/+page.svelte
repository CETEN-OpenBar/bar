<script lang="ts">
	import { page } from '$app/stores';
	import Spinner from '$lib/components/spinner.svelte';
	import { refillsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import { derived } from 'svelte/store';

	// derive orderId and code from the query parameters
	const checkoutIntentId = $page.url.searchParams.get('checkoutIntentId');
	const code = $page.url.searchParams.get('code');

    let loading: boolean = false;
    let success: boolean = false;
    let error: string | null = $page.url.searchParams.get('error');

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

                switch (res.status) {
                    case 200:
                        success = true;
                        break;
                    case 409:
                        error = "Ce rechargement a déjà été validé"
                        break;
                    case 404:
                        error = "Rechargement introuvable"
                        break;
                    case 402:
                        error = "Le paiment n'a pas encore été validé, merci de réessayer plus tard"
                        break;
                    default:
                        error = "Une erreur innatendue est survenue"
                        break;
                }

                loading = false;
            }
        ).catch((err) => {
            loading = false;
            error = "Une erreur innatendue est survenue";
        })
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
            <span class="font-extrabold">Erreur : </span>{error}
        </div>
    {/if}
    
	<a class="text-lg font-bold text-white underline" href="/client/index"> Retour </a>
</div>
