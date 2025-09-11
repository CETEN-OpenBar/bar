<script lang="ts">
	import { goto } from '$app/navigation';
	import Habutton from '$lib/components/client/habutton.svelte';
	import Spinner from '$lib/components/spinner.svelte';
	import { refillsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	let value: string = '';

    let remote_refills_available: boolean = true;
    onMount(() => {
        refillsApi().getRemoteRefillStatus({withCredentials: true})
        .then((resp) => {
            remote_refills_available = resp.status == 200;
        })
        .catch((_) => {
            remote_refills_available = false;
        });
    })

	let isValid: boolean = false;
	let loading: boolean = false;
	let error: string | undefined;

	// Regex: positive integer or float with up to 2 decimals
	const regex: RegExp = /^(?:\d+)(?:\.\d{1,2})?$/;

	function validate(): void {
		isValid = regex.test(value) && parseFloat(value) > 0.5;
	}

	function startRemoteRefill(): void {
        if (!remote_refills_available) return;
		if (!isValid) return;

		loading = true;
		error = undefined;

        // Convert amount to cents
        let amount: number = parseFloat(value) * 100

        // Get redirect url
        refillsApi().startRemoteRefill(amount, {withCredentials: true}).then(
            (resp) => {
                goto(resp.data.redirect_url);
            }
        ).catch((reason) => {
            error = reason.message;
            loading = false;
        })
	}
</script>

<div class="text-center p-3">
	<a
		class="text-xl lg:text-2xl font-extrabold underline underline-offset-4 text-white hover:text-gray-400"
		href="/client/index"
	>
		&lt&nbsp Retour&nbsp
	</a>
</div>

{#if !remote_refills_available}
<div class="w-full flex flex-col items-center px-5 mb-5">
    <div class="bg-red-200 border-red-400 border-4 rounded-lg p-3">
        <div class="font-bold text-center text-lg">Les rechargements en ligne sont actuellement indisponibles</div>
        <div class="text-center">
            Merci de réessayer plus tard
        </div>
    </div>
</div>
{/if}

<div class="flex justify-center">
	<form
		class="w-[90%] lg:max-w-3xl bg-gray-200 rounded py-3 px-5 grid gap-5 items-center justify-items-center relative"
		on:submit|preventDefault={startRemoteRefill}
	>
        {#if !remote_refills_available}
            <div class="absolute w-full h-full top-0 left-0 bg-gray-700 rounded opacity-50"></div>
        {/if}
		<h1 class="text-xl font-bold text-center">Recharger</h1>

		<label class="flex flex-col lg:flex-row justify-center lg:gap-3 lg:items-center">
			Montant
			<div class="flex flex-row items-center gap-3">
				<input
					class="h-8 rounded p-2 flex-grow"
					bind:value
					on:input={validate}
					type="number"
					step="0.01"
					min="0.50"
					placeholder="10.0"
                    disabled={!remote_refills_available}
				/>
				<div>€</div>
			</div>
		</label>

        <div class="text-red-500 text-center">
            N'oublie pas de retirer la contribution à HelloAsso !
        </div>

		<Habutton disabled={!isValid || loading} type="submit" />

		{#if loading}
            <div class="mb-3">
                <Spinner />
            </div>
		{/if}

		{#if error}
			<div class="p-3 bg-red-100 border border-red-400 text-red-700 rounded mb-3">
				{{ error }}
			</div>
		{/if}
	</form>
</div>
