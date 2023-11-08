<script lang="ts">
	import { goto } from '$app/navigation';
	import { authApi } from '$lib/requests/requests';
	import { onDestroy, onMount } from 'svelte';

	export let delay: number = 60000;
	export let location: string = '/borne';

	let disconnectInterval: number | undefined = undefined;

	onMount(() => {
        disconnectInterval = setInterval(logout, delay);

        // trigger action on any event
        let events = [
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
            'gestureend',
        ];
        for (let i in events) {
            window.addEventListener(events[i], onAction);
        }
    });

    onDestroy(() => {
        clearInterval(disconnectInterval);
    });

	async function logout() {
		// TODO: don't forget to uncomment this
		let _ = await authApi().logout({ withCredentials: true });
		clearInterval(disconnectInterval);
		goto(location);
	}

	function onAction() {
		clearInterval(disconnectInterval);
		disconnectInterval = setInterval(logout, delay);
	}
</script>
