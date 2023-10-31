<script lang="ts">
	import { dev } from '$app/environment';
	export let callback: (card: string) => void = () => {};

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

	let buffer = '';
	function handleInput(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			callback(buffer);
			buffer = '';
		} else {
			buffer += event.key;
		}
	}

	defineSocket(socket);
</script>

<svelte:window on:keydown={handleInput} />
