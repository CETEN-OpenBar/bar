<script lang="ts">
	import Refills from '$lib/components/comptoir/refills.svelte';
	import ComptoirHeaderControls from '$lib/components/comptoir/headerControls.svelte';
</script>

<div class="content-wrapper">
	<ComptoirHeaderControls activeTab="refills" />

	<main class="content-area">
		<Refills amount={6} />
	</main>
</div>

<style>
	.content-wrapper {
		padding: 15px 15px 0 15px;
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.content-area {
		flex: 1;
		background-color: var(--bg-secondary);
		border-radius: 0.75rem;
		border: 1px solid var(--border-color);
		overflow: hidden;
	}

	:root {
		--bg-secondary: #1a1a1a;
		--border-color: #404040;
	}

	@media (prefers-color-scheme: light) {
		:root {
			--bg-secondary: #f8fafc;
			--border-color: #e2e8f0;
		}
	}
</style>