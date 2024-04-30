<script lang="ts">
	import type { Account, ConnectCard200Response, ConnectCardRequest } from '$lib/api';
	import { accountsApi, authApi } from '$lib/requests/requests';
	import Button from '$lib/components/button.svelte';

	let account: Account | undefined = undefined;

	const callback = async () => {
		let req: ConnectCardRequest = {
			card_id: '1',
			card_pin: '1111'
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
						let data: ConnectCard200Response = res.data;
						account = data.account;
					});
			});
	};
</script>

<h1 class="text-3xl font-bold underline">Hello world!</h1>

<Button class="text-white bg-blue-500 hover:bg-blue-700" text="Click me!" {callback} />

{#if account}
	<p class="text-xl">You are logged in as {account.id}.</p>
{:else}
	<p class="text-xl">You are not logged in.</p>
{/if}
