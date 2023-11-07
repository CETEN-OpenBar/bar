<script lang="ts">
	import type { Account, UpdateAccountAdmin } from '$lib/api';
	import { accountsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	let accounts: Account[] = [];
	let selectedAccount: Account | null = null;

	async function getAccounts(search: string) {
		let res = await accountsApi().getAccounts(0, 5, search, { withCredentials: true });
		if (Array.isArray(res.data.accounts)) accounts = res.data.accounts;
	}

    function isAccountBlocked(account: Account) {
        for (let restriction of account.restrictions??[]) {
            if (restriction === "blocked") return true;
        }
    }

    async function toggleBlock(account: Account) {
        if (isAccountBlocked(account)) {
            // remove blocked from restrictions
            account.restrictions = account.restrictions.filter((r) => r !== "blocked");
            let update: UpdateAccountAdmin = {
                restrictions: account.restrictions
            };
            await accountsApi().patchAccountId(account.id, update, { withCredentials: true });
        } else {
            // block
            if (!account.restrictions) account.restrictions = [];
            account.restrictions.push("blocked");
            let update: UpdateAccountAdmin = {
                restrictions: account.restrictions
            };
            await accountsApi().patchAccountId(account.id, update, { withCredentials: true });
        }

        // fetch account to update
        let res = await accountsApi().getAccountId(account.id, { withCredentials: true });
        if (res.data) selectedAccount = res.data;
    }

	onMount(async () => {
		await getAccounts('');
	});
</script>

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
                    {#if isAccountBlocked(selectedAccount)}
                        (bloqué)
                    {/if}
				</div>

				<button
					class="rounded-sm p-5 text-center text-lg h-24 justify-center {isAccountBlocked(selectedAccount) ? 'bg-red-200 hover:bg-red-300' : 'bg-green-200 hover:bg-green-300'}"
                    on:click={() => {
                        if (selectedAccount) toggleBlock(selectedAccount);
                    }}
				>
                {#if isAccountBlocked(selectedAccount)}
                    Débloquer
                {:else}
                    Bloquer
                {/if}
				</button>
			</div>
		{/if}
	</div>
</div>
