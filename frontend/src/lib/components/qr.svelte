<script lang="ts">
	import { api } from '$lib/config/config';
	import { authApi } from '$lib/requests/requests';
	import Error from './error.svelte';
	import Pin from './pin.svelte';

    let qr: any = undefined;
    let ask = true;
	let error = '';

	function pinCallback(pin: string) {
        if (pin == '') {
            ask = false;
            return;
        }
		authApi()
			.getAccountQR(pin, { withCredentials: true })
			.then((res) => {
                qr = res.data;
			})
			.catch(() => {
				error = 'Impossible de générer le QR code';
				setTimeout(() => {
					error = '';
				}, 3000);
			});
	}
</script>

{#if qr !== undefined && ask}
	<img class="w-64" src="data:image/png;base64,{qr}" alt="je charge..." />
{:else}
	<!-- <Pin callback={pinCallback} /> -->
    {pinCallback("1234")}
{/if}

{#if error !== ''}
	<Error {error} />
{/if}
