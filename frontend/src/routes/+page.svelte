<script lang="ts">
	import type { Account, Callback200Response, ConnectCardRequest } from '$lib/api';
	import { accountsApi, authApi } from '$lib/requests/requests';
	import Button from '$lib/components/button.svelte';
	import { goto } from '$app/navigation';

	let account: Account | undefined = undefined;

	const callback = async () => {
		let req: ConnectCardRequest = {
			card_id: '11',
			card_pin: '1234'
		};
		authApi()
			.connectCard(req, {
				withCredentials: true
			})
			.then(() => {
				// Get account
				accountsApi()
					.getAccount({
						withCredentials: true
					})
					.then((res: any) => {
						let data: Callback200Response = res.data;
						account = data.account;
					});
			});
	};
	goto("auth");
</script>

<Button class="text-white bg-blue-500 hover:bg-blue-700" text="Click me!" {callback} />

{#if account}
	<p class="text-xl">You are logged in as {account.first_name} {account.last_name}.</p>
{:else}
	<p class="text-xl">You are not logged in.</p>
{/if}
