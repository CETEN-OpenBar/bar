<script lang="ts">
	import { refillsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onDestroy, onMount } from 'svelte';
	import { RefillType, type Refill } from '$lib/api';
	import { dragscroll } from '@svelte-put/dragscroll';
	import ComptoirHeaderControls from '$lib/components/comptoir/headerControls.svelte';

	let refills: Array<Refill> = [];
	let refillAmount: number = 4;
	let interval: number;

	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadRefills();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadRefills();
		}
	};

	function handleAvatarError(e: Event) {
		const target = e.currentTarget as HTMLImageElement;
		target.style.display = 'none';
		target.nextElementSibling?.classList.remove('hidden');
	}

	onMount(() => {
		reloadRefills();
		interval = setInterval(() => {
			reloadRefills();
		}, 2000);

		window.addEventListener('keydown', handleKeyDown);

		const refillsList = document.querySelector('.refills-list');
		if (refillsList) {
			let touchStartX = 0;
			let touchEndX = 0;
			let isTouch = false;

			refillsList.addEventListener(
				'touchstart',
				(e: Event) => {
					const touchEvent = e as TouchEvent;
					touchStartX = touchEvent.changedTouches[0].screenX;
					isTouch = true;
				},
				{ passive: true }
			);

			refillsList.addEventListener(
				'touchend',
				(e: Event) => {
					if (isTouch) {
						const touchEvent = e as TouchEvent;
						touchEndX = touchEvent.changedTouches[0].screenX;
						handleSwipe();
						isTouch = false;
					}
				},
				{ passive: true }
			);

			refillsList.addEventListener('mousedown', (e: Event) => {
				const mouseEvent = e as MouseEvent;
				if (!isTouch) {
					touchStartX = mouseEvent.screenX;
					isTouch = false;
				}
			});

			refillsList.addEventListener('mouseup', (e: Event) => {
				const mouseEvent = e as MouseEvent;
				if (!isTouch) {
					touchEndX = mouseEvent.screenX;
					handleMouseSwipe();
				}
			});

			const handleSwipe = () => {
				const swipeThreshold = 50;
				const diff = touchStartX - touchEndX;

				if (Math.abs(diff) > swipeThreshold) {
					if (diff > 0) {
						if (page < maxPage) nextPage();
					} else {
						if (page > 1) prevPage();
					}
				}
			};

			const handleMouseSwipe = () => {
				const swipeThreshold = 100;
				const diff = touchStartX - touchEndX;

				if (Math.abs(diff) > swipeThreshold) {
					if (diff > 0) {
						if (page < maxPage) nextPage();
					} else {
						if (page > 1) prevPage();
					}
				}
			};
		}
	});

	onDestroy(() => {
		clearInterval(interval);
		if (typeof window !== 'undefined') {
			window.removeEventListener('keydown', handleKeyDown);
		}
	});

	function reloadRefills() {
		refillsApi()
			.getRefills(page, refillAmount, undefined, undefined, { withCredentials: true })
			.then((res) => {
				page = res.data.page ?? 0;
				maxPage = res.data.max_page ?? 0;
				refills = res.data.refills ?? [];
			});
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (
			event.target instanceof HTMLInputElement ||
			event.target instanceof HTMLTextAreaElement ||
			event.target instanceof HTMLSelectElement
		) {
			return;
		}

		switch (event.key) {
			case 'ArrowLeft':
				event.preventDefault();
				if (page > 0) prevPage();
				break;
			case 'ArrowRight':
				event.preventDefault();
				if (page < maxPage) nextPage();
				break;
			case 'PageUp':
				event.preventDefault();
				if (page > 0) prevPage();
				break;
			case 'PageDown':
				event.preventDefault();
				if (page < maxPage) nextPage();
				break;
			case 'Home':
				event.preventDefault();
				if (page > 1) {
					page = 1;
					reloadRefills();
				}
				break;
			case 'End':
				event.preventDefault();
				if (page < maxPage) {
					page = maxPage;
					reloadRefills();
				}
				break;
		}
	}

	function handleCancelRefill(refill: Refill) {
		refillsApi()
			.patchRefillId(refill.account_id, refill.id, 'canceled', refill.type, {
				withCredentials: true
			})
			.then(() => {
				reloadRefills();
			});
	}

	function handleReValidateRefill(refill: Refill) {
		refillsApi()
			.patchRefillId(refill.account_id, refill.id, 'valid', refill.type, { withCredentials: true })
			.then(() => {
				reloadRefills();
			});
	}
</script>

