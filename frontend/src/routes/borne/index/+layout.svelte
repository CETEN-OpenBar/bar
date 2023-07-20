<script>
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	onMount(() => {
		accountsApi()
			.getAccount({
				withCredentials: true
			})
			.catch((err) => {
				goto('/borne');
			});
	});
</script>

<slot />
