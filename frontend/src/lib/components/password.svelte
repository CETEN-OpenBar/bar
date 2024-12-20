<script lang="ts">
	let password = '';
	export let callback: (pin: string) => void;
	export let custom_text: string = 'Entrez votre mot de passe';

	// onEnter
	function onEnter(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			// console.log('a' + password);
			callback(password);
			password = '';
		}
	}
	let warnCapsLockOn: boolean = false

    function checkCapsLock(e: KeyboardEvent) {
		warnCapsLockOn = e.getModifierState("CapsLock")
    }
</script>

<svelte:document on:keydown={onEnter}/>

<!-- Display a popup that asks for a pin -->

<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-40"
	on:click={() => {
		password = '';
		callback('');
		console.log('click');
	}}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<!-- Put a title and the numpad -->
	<div class="flex flex-col items-center bg-neutral-700 rounded-lg shadow-lg p-4 z-40">
		<h1 class="text-2xl font-bold mb-4 text-white">{custom_text}</h1>
		{#if warnCapsLockOn}
			<h2 class="text-2xl font-bold mb-4 text-red-600">Attention, caps lock on!</h2>
		{/if}
		<div class="flex flex-col items-center">
			<!-- Display the numpad -->
			<!-- svelte-ignore a11y-autofocus -->
			<input
				name="password"
				bind:value={password}
				on:keyup={checkCapsLock}
				type="password"
				class="w-full h-20 m-3 bg-neutral-800 rounded-xl text-white text-4xl text-center"
				autofocus
			/>

			<button
				class="w-36 h-12 text-xl rounded-xl text-white bg-green-600 active:bg-green-800"
				on:click={() => {
					console.log('g' + password);
					callback(password);
					password = '';
				}}
			>
				Valider
			</button>
		</div>
	</div>
</div>
