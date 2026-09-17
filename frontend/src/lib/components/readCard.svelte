<script lang="ts">
	import { onMount } from 'svelte';
	export let callback: (card: string) => void = () => {};

	onMount(() => {
		let mounted = true;
		let socket: WebSocket | undefined;
		let reconnectTimer: ReturnType<typeof setTimeout> | undefined;

		function closeSocket() {
			if (!socket) return;
			socket.onmessage = null;
			socket.onerror = null;
			socket.onclose = null;
			socket.close();
			socket = undefined;
		}

		function reconnect() {
			if (!mounted || reconnectTimer !== undefined) return;
			reconnectTimer = setTimeout(() => {
				reconnectTimer = undefined;
				if (!mounted) return;
				closeSocket();
				connect();
			}, 1000);
		}

		function connect() {
			const currentSocket = new WebSocket('ws://localhost:3737/');
			socket = currentSocket;
			currentSocket.onmessage = (event) => {
				if (!mounted || socket !== currentSocket) return;
				const data: { uid: string } = JSON.parse(event.data);
				callback(data.uid);
			};
			currentSocket.onerror = reconnect;
			currentSocket.onclose = reconnect;
		}

		connect();
		return () => {
			mounted = false;
			clearTimeout(reconnectTimer);
			closeSocket();
			buffer = '';
		};
	});

	let buffer = '';
	function handleInput(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			callback(buffer);
			buffer = '';
		} else {
			buffer += event.key;
		}
	}
</script>

<svelte:window on:keydown={handleInput} />
