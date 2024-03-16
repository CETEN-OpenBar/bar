<script lang="ts">
	import type { CarouselText, CarouselTextCreate } from '$lib/api';
	import { carouselApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import ConfirmationPopup from '$lib/components/confirmationPopup.svelte';

	let allCarouselTexts: CarouselText[] = [];
	let carouselTexts: CarouselText[] = [];
	let newText: CarouselTextCreate = {
		text: '',
		color: ''
	};

	let searchQuery = '';
	let page = 0;
	let textsPerPage = 10;

	let showConfirmation: boolean = false;
	let confirmationCallback: VoidFunction = () => {}

	onMount(() => {
		carouselApi()
			.getCarouselTexts({ withCredentials: true })
			.then((res) => {
				carouselTexts = res.data;
				allCarouselTexts = res.data;
			});
	});

	function createNewCarouselText() {
		carouselApi()
			.addCarouselText(newText, { withCredentials: true })
			.then((res) => {
				carouselTexts = [...carouselTexts, res.data];
				allCarouselTexts = [...allCarouselTexts, res.data];
			});
	}

	function deleteCarouselText(id: string) {
		carouselApi()
			.markDeleteCarouselText(id, { withCredentials: true })
			.then(() => {
				carouselTexts = carouselTexts.filter((ct) => ct.id !== id);
				allCarouselTexts = allCarouselTexts.filter((ct) => ct.id !== id);
			});
	}
</script>

<!-- Popup -->
<div
	id="hs-modal-new-text"
	class="hs-overlay hidden w-full h-full fixed top-0 left-0 z-[60] overflow-x-hidden overflow-y-auto"
>
	<div
		class="hs-overlay-open:mt-7 hs-overlay-open:opacity-100 hs-overlay-open:duration-500 mt-0 opacity-0 ease-out transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto"
	>
		<div
			class="bg-white border border-gray-200 rounded-xl shadow-sm dark:bg-gray-800 dark:border-gray-700"
		>
			<div class="p-4 sm:p-7">
				<div class="text-center">
					<h2 class="block text-2xl font-bold text-gray-800 dark:text-gray-200">
						Ajouter un texte
					</h2>
				</div>

				<div class="mt-5">
					<!-- Form -->
					<div class="grid gap-y-4">
						<!-- Form Group -->
						<div>
							<label for="text" class="block text-sm mb-2 dark:text-white">Texte</label>
							<div class="relative">
								<input
									type="text"
									id="text"
									name="text"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newText.text}
								/>

								<!-- select for color with basic css colors -->
								<label for="color" class="block text-sm mb-2 dark:text-white">Couleur</label>
								<select
									id="color"
									name="color"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="color-error"
									bind:value={newText.color}
								>
									<option value="black">Noir</option>
									<option value="white">Blanc</option>
									<option value="red">Rouge</option>
									<option value="blue">Bleu</option>
									<option value="green">Vert</option>
									<option value="yellow">Jaune</option>
									<option value="orange">Orange</option>
									<option value="purple">Violet</option>
									<option value="pink">Rose</option>
									<option value="brown">Marron</option>
									<option value="gray">Gris</option>
								</select>
							</div>
							<!-- End Form Group -->

							<button
								type="submit"
								class="mt-4 py-3 px-4 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
								on:click={() => createNewCarouselText()}
								data-hs-overlay="#hs-modal-new-text">Créer</button
							>
						</div>
					</div>
					<!-- End Form -->
				</div>
			</div>
		</div>
	</div>
</div>

