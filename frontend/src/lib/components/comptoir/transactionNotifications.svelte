<script context="module" lang="ts">
	export type TransactionNotification = {
		kind: 'success' | 'error';
		message: string;
		customerName: string;
		customerNickname?: string;
		description?: string;
	};
</script>

<script lang="ts">
	import { onDestroy } from 'svelte';

	const duration = 2500;
	type Notice = TransactionNotification & { id: number };
	let notices: Notice[] = [];
	let nextId = 0;
	const timers = new Set<ReturnType<typeof setTimeout>>();

	export function show(notification: TransactionNotification) {
		const id = ++nextId;
		notices = [...notices, { ...notification, id }];
		const timer = setTimeout(() => {
			timers.delete(timer);
			notices = notices.filter((notice) => notice.id !== id);
		}, duration);
		timers.add(timer);
	}

	onDestroy(() => {
		for (const timer of timers) clearTimeout(timer);
	});
</script>

<div class="transaction-notifications" role="status" aria-live="polite" aria-relevant="additions">
	{#each notices as notice (notice.id)}
		<div
			class="transaction-notice"
			class:error={notice.kind === 'error'}
			style:--notice-duration={`${duration}ms`}
		>
			<div class="notice-icon" aria-hidden="true">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					{#if notice.kind === 'success'}
						<path d="m6 12 4 4 8-8" />
					{:else}
						<path d="M12 8v5m0 3h.01" />
						<circle cx="12" cy="12" r="9" />
					{/if}
				</svg>
			</div>
			<div class="notice-content">
				<p class="notice-title">{notice.message}</p>
				<p class="notice-customer">
					<strong>{notice.customerNickname || notice.customerName}</strong>
					{#if notice.customerNickname && notice.customerName && notice.customerNickname !== notice.customerName}
						<span>({notice.customerName})</span>
					{/if}
				</p>
				{#if notice.description}
					<p class="notice-description">{notice.description}</p>
				{/if}
			</div>
			<div class="notice-progress" aria-hidden="true" />
		</div>
	{/each}
</div>

<style>
	.transaction-notifications {
		position: fixed;
		top: max(16px, env(safe-area-inset-top));
		right: max(16px, env(safe-area-inset-right));
		z-index: 100;
		display: flex;
		flex-direction: column;
		gap: 10px;
		width: 360px;
		max-width: calc(100vw - 32px);
		pointer-events: none;
	}

	.transaction-notice {
		--notice-accent: #15803d;
		--notice-icon-bg: #dcfce7;
		position: relative;
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 16px 18px;
		overflow: hidden;
		border: 1px solid #e5e7eb;
		border-radius: 14px;
		background: #fff;
		color: #111827;
		box-shadow: 0 8px 24px -6px rgb(15 23 42 / 18%), 0 2px 6px rgb(15 23 42 / 4%);
	}

	.transaction-notice.error {
		--notice-accent: #dc2626;
		--notice-icon-bg: #fee2e2;
	}

	.notice-icon {
		display: grid;
		place-items: center;
		flex-shrink: 0;
		width: 36px;
		height: 36px;
		border-radius: 50%;
		background: var(--notice-icon-bg);
		color: var(--notice-accent);
	}

	.notice-icon svg {
		width: 20px;
		height: 20px;
	}

	.notice-content {
		min-width: 0;
		overflow-wrap: anywhere;
		line-height: 1.45;
	}

	p {
		margin: 0;
	}

	.notice-title {
		font-size: 14px;
		font-weight: 650;
	}

	.notice-customer {
		margin-top: 3px;
		font-size: 13px;
		color: #6b7280;
	}

	.notice-customer strong {
		font-weight: 600;
		color: #374151;
	}

	.notice-description {
		margin-top: 5px;
		font-size: 12px;
		color: #6b7280;
	}

	.notice-progress {
		position: absolute;
		bottom: 0;
		left: 0;
		height: 3px;
		width: 100%;
		background: var(--notice-accent);
		opacity: 0.65;
		transform-origin: left;
		animation: notice-countdown var(--notice-duration) linear forwards;
	}

	@keyframes notice-countdown {
		to { transform: scaleX(0); }
	}

	@media (prefers-color-scheme: dark) {
		.transaction-notice {
			--notice-accent: #4ade80;
			--notice-icon-bg: #14382a;
			background: #1f2937;
			border-color: #374151;
			color: #f9fafb;
		}

		.transaction-notice.error {
			--notice-accent: #f87171;
			--notice-icon-bg: #48252a;
		}

		.notice-customer,
		.notice-description {
			color: #9ca3af;
		}

		.notice-customer strong {
			color: #e5e7eb;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.notice-progress {
			animation: none;
		}
	}
</style>
