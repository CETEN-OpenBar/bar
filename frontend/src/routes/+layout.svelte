<script>
	import '../app.css';
	import 'preline/dist/preline.js';
	import { loadConfig } from '$lib/config/config';

	let loading = true;

	loadConfig().then(() => {
		loading = false;
	})
</script>

{#if loading}
	loading...
{:else}
	<slot data-sveltekit-preload-data="off" />
{/if}
