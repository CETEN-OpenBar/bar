<script lang="ts">
	import { fly } from 'svelte/transition';
	import QrOnboarding from '$lib/components/borne/qr_onboarding.svelte';
	import { onMount } from 'svelte';
	import { elasticIn } from 'svelte/easing';
</script>

<div class="gap-16 p-16 h-screen w-screen bg-slate-700">
	<div class="justify-center absolute top-[50%] left-[50%] -translate-x-[50%] -translate-y-[50%]">
		<div class="text-center text-4xl text-white font-bold mb-10">
			Connectez-vous avec tn.net pour continuer la création du compte.
		</div>
		<div class="w-full flex justify-center">
			<div class="w-96 h-96">

				<QrOnboarding />
			</div>
		</div>
	</div>
</div>
