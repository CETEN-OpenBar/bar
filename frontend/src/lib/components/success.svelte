<script lang="ts">
	export let message: string;
</script>

<!-- Display a popup that asks for a pin -->

<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-40"
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<div class="flex flex-row items-center bg-green-400 rounded-lg shadow-lg p-6 z-40">
		<!-- success icon -->
		<svg class="w-10 h-10 text-white mb-10 mt-10 mr-3" viewBox="0 0 24 24">
			<path
				fill="currentColor"
				d="M12 2C6.486 2 2 6.486 2 12s4.486 10 10 10 10-4.486 10-10S17.514 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"
			/>
		</svg>
		<h1 class="text-2xl font-bold text-white">{message}</h1>
	</div>
</div>
