<script lang="ts">
	import type { Account } from '$lib/api';
	import { onMount, onDestroy } from 'svelte';
	import { store } from '$lib/store/store';
	import { fly, scale } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';
	import 'iconify-icon';
	import Qr from '$lib/components/qr.svelte';

	let account: Account | undefined = undefined;
	let unsub: () => void;

	onMount(() => {
		unsub = store.subscribe((state) => {
			account = state.account;
		});
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

{#if account !== undefined}
	<div class="grid grid-cols-3 gap-16 p-16 w-full">
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
				</div>
			{/if}
		</div>
	</div>
{/if}
