<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi, authApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import 'iconify-icon';

	let account: Account | undefined = undefined;

	onMount(() => {
		accountsApi()
			.getAccount({
				withCredentials: true
			})
			.then((res) => {
				account = res.data.account;
				store.set({ account });
			})
			.catch(() => {
				goto('/comptoir');
			});

		disconnectInterval = setInterval(logout, 15000);

		// trigger action on any event
		let events = ['mousemove', 'mousedown', 'keypress', 'DOMMouseScroll', 'mousewheel', 'touchmove', 'MSPointerMove', 'click'];
		for (let i in events) {
			window.addEventListener(events[i], onAction);
		}
	});

	let disconnectInterval: number | undefined = undefined;

	function logout() {
		authApi()
			.logout({ withCredentials: true })
			.then(() => {
				clearInterval(disconnectInterval);
				goto('/comptoir');
			});
	}

	function onAction() {
		console.log('action');
		clearInterval(disconnectInterval);
		disconnectInterval = setInterval(logout, 15000);
	}
</script>


{#if account !== undefined}
	<div id="main" class="absolute top-0 left-0 w-screen h-screen" style="background-color:black;">
		<slot />
	</div>
{/if}
