<script lang="ts">
	export let callback: (response: Boolean) => void;
	export let custom_text: string = 'Êtes-vous sûr ?';
</script>

<!-- Display a popup that asks for a pin -->

<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-40"
	on:click={() => {
		callback(false);
	}}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<!-- Put a title and the numpad -->
	<div class="flex flex-col items-center bg-white rounded-lg shadow-lg p-4 z-40">
		<h1 class="text-2xl font-bold mb-4">{custom_text}</h1>
		<div class="flex flex-col items-center">
			<!-- Display the numpad -->	
			<div class="grid grid-cols-3 gap-2">
				<button
					class="w-16 h-16 text-xl border-2 border-gray-300 rounded-full"
					on:click={() => {
						callback(false);
					}}
				>
					<iconify-icon class="text-red align-middle text-2xl" icon="maki:cross" />
				</button>
				<div></div>
				<button
					class="w-16 h-16 text-xl border-2 border-gray-300 rounded-full"
					on:click={() => {
						callback(true);
					}}
				>
					<iconify-icon class="text-green align-middle text-2xl" icon="el:ok" />
				</button>
			</div>
		</div>
	</div>
</div>