<div class="refills-wrapper">
	<div class="refills-content">
		<ComptoirHeaderControls activeTab="refills" showSearch={false} />

		<div use:dragscroll class="refills-list">
			{#each refills as refill}
				<div class="refill-card {refill.state}">
					<div class="refill-header">
						{#if refill.account_google_picture}
							<img
								src={refill.account_google_picture}
								alt="Avatar"
								class="refill-avatar"
								on:error={handleAvatarError}
							/>
						{/if}
						<iconify-icon
							icon="mdi:account-circle"
							class="refill-avatar placeholder-icon {refill.account_google_picture
								? 'hidden'
								: ''}"
						/>
						<div class="refill-user-info">
							<b>{refill.account_name}</b>
							<span class="refill-date">{new Date(refill.issued_at * 1000).toLocaleString()}</span>
						</div>
						{#if refill.type == RefillType.RefillHelloAsso}
							<iconify-icon icon="mdi:credit-card-outline" class="type-icon" />
						{/if}
					</div>

					<div class="refill-content">
						<div class="refill-info">
							<div class="refill-row">
								<span class="refill-label">De :</span>
								<span class="refill-value">{refill.issued_by_name}</span>
							</div>
							{#if refill.canceled_by}
								<div class="refill-row canceled">
									<span class="refill-label">Annulé par :</span>
									<span class="refill-value">{refill.canceled_by_name}</span>
								</div>
							{/if}
						</div>
						<div class="refill-amount-section">
							<span class="refill-amount">{formatPrice(refill.amount)}</span>
						</div>
					</div>

					<div class="refill-actions">
						{#if refill.type != RefillType.RefillHelloAsso}
							{#if refill.state == 'valid'}
								<button class="action-button cancel" on:click={() => handleCancelRefill(refill)}>
									<iconify-icon icon="mdi:close-circle" />
									<span>Annuler</span>
								</button>
							{:else if refill.state == 'canceled'}
								<button class="action-button validate" on:click={() => handleReValidateRefill(refill)}>
									<iconify-icon icon="mdi:check-circle" />
									<span>Re-valider</span>
								</button>
							{/if}
                        {:else}
                            <span>Recharge HelloAsso - Pas d'actions possibles</span>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</div>

	<div class="pagination">
		<div class="pagination-results">
			<span class="font-semibold">{refills.length}</span> résultats
		</div>
		<div class="pagination-controls">
			<p class="pagination-info">
				Page <span class="font-bold">{page}</span> sur <span class="font-bold">{maxPage}</span>
			</p>
			<div class="pagination-buttons">
				<button class="pagination-button" on:click={prevPage} disabled={page === 1}>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M15 19l-7-7 7-7"
						/>
					</svg>
					<span>Précédent</span>
				</button>
				<button class="pagination-button" on:click={nextPage} disabled={page === maxPage}>
					<span>Suivant</span>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 5l7 7-7 7"
						/>
					</svg>
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	.refills-wrapper {
		width: 100%;
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.refills-content {
		display: flex;
		flex-direction: column;
		flex-grow: 1;
		flex: 1;
		padding: 15px;
		padding-bottom: 0;
	}

	.refills-list {
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		flex-grow: 1;
		gap: 20px;
		max-height: calc(100vh - 260px);
	}

	.refill-card {
		display: flex;
		flex-direction: column;
		border: 2px solid #e5e7eb;
		border-radius: 12px;
		transition: all 0.2s ease;
		position: relative;
		background-color: white;
		flex-shrink: 0;
	}

	@media (prefers-color-scheme: dark) {
		.refill-card {
			background-color: #2f2f41;
			border: 2px solid #3c4554;
		}
	}

	.refill-card:hover {
		border: 2px solid #9ca3af;
	}

	.refill-card.canceled {
		border-color: #ef4444;
	}

	.refill-card.valid {
		border-color: #22c55e;
	}

	.refill-card.canceled .refill-actions {
		opacity: 0.5;
	}

	.refill-header {
		color: black;
		padding: 12px 12px 12px 12px;
		display: flex;
		align-items: center;
		gap: 12px;
		position: relative;
		border-bottom: 1px solid #9ca3af;
	}

	@media (prefers-color-scheme: dark) {
		.refill-header {
			color: white;
			border-bottom-color: #4b5563;
		}
	}

	.refill-avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		object-fit: cover;
		flex-shrink: 0;
	}

	.placeholder-icon {
		color: #9ca3af;
		font-size: 40px;
		flex-shrink: 0;
	}

	@media (prefers-color-scheme: dark) {
		.placeholder-icon {
			color: #6b7280;
		}
	}

	.refill-user-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.refill-date {
		font-size: 12px;
		color: #6b7280;
	}

	@media (prefers-color-scheme: dark) {
		.refill-date {
			color: #9ca3af;
		}
	}

	.type-icon {
		height: 20px;
		width: 20px;
		margin-left: auto;
		color: #374151;
	}

	@media (prefers-color-scheme: dark) {
		.type-icon {
			color: #e5e7eb;
		}
	}

	.refill-content {
		display: grid;
		grid-template-columns: 10fr 1fr;
		gap: 8px;
	}

	.refill-info {
		padding: 8px;
		padding-right: 8px;
		height: 100%;
		width: 100%;
		border-right: 1px solid #9ca3af;
	}

	@media (prefers-color-scheme: dark) {
		.refill-info {
			border-right-color: #4b5563;
		}
	}

	.refill-row {
		display: flex;
		gap: 8px;
		font-size: 14px;
		color: black;
	}

	@media (prefers-color-scheme: dark) {
		.refill-row {
			color: white;
		}
	}

	.refill-row.canceled {
		color: #ef4444;
	}

	.refill-label {
		font-weight: 500;
	}

	.refill-value {
		color: #374151;
	}

	@media (prefers-color-scheme: dark) {
		.refill-value {
			color: #e5e7eb;
		}
	}

	.refill-amount-section {
		padding: 8px;
		padding-left: 8px;
		width: 100%;
		font-size: 18px;
		font-weight: 600;
		align-self: center;
		text-align: center;
		color: black;
	}

	@media (prefers-color-scheme: dark) {
		.refill-amount-section {
			color: white;
		}
	}

	.refill-actions {
		display: flex;
		justify-content: flex-end;
		padding: 8px 12px;
		border-top: 1px solid #9ca3af;
	}

	@media (prefers-color-scheme: dark) {
		.refill-actions {
			border-top-color: #4b5563;
		}
	}

	.action-button {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 16px;
		border-radius: 8px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		border: none;
	}

	.action-button.cancel {
		background-color: #ef4444;
		color: white;
	}

	.action-button.cancel:hover {
		background-color: #dc2626;
	}

	.action-button.validate {
		background-color: #22c55e;
		color: white;
	}

	.action-button.validate:hover {
		background-color: #16a34a;
	}

	.pagination {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 16px;
		background-color: white;
		border-top: 1px solid #e5e7eb;
	}

	.pagination-results {
		font-size: 14px;
		color: #6b7280;
	}

	.pagination-controls {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.pagination-info {
		font-size: 14px;
		font-weight: 500;
		color: #374151;
	}

	.pagination-buttons {
		display: flex;
		gap: 8px;
	}

	.pagination-button {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 16px;
		background-color: white;
		color: #374151;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.pagination-button:hover:not(:disabled) {
		background-color: #f9fafb;
		border-color: #9ca3af;
	}

	.pagination-button:disabled {
		background-color: #f9fafb;
		color: #9ca3af;
		cursor: not-allowed;
	}

	@media (prefers-color-scheme: dark) {
		.pagination {
			background-color: #1f2937;
			border-color: #374151;
		}

		.pagination-results {
			color: #9ca3af;
		}

		.pagination-info {
			color: #e5e7eb;
		}

		.pagination-button {
			background-color: #374151;
			color: #e5e7eb;
			border-color: #4b5563;
		}

		.pagination-button:hover:not(:disabled) {
			background-color: #4b5563;
			border-color: #6b7280;
		}

		.pagination-button:disabled {
			background-color: #374151;
			color: #6b7280;
		}
	}

	@media (max-width: 768px) {
		.refill-content {
			grid-template-columns: 8fr 1fr;
		}

		.refill-row {
			font-size: 12px;
		}

		.refill-actions {
			justify-content: center;
		}

		.action-button {
			padding: 6px 12px;
			font-size: 12px;
		}
	}

	@media (max-width: 480px) {
		.refill-content {
			grid-template-columns: 1fr;
		}

		.refill-info {
			border-right: none;
			border-bottom: 1px solid #9ca3af;
		}

		@media (prefers-color-scheme: dark) {
			.refill-info {
				border-bottom-color: #4b5563;
			}
		}

		.refill-amount-section {
			text-align: center;
			padding-top: 8px;
		}

		.refill-header {
			font-size: 14px;
		}

		.refill-avatar {
			width: 32px;
			height: 32px;
			font-size: 32px;
		}

		.placeholder-icon {
			font-size: 32px;
		}

		.pagination {
			flex-direction: column;
			gap: 12px;
		}

		.pagination-controls {
			flex-direction: column;
			gap: 5px;
		}

		.refills-list {
			max-height: calc(100vh - 340px);
		}
	}
</style>