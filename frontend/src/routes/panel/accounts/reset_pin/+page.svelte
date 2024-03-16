<script lang="ts">
	import type { Account } from '$lib/api';
	import { accountsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
    import Error from '$lib/components/error.svelte';
    import Success from '$lib/components/success.svelte';
	import { goto } from '$app/navigation';

	let accounts: Account[] = [];
	let selectedAccount: Account | null = null;

    let success = '';
	let error = '';

	async function getAccounts(search: string) {
		let res = await accountsApi().getAccounts(0, 5, search, { withCredentials: true });
		if (Array.isArray(res.data.accounts)) accounts = res.data.accounts;
	}

	async function resetPin() {
		if (!selectedAccount) return;
		await accountsApi().resetAccountPin(selectedAccount.id, { withCredentials: true })
			.then((res) => {
				if (res.status == 200) {
					success = 'Code pin réinitialisé'
				} else {
					error = 'Une erreur est survenue'
				}
			})
			.catch((err) => {
				error = 'Une erreur est survenue'
			});
        setTimeout(() => {
            goto("/panel");
        }, 2000);
	}

	onMount(async () => {
		await getAccounts('');
	});
</script>

{#if success != ''}
	<Success message={success} />
{/if}

{#if error != ''}
	<Error {error} />
{/if}

<div class="w-full flex flex-row justify-center mt-5">
	<div class="w-[50%] flex flex-col gap-5">
		<div class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
			<div class="text-3xl text-center font-bold p-5">Sélection du compte</div>

			<div class="flex flex-row justify-center">
				<input
					type="text"
					class="w-1/2 rounded-lg p-2 border-2 border-slate-200"
					placeholder="nom du compte"
					on:keyup={async (e) => {
						// @ts-ignore
						let val = e.target?.value;
						await getAccounts(val);
					}}
				/>
			</div>

			<div class="grid grid-cols-4 items-center gap-3 overflow-x-auto mt-2 p-5 h-56">
				{#each accounts as a}
					<button
						class="flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg h-24 justify-center"
						on:click={() => {
							selectedAccount = a;
						}}
					>
						{a.first_name}
						{a.last_name}
					</button>
				{/each}
			</div>
		</div>

		{#if selectedAccount}
			<div class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
				<div class="text-3xl text-center font-bold p-5">
					{selectedAccount?.first_name}
					{selectedAccount?.last_name}
					{#if selectedAccount?.nickname}
						(Surnom : {selectedAccount?.nickname})
					{/if}
				</div>

				<button
					class="rounded-md p-2 text-center bg-red-200 hover:bg-red-300 mx-32 mt-5"
					on:click={() => {
						resetPin();
					}}
				>
					Réinitialiser le code PIN
				</button>
			</div>
		{/if}
	</div>
</div>
