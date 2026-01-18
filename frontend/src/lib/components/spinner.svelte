<script lang="ts">
    export let size = "2rem";
</script>

<div class="text-center">
    <div
        class="animate-spin rounded-full border-t-2 border-b-2 border-blue-500 mx-auto"
        style="width: {size}; height: {size}"
    />
</div>

<style>
	.animate-spin {
		animation: spin 1s linear infinite;
	}
	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
</style>