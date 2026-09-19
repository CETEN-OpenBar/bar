<script lang="ts">
	import { onDestroy } from 'svelte';

	type Notice = { id: number; message: string; kind: 'success' | 'error' };
	let notices: Notice[] = [];
	let nextId = 0;
	const timers = new Set<ReturnType<typeof setTimeout>>();

	export function show(message: string, kind: Notice['kind']) {
		const id = ++nextId;
		notices = [...notices, { id, message, kind }];
		const timer = setTimeout(() => {
			timers.delete(timer);
			notices = notices.filter((notice) => notice.id !== id);
		}, 1500);
		timers.add(timer);
	}

	onDestroy(() => {
		for (const timer of timers) clearTimeout(timer);
	});
</script>

<div class="transaction-notifications" role="status" aria-live="polite" aria-relevant="additions">
	{#each notices as notice (notice.id)}
		<div class="transaction-notice" class:error={notice.kind === 'error'}>
			<iconify-icon
				icon={notice.kind === 'success' ? 'mdi:check-circle' : 'mdi:alert-circle'}
				aria-hidden="true"
			/>
			<span>{notice.message}</span>
		</div>
	{/each}
</div>

<style>
	.transaction-notifications {
		position: fixed;
		top: max(12px, env(safe-area-inset-top));
		left: 50%;
		transform: translateX(-50%);
		z-index: 100;
		display: flex;
		flex-direction: column;
		gap: 8px;
		width: max-content;
		max-width: min(560px, calc(100vw - 24px));
		pointer-events: none;
	}

	.transaction-notice {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 12px 16px;
		border: 1px solid #86efac;
		border-radius: 12px;
		background: #f0fdf4;
		color: #14532d;
		box-shadow: 0 4px 16px rgb(0 0 0 / 12%);
		font-size: 15px;
		font-weight: 600;
		line-height: 1.4;
		overflow-wrap: anywhere;
	}

	.transaction-notice.error {
		border-color: #fca5a5;
		background: #fef2f2;
		color: #991b1b;
	}

	iconify-icon {
		flex-shrink: 0;
		font-size: 22px;
	}
</style>
