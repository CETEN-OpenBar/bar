<div class="fixed inset-0 flex items-center justify-center z-50">
    <div class="flex flex-col">
        <h1 class="text-2xl font-bold mb-4 animate-pulse">Chargement du bar</h1>
        <div class="self-center animate-spin rounded-full h-16 w-16 border-t-4 border-blue-500" />
    </div>
</div>

<style>
	/* Style for the loading animation */
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.animate-spin {
		animation: spin 1s linear infinite;
	}
</style>
