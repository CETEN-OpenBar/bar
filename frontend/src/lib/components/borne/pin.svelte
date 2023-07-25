<script lang="ts">
	let temp_pin = '';
	export let callback: (pin: string) => void;
	export let custom_text: string = 'Entrez votre code pin';
</script>

<!-- Display a popup that asks for a pin -->

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
	<div class="flex flex-col items-center bg-white rounded-lg shadow-lg p-4 z-40">
		<h1 class="text-2xl font-bold mb-4">{custom_text}</h1>
		<div class="flex flex-col items-center">
			<!-- Display the pin -->
			<div class="flex flex-row mb-4">
				{#each temp_pin.split('') as i}
					<p class="font-weight-bold text-3xl">*</p>
				{/each}
				{#if temp_pin == ''}
					<p class="font-weight-bold text-3xl text-transparent">*</p>
				{/if}
			</div>
			<!-- Display the numpad -->
			<div class="grid grid-cols-3 gap-2">
				{#each Array.from({ length: 9 }, (_, i) => i + 1) as i}
					<button
						class="w-16 h-16 text-xl border-2 border-gray-300 rounded-full hover:bg-gray-200 active:bg-gray-200"
						on:click={() => {
							temp_pin += i;
							if (temp_pin.length === 15) {
								callback(temp_pin);
								temp_pin = '';
							}
						}}
					>
						{i}
					</button>
				{/each}
				<button
					class="w-16 h-16 text-xl border-2 border-gray-300 rounded-full"
					on:click={() => {
						temp_pin = temp_pin.slice(0, -1);
					}}
				>
					←
				</button>
				<button
					class="w-16 h-16 text-xl border-2 border-gray-300 rounded-full"
					on:click={() => {
						temp_pin += 0;
						if (temp_pin.length === 15) {
							callback(temp_pin);
							temp_pin = '';
						}
					}}
				>
					0
				</button>
				<button
					class="w-16 h-16 text-xl border-2 border-gray-300 rounded-full"
					on:click={() => {
						callback(temp_pin);
						temp_pin = '';
					}}
				>
					✔
				</button>
			</div>
		</div>
	</div>
</div>
