<script lang="ts">
	export let text: string;

	let clazz = '';
	export {clazz as class};

	export let callback: () => void;
</script>

<button class="inline-block px-4 py-2 rounded hover:cursor-pointer {clazz}" type="button" on:click={callback}>
	{text}
</button>
