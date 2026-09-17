<script lang="ts">
	import { goto } from '$app/navigation';
	import { authApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	export let delay: number = 60000;
	export let location: string = '/borne';

	onMount(() => {
		let mounted = true;
		let disconnecting = false;
		let disconnectTimer: ReturnType<typeof setTimeout> | undefined;
		const events = [
			'mousemove',
			'mousedown',
			'keypress',
			'DOMMouseScroll',
			'mousewheel',
			'touchmove',
			'MSPointerMove',
			'click',
			'drag',
			'dragend',
			'dragenter',
			'dragleave',
			'dragover',
			'dragstart',
			'touchstart',
			'touchend',
			'gesturestart',
			'gesturechange',
			'gestureend'
		];

		async function logout() {
			if (!mounted || disconnecting) return;
			disconnecting = true;
			clearTimeout(disconnectTimer);
			try {
				await authApi().logout({ withCredentials: true });
				if (mounted) goto(location);
			} catch (error) {
				if (mounted) {
					console.error('Impossible de fermer la session', error);
					disconnecting = false;
					onAction();
				}
			}
		}

		function onAction() {
			if (!mounted || disconnecting) return;
			clearTimeout(disconnectTimer);
			disconnectTimer = setTimeout(logout, delay);
		}

		onAction();
		for (const event of events) {
			window.addEventListener(event, onAction);
		}

		return () => {
			mounted = false;
			clearTimeout(disconnectTimer);
			for (const event of events) {
				window.removeEventListener(event, onAction);
			}
		};
	});
</script>
