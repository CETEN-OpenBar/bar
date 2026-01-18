<script lang="ts">
	import type { Account, RemoteRefill } from '$lib/api';
	import { onMount, onDestroy } from 'svelte';
	import { store } from '$lib/store/store';
	import { fly, scale } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';
	import 'iconify-icon';
	import Qr from '$lib/components/client/qr.svelte';
	import Transactions from '$lib/components/client/transactions.svelte';
	import Refills from '$lib/components/client/refills.svelte';
	import { goto } from '$app/navigation';
	import { refillsApi } from '$lib/requests/requests';

	let account: Account | undefined = undefined;
    let pending_refills: RemoteRefill[] = [];
    let remote_refills_available: boolean = true;
	let unsub: () => void;

	onMount(() => {
		unsub = store.subscribe((state) => {
			account = state.account;
		});
        refillsApi().getPendingRemoteRefills({withCredentials: true})
        .then((resp) => {
            pending_refills = resp.data.remote_refills ?? [];
        })
        refillsApi().getRemoteRefillStatus({withCredentials: true})
        .then((resp) => {
            remote_refills_available = resp.status == 200;
        })
        .catch((_) => {
            remote_refills_available = false;
        })
	});

	onDestroy(() => {
		unsub();
	});

	let showPreviousOrders = false;
	let showPreviousRefills = false;
	let showLinkWithGoogle = false;

	function togglePreviousOrders() {
		showPreviousOrders = !showPreviousOrders;
	}

	function togglePreviousRefills() {
		showPreviousRefills = !showPreviousRefills;
	}

	function toggleLinkWithGoogle() {
		showLinkWithGoogle = !showLinkWithGoogle;
	}
</script>

<svelte:head>
    <title>OpenBar - Mon compte</title>
</svelte:head>


{#if account !== undefined}
    <div class="flex justify-center">
        <button 
            class="bg-green-500 hover:bg-green-700 rounded p-2 m-3 text-lg font-bold text-white lg:w-1/2 lg:h-20 disabled:bg-gray-400 disabled:text-gray-700"
            on:click={() => {
				goto('/client/index/refill');
			}}
            disabled={!remote_refills_available}
        >
            Recharger mon compte
        </button>
    </div>
    {#if !remote_refills_available}
    <div class="w-full flex flex-col items-center px-5">
        <div class="bg-red-200 border-red-400 border-4 rounded-lg p-3">
            <div class="font-bold text-center text-lg">Les rechargements en ligne sont actuellement indisponibles</div>
            <div class="text-center">
                Merci de réessayer plus tard
            </div>
        </div>
    </div>
    {/if}

    {#if pending_refills.length > 0}
    <div class="w-full flex flex-col items-center px-5">
        <div class="bg-yellow-200 border-yellow-400 border-4 rounded-lg p-3">
            <div class="font-bold text-center text-lg">Vous avez {pending_refills.length} paiement{pending_refills.length > 1 ? 's' : ''} en attente de validation</div>
            <div class="my-2 text-center">
                {#each pending_refills as refill}
                    <div>
                        Le {new Date(refill.created_at * 1000).toLocaleDateString('fr-FR', {
                        day: 'numeric',
                        month: 'long',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit'
                        })} : {refill.amount / 100} €
                    </div>
                {/each}
            </div>
            <div class="text-center">
                Les paiements sont revérifiés automatiquement.<br>
                HelloAsso peut mettre jusqu'à 30 minutes pour valider un paiement, au delà de ce délai, contactez un membre du bar.
            </div>
        </div>
    </div>
    {/if}

	<div class="grid lg:grid-cols-3 grid-cols-1 gap-16 p-5 w-full">
		<!-- Previous orders column -->
		<div
			class="flex flex-col flex-grow transition-all ease-in-out"
			transition:scale={{ delay: 250, duration: 300, easing: quintOut }}
		>
			<button
				class="p-4 flex flex-col text-lg font-semibold {showPreviousOrders
					? 'rounded-t'
					: 'rounded'}  cursor-pointer z-10"
				style="background-color:#EEEEEE"
				on:click={togglePreviousOrders}
			>
				Commandes passées
				<iconify-icon class="text-white text-2xl self-center" icon="flat-color-icons:expand" />
			</button>
			{#if showPreviousOrders}
				<div
					class="flex flex-col rounded-br rounded-bl p-10"
					style="background-color:#EEEEEE"
					transition:fly={{ y: -80, duration: 300, easing: quintOut }}
				>
					<!-- Add content for previous orders here -->
					<hr class="my-2 border-gray-400" />

					<Transactions />
				</div>
			{/if}
		</div>

		<!-- Link with Google column -->
		<div
			class="flex flex-col flex-grow transition-all ease-in-out"
			transition:scale={{ delay: 300, duration: 300, easing: quintOut }}
		>
			<button
				class="p-4 flex flex-col text-lg font-semibold {showLinkWithGoogle
					? 'rounded-t'
					: 'rounded'}  cursor-pointer z-10"
				style="background-color:#EEEEEE"
				on:click={toggleLinkWithGoogle}
			>
				Connexion avec Google
				<iconify-icon class="text-white text-2xl self-center" icon="flat-color-icons:expand" />
			</button>
			{#if showLinkWithGoogle}
				<div
					class="flex flex-col rounded-br rounded-bl p-10"
					style="background-color:#EEEEEE"
					transition:fly={{ y: -80, duration: 300, easing: quintOut }}
				>
					<!-- Add content for previous orders here -->
					<hr class="my-2 border-gray-400" />
					<div class="w-full flex justify-center">
						<div class="flex flex-col gap-5">
							<Qr />
							Scan ce code pour lier ton compte !
						</div>
					</div>
				</div>
			{/if}
		</div>

		<!-- Previous refills column -->
		<div
			class="flex flex-col flex-grow transition-all ease-in-out"
			transition:scale={{ delay: 300, duration: 300, easing: quintOut }}
		>
			<button
				class="p-4 flex flex-col text-lg font-semibold {showPreviousRefills
					? 'rounded-t'
					: 'rounded'}  cursor-pointer z-10"
				style="background-color:#EEEEEE"
				on:click={togglePreviousRefills}
			>
				Transactions
				<iconify-icon class="text-white text-2xl self-center" icon="flat-color-icons:expand" />
			</button>
			{#if showPreviousRefills}
				<div
					class="flex flex-col rounded-br rounded-bl p-10"
					style="background-color:#EEEEEE"
					transition:fly={{ y: -80, duration: 300, easing: quintOut }}
				>
					<!-- Add content for previous orders here -->
					<hr class="my-2 border-gray-400" />

					<Refills />
				</div>
			{/if}
		</div>
	</div>
{/if}
