<script>
	// If we don't have local_token, go back to /
	import { local_token } from '$lib/config/config';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		if (local_token() === '') {
			goto('/');
		}
	});
</script>

<slot />
