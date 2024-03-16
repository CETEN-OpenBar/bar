<script lang="ts">
	export let message: string = "Voulez-vous vraiment faire ceci ?";
    export let confirm_text: string = "Confirmer";
    export let cancel_text: string = "Annuler";
    export let cancel_callback: VoidFunction;
    export let confirm_callback: VoidFunction;
</script>

<!-- Display a popup that asks to confirm an action -->

<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-40"
    on:click={cancel_callback}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<div class="flex flex-col items-center bg-slate-700 rounded-lg shadow-lg p-6 z-40 space-y-5 max-w-[30vw]">
		<h1 class="text-2xl font-bold text-white">{message}</h1>
        <div class="flex flex-row w-full justify-end space-x-5">
            <button on:click={cancel_callback} class="bg-slate-400 text-black rounded-md p-4">{cancel_text}</button>
            <button on:click={confirm_callback} class="bg-red-500 text-black rounded-md p-4">{confirm_text}</button>
        </div>
	</div>
</div>