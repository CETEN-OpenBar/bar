
<div class="fixed inset-0 flex items-center justify-center z-50">
    <div class="flex flex-col roll-in-bottom">
        <h1 class="text-2xl font-bold mb-4">Chargement du bar...</h1>
        <div class="self-center animate-spin rounded-full h-16 w-16 border-t-4 border-blue-500" />
    </div>
</div>

<style>


	/* ----------------------------------------------
 * Generated by Animista on 2023-8-12 0:34:45
 * Licensed under FreeBSD License.
 * See http://animista.net/license for more info. 
 * w: http://animista.net, t: @cssanimista
 * ---------------------------------------------- */

	@-webkit-keyframes roll-in-bottom {
		0% {
			-webkit-transform: translateY(800px) rotate(540deg);
			transform: translateY(800px) rotate(540deg);
			opacity: 0;
		}
		100% {
			-webkit-transform: translateY(0) rotate(0deg);
			transform: translateY(0) rotate(0deg);
			opacity: 1;
		}
	}
	@keyframes roll-in-bottom {
		0% {
			-webkit-transform: translateY(800px) rotate(540deg);
			transform: translateY(800px) rotate(540deg);
			opacity: 0;
		}
		100% {
			-webkit-transform: translateY(0) rotate(0deg);
			transform: translateY(0) rotate(0deg);
			opacity: 1;
		}
	}

	.roll-in-bottom {
		-webkit-animation: roll-in-bottom 1s ease-out both;
		animation: roll-in-bottom 1s ease-out both;
	}
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