{#if showConfirmation}
	<ConfirmationPopup
		message="Supprimer ce texte ?"
		confirm_text="Supprimer"
		cancel_callback={() => {
			showConfirmation = false;
		}} 
		confirm_callback={confirmationCallback}
	/>
{/if}

<!-- Table Section -->
<div class="max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto">
	<!-- Card -->
	<div class="flex flex-col">
		<div class="-m-1.5 overflow-x-auto">
			<div class="p-1.5 min-w-full inline-block align-middle">
				<div
					class="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden dark:bg-slate-900 dark:border-gray-700"
				>
					<!-- Header -->
					<div
						class="px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-b border-gray-200 dark:border-gray-700"
					>
						<div>
							<h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Textes</h2>
							<p class="text-sm text-gray-600 dark:text-gray-400">Ajouter des textes au carousel</p>
						</div>



						<!-- search bar -->
						<div class="relative mt-4 w-96 md:mt-0">
							<input
								type="text"
								class="py-3 px-4 w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								placeholder="Rechercher"
								aria-label="Rechercher"
								on:input={(e) => {
									// @ts-ignore
									searchQuery = e.target.value.toLowerCase();
									
									// filter the carousel texts
									carouselTexts = allCarouselTexts.filter((text) => {
										return text.text.toLowerCase().includes(searchQuery);
									});
								}}
							/>
							<svg
								class="absolute w-4 h-4 right-3 top-3 text-gray-400 dark:text-gray-300 pointer-events-none"
								xmlns="http://www.w3.org/2000/svg"
								width="16"
								height="16"
								viewBox="0 0 16 16"
								fill="none"
							>
								<path
									d="M11.6667 11.6667L15.3333 15.3333"
									stroke="currentColor"
									stroke-width="1.5"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
								<path
									d="M6.66663 12.6667C9.53763 12.6667 12 10.2037 12 7.33337C12 4.46337 9.53763 2.00004 6.66663 2.00004C3.79563 2.00004 1.33329 4.46337 1.33329 7.33337C1.33329 10.2037 3.79563 12.6667 6.66663 12.6667Z"
									stroke="currentColor"
									stroke-width="1.5"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
							</svg>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
									data-hs-overlay="#hs-modal-new-text"
								>
									<svg
										class="w-3 h-3"
										xmlns="http://www.w3.org/2000/svg"
										width="16"
										height="16"
										viewBox="0 0 16 16"
										fill="none"
									>
										<path
											d="M2.63452 7.50001L13.6345 7.5M8.13452 13V2"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
										/>
									</svg>
									Ajouter un texte
								</button>
							</div>
						</div>
					</div>
					<!-- End Header -->

					<!-- Table -->
					<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
						<thead class="bg-gray-50 dark:bg-slate-800">
							<tr>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Texte
										</span>
									</div>
								</th>

								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Couleur
										</span>
									</div>
								</th>

								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each carouselTexts.slice(page * textsPerPage, (page + 1) * textsPerPage) as carouselText}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p class="block text-sm text-gray-500 break-words">{carouselText.text}</p>
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-3">
											<p
												class="inline-flex items-center gap-1.5 py-2.5 px-5 rounded-full text-xs font-medium border-2 border-black"
												style="background-color: {carouselText.color};"
											/>
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-1.5">
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => {
													confirmationCallback = () => {
														showConfirmation = false;
														deleteCarouselText(carouselText.id);
													}
													showConfirmation = true;
												}}
											>
												Supprimer
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
					<!-- End Table -->

					<!-- Footer -->
					<div
						class="px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-t border-gray-200 dark:border-gray-700"
					>
						<div>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								<span class="font-semibold text-gray-800 dark:text-gray-200"
									>{carouselTexts.length}</span
								> résultats
							</p>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page > 0) page--;
									}}
								>
									<svg
										class="w-3 h-3"
										xmlns="http://www.w3.org/2000/svg"
										width="16"
										height="16"
										fill="currentColor"
										viewBox="0 0 16 16"
									>
										<path
											fill-rule="evenodd"
											d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"
										/>
									</svg>
									Précédent
								</button>

								<p class="text-sm self-center text-gray-600 dark:text-gray-400">
									Page {page + 1} / {Math.ceil(carouselTexts.length / textsPerPage)}
								</p>

								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page < Math.ceil(carouselTexts.length / textsPerPage) - 1) page++;
									}}
								>
									Suivant
									<svg
										class="w-3 h-3"
										xmlns="http://www.w3.org/2000/svg"
										width="16"
										height="16"
										fill="currentColor"
										viewBox="0 0 16 16"
									>
										<path
											fill-rule="evenodd"
											d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"
										/>
									</svg>
								</button>
							</div>
						</div>
					</div>
					<!-- End Footer -->
				</div>
			</div>
		</div>
	</div>
	<!-- End Card -->
</div>
<!-- End Table Section -->
