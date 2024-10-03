<script lang="ts">
	let temp_pin = '';
	export let callback: (pin: string) => void;
	export let custom_text: string = 'Entrez votre code pin';
	function enterPin(i: number) {
		temp_pin += i;
		if (temp_pin.length === 20) {
			console.log('e' + temp_pin);
			callback(temp_pin);
			temp_pin = '';
		}
	}

	function deletePinChar() {
		temp_pin = temp_pin.slice(0, -1);
	}

	function validatePin() {
		callback(temp_pin);
		temp_pin = '';
	}

	function controlPinWithKbd(k: string) {
		switch (k) {
			case '1':
				enterPin(1);
				break;
			case '2':
				enterPin(2);
				break;
			case '3':
				enterPin(3);
				break;
			case '4':
				enterPin(4);
				break;
			case '5':
				enterPin(5);
				break;
			case '6':
				enterPin(6);
				break;
			case '7':
				enterPin(7);
			case '8':
				enterPin(8);
				break;
			case '9':
				enterPin(9);
				break;
			case '0':
				enterPin(0);
				break;
			case 'Backspace':
				deletePinChar();
				break;
			case 'Enter':
				validatePin();
				break;
			default:
				break;
		}
	}
</script>

<!-- Display a popup that asks for a pin -->

<svelte:window on:keydown={(e) => controlPinWithKbd(e.key)} />

<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-40"
	on:click={() => {
		temp_pin = '';
		callback('');
		console.log('click');
	}}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<!-- Put a title and the numpad -->
	<div class="flex flex-col items-center bg-neutral-700 rounded-lg shadow-lg p-4 z-40">
		<h1 class="text-2xl font-bold mb-4 text-white">{custom_text}</h1>
		<div class="flex flex-col items-center">
			<!-- Display the pin -->
			<input
				name="pin"
				bind:value={temp_pin}
				disabled
				type="password"
				class="w-full h-20 m-3 bg-neutral-800 rounded-xl text-white text-4xl text-center"
			/>
			<!-- Display the numpad -->
			<div class="grid grid-cols-3 gap-2">
				{#each Array.from({ length: 9 }, (_, i) => i + 1) as i}
					<button
						class="w-36 h-32 text-4xl border-2 text-white border-gray-300 rounded-xl hover:bg-gray-200/[0.5] active:bg-gray-200/[0.5]"
						on:click={() => enterPin(i)}
					>
						{i}
					</button>
				{/each}
				<button
					class="w-36 h-32 text-4xl rounded-xl bg-yellow-600 text-white"
					on:click={() => deletePinChar}
				>
					←
				</button>
				<button
					class="w-36 h-32 text-4xl border-2 rounded-xl text-white border-gray-300 hover:bg-gray-200 active:bg-gray-200"
					on:click={() => enterPin(0)}
				>
					0
				</button>
				<button
					class="w-36 h-32 text-xl rounded-xl text-white bg-green-600 active:bg-green-800"
					on:click={() => validatePin}
				>
					Valider
				</button>
			</div>
		</div>
	</div>
</div>
