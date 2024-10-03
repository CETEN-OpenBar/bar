<script lang="ts">
	import { dev } from '$app/environment';
	import { pinIsShown } from '$lib/store/store';
	import { onDestroy } from 'svelte';
	export let callback: (card: string) => void = () => {};
	$: $pinIsShown

	let socket = new WebSocket('ws://localhost:3737/');

	function defineSocket(socket: WebSocket) {
		type WsData = {
			uid: string;
		};

		socket.onmessage = (event) => {
			const data: WsData = JSON.parse(event.data);
			callback(data.uid);
		};

		// on any error the websocket should try to reconnect 1 second later
		socket.onerror = () => {
			setTimeout(() => {
				socket = new WebSocket('ws://localhost:3737/');
				defineSocket(socket);
			}, 1000);
		};
	}

	onDestroy(() => {
		socket.close();
		// remove the event listener
		window.removeEventListener('keydown', handleInput);
	});

	let buffer = '';
	function handleInput(event: KeyboardEvent) {
		if (event.key === 'Enter' && !$pinIsShown) {
			callback(buffer);
			buffer = '';
		} else if (!$pinIsShown){
			buffer += event.key;
		}
	}

	defineSocket(socket);
</script>

<svelte:window on:keydown={handleInput} />
