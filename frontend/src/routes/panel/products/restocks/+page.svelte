<script lang="ts">
	import {		
		RestockState,
		type Item,
		type NewRestock,
		type NewRestockItem,
		type Restock
	} from '$lib/api';
	import type {
		Fournisseur
	} from '$lib/api';
	import ConfirmationPopup from '$lib/components/confirmationPopup.svelte';
	import { itemsApi, restocksApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice, restockTypeIterator } from '$lib/utils';
	import { onMount } from 'svelte';
	import Time from 'svelte-time';
	import { getDocument, GlobalWorkerOptions } from 'pdfjs-dist';
	import pdfjsWorker from "pdfjs-dist/build/pdf.worker?url";

	let sure: boolean = false;
	let items: Item[] = [];
	// Array of items which are unkown by their reference for the automatic restock with the invoice
	let unknownItems: Array<{name : string, reference : number}> = [];

	let restocks: Restock[] = [];
	let newRestock: NewRestock = {
		total_cost_ht: 0,
		total_cost_ttc: 0,
		driver_id: '',
		type: 'promocash',
		state: RestockState.RestockPending,
		items: []
	};

	let page: number = 0;
	let max_page: number = 0;
	let maxPageRestock: number = 0;
	let pageRestock: number = 0;
	let itemsPerPage: number = 4;

	let nextPage = () => {
		if (pageRestock < maxPageRestock) {
			pageRestock++;
			reloadRestocks();
		}
	};
	let prevPage = () => {
		if (pageRestock > 0) {
			pageRestock--;
			reloadRestocks();
		}
	};

	let newItem: NewRestockItem = {
		item_id: '',
		item_name: '',
		amount_of_bundle: 0,
		amount_per_bundle: 0,
		bundle_cost_ht: 0,
		bundle_cost_ttc: 0,
		bundle_cost_float_ttc: 0,
		tva: 0
	};
	let searchName: string = '';	

	type dV = {
		name: string;
		item_price_calc: number;
		item_price: string;
		item_price_ht: string;
		amount_of_bundle: string;
		amount_per_bundle: string;
		bundle_cost_ht: string;
		tva: string;
		bundle_cost_ttc: string;
	};

	let displayedValues: dV = {
		name: 'Nom du produit',
		item_price_calc: 0,
		item_price: 'Prix coûtant TTC',
		item_price_ht: 'Prix coûtant HT',
		amount_of_bundle: 'Nombre de lots',
		amount_per_bundle: 'Nombre de produits par lots',
		bundle_cost_ht: "Prix d'un lot HT",
		tva: '0',
		bundle_cost_ttc: "Prix d'un lot TTC"
	};

	let deletingRestock = false;
	let deleteRestockCallback: VoidFunction = () => {};
	let confirmationMessage: string | undefined = undefined;

	let selectedViewRestock: Restock | undefined = undefined;
	let selectedEditRestock: Restock | undefined = undefined;

	onMount(() => {
		reloadItems();
		reloadRestocks();
	});

	function reloadRestocks() {
		restocksApi()
			.getRestocks(pageRestock, itemsPerPage, undefined, {
				withCredentials: true
			})
			.then((res) => {
				restocks = res.data.restocks ?? [];
				pageRestock = res.data.page ?? 0;
				maxPageRestock = res.data.max_page ?? 0;
			});
	}

	function reloadItems() {
		itemsApi()
			.getAllItems(page, itemsPerPage, undefined, undefined, searchName, undefined, {
				withCredentials: true
			})
			.then((res) => {
				max_page = res.data.max_page ?? 0;
				page = res.data.page ?? 0;
				itemsPerPage = res.data.limit ?? 0;
				items = res.data.items ?? [];
			});
	}

	function updateTotalHTandTTC() {
		newRestock.total_cost_ht = 0;
		newRestock.items.forEach((item) => {
			newRestock.total_cost_ht += item.amount_of_bundle * item.bundle_cost_ht;
		});
		newRestock.total_cost_ttc = 0.0;
		newRestock.items.forEach((item) => {
			if (item.bundle_cost_float_ttc === 0.0) {
				newRestock.total_cost_ttc += item.amount_of_bundle * item.bundle_cost_ttc;
			} else {
				newRestock.total_cost_ttc +=
					item.amount_of_bundle * (item.bundle_cost_float_ttc ?? item.bundle_cost_ttc);
			}
		});
	}

	async function applyRestock(state: RestockState) {
		newRestock.driver_id = undefined;
		newRestock.total_cost_ttc = Math.round(newRestock.total_cost_ttc);
		newRestock.total_cost_ht = Math.round(newRestock.total_cost_ht);
		newRestock.state = state;
		restocksApi()
			.createRestock(newRestock, { withCredentials: true })
			.then((res) => {
				restocks = [res.data, ...restocks];
				newRestock = {
					total_cost_ht: 0,
					total_cost_ttc: 0,
					driver_id: '',
					type: newRestock.type,
					state: newRestock.state,
					items: []
				};
				displayedValues = {
					name: 'Nom du produit',
					item_price_calc: 0,
					item_price: 'Prix coûtant TTC',
					item_price_ht: 'Prix coûtant HT',
					amount_of_bundle: 'Nombre de lots',
					amount_per_bundle: 'Nombre de produits par lots',
					bundle_cost_ht: "Prix d'un lot HT",
					tva: '0',
					bundle_cost_ttc: "Prix d'un lot TTC"
				};
				newItem = {
					item_id: '',
					item_name: '',
					amount_of_bundle: 0,
					amount_per_bundle: 0,
					bundle_cost_ht: 0,
					bundle_cost_ttc: 0,
					bundle_cost_float_ttc: 0.0,
					tva: 0
				};
				sure = false;
				reloadRestocks();
			});
	}

	async function applyEditRestock() {
		newRestock.state = RestockState.RestockFinished;
		saveEditRestock().then(() => {
			selectedEditRestock = undefined;
			sure = false;
			reloadRestocks();
		});
	}

	async function saveEditRestock() {
		if (!selectedEditRestock) return;		
		restocksApi()
			.updateRestock(selectedEditRestock.id, newRestock, { withCredentials: true })
			.then((res) => {
				newRestock = {
					total_cost_ht: 0,
					total_cost_ttc: 0,
					driver_id: '',
					type: newRestock.type,
					state: newRestock.state,
					items: []
				};
				reloadRestocks();
			});
	}

	function updatePrices() {
		// Calculate from displayedValues.item_price_calc, displayedValues.amount_of_bundle and TVA
		if (newItem.amount_of_bundle === 0 || newItem.amount_per_bundle === 0) return;

		if (displayedValues.bundle_cost_ht === "Prix d'un lot HT") {
			newItem.bundle_cost_ht = Math.round(
				(displayedValues.item_price_calc * newItem.amount_per_bundle) / (1 + newItem.tva / 10000)
			);
			displayedValues.bundle_cost_ht = formatPrice(newItem.bundle_cost_ht);
		}
		if (displayedValues.bundle_cost_ttc === "Prix d'un lot TTC") {
			newItem.bundle_cost_ttc = Math.round(
				displayedValues.item_price_calc * newItem.amount_per_bundle
			);
			displayedValues.bundle_cost_ttc = formatPrice(newItem.bundle_cost_ttc);
		} else {
			newItem.bundle_cost_ttc = Math.round(newItem.bundle_cost_ht * (1 + newItem.tva / 10000));
			newItem.bundle_cost_float_ttc = newItem.bundle_cost_ht * (1 + newItem.tva / 10000);
			displayedValues.bundle_cost_ttc = formatPrice(newItem.bundle_cost_ttc);
		}
		displayedValues.item_price_ht = formatPrice(newItem.bundle_cost_ht * newItem.amount_of_bundle);
		displayedValues.item_price = formatPrice(newItem.bundle_cost_ttc * newItem.amount_of_bundle);
	}

	function deleteRestock(restockId: string) {
		if (selectedEditRestock && selectedEditRestock.id === restockId) {
			selectedEditRestock = undefined;
		}
		restocksApi()
			.deleteRestock(restockId, { withCredentials: true })
			.then(() => {
				restocks = restocks.filter((ct) => ct.id !== restockId);
				reloadRestocks();
			});
	}

	function calculateTtc(item: NewRestockItem) {
		return Math.round(item.bundle_cost_ht * (1 + item.tva / 10000));
	}
	
	function calculateHt(item: NewRestockItem) {
		return Math.round(item.bundle_cost_ttc / (1 + item.tva / 10000));
	}

	onMount(() => {
    	GlobalWorkerOptions.workerSrc = pdfjsWorker;
  });

  let fileInput: HTMLInputElement;
  let pdfText = '';

  const handleFileChange = async (event: Event) => {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      const fileReader = new FileReader();

      fileReader.onload = async (event) => {
        const arrayBuffer = event.target?.result as ArrayBuffer;
        if (arrayBuffer) {
          const pdf = await getDocument(arrayBuffer).promise;
          let fullText = '';

          for (let i = 1; i <= pdf.numPages; i++) {
            const page = await pdf.getPage(i);
            const textContent = await page.getTextContent();
            const text = textContent.items.map((item: any) => item.str).join(' ');
            fullText += text + ' ';
          }

          pdfText = fullText;          
        }
      };

      fileReader.readAsArrayBuffer(file);
    }
  };

	function isDigit(chr: string): boolean{
		return chr >= '0' && chr <= '9';
	}

	function isPrice(str: string): boolean{
		const len = str.length;
		return len >= 4 && isDigit(str[len-1]) && isDigit(str[len-2]) && str[len-3] == ',' && isDigit(str[len-4]);
	}

	function restockAuchanDrive(input: string): Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> {
		const splitPage = "Référence   Caractéristiques produit   Prix U. (HT) € Remises U. (HT) €   Qte.   Prix total (HT) € Taux TVA % Cagnotte WAAOH! Prix total (TTC) €  "
		const splitLine = "   ";
		const endPage1 = "Votre commande :";
		const endPage2 = "Taux TVA";	
		const endPage3 = "Votre facture :";
		const pages = input.split(splitPage);
		const itemsList: Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> = [];
		for (let i=1;i<pages.length;i++){
			const pageSplit = pages[i].split(splitLine);
			for (let j = 0; j < pageSplit.length; ) {
				if (pageSplit[j].includes(endPage1) || pageSplit[j].includes(endPage2) || pageSplit[j].includes(endPage3)) {				
					break;
				}
				const reference = parseInt(pageSplit[j].split(" ").slice(-1)[0]);
				j++;
				let name;
				let unitPriceHT;
				let splitName = pageSplit[j].split(" ");
				if (isPrice(splitName[splitName.length-1])){
					name = splitName.slice(0, -1).join(" ");
					unitPriceHT = parseFloat(splitName[splitName.length-1].replace(",","."));
				}
				else{
					name = pageSplit[j];
					j++;
					unitPriceHT = parseFloat(pageSplit[j].replace(",","."));
				}
				j++;
				if (isPrice(pageSplit[j])){ // Check if there is a discount
					unitPriceHT -= parseFloat(pageSplit[j].replace(",","."));
					j++;
				}
				const quantity = parseInt(pageSplit[j]);
				j+=2;
				const tvaRate = parseFloat(pageSplit[j].replace(",","."));
				j++;
				if (isPrice(pageSplit[j])){ // Check if there is a cagnotte WHAOOH					
					j++;
				}				
				itemsList.push({
					reference,
					name,
					quantity,
					unitPriceHT,
					tvaRate
				});

			}
		}
		return itemsList;
	}

	function restockAuchan(input: string): Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> {
		const splitPage = "Référence   Caractéristiques produit   Prix U. (HT) € Remises U. (HT) €   Qte.   Prix total (HT) € Taux TVA % Prix total (TTC) €  ";
		const splitLine = "   ";
		const splitNameQuantity = "  Dont éco-participation :  ";
		const endPage1 = "Votre commande :";
		const endPage2 = "TVA déjà collectée  Mode de paiement";
		const endPage3 = "Taux TVA (%)";
		const pages = input.split(splitPage);
		const itemsList: Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> = [];
		for (let i=1;i<pages.length;i++){
			const pageSplit = pages[i].split(splitLine);
			for (let j = 0; j < pageSplit.length;) {
				if (pageSplit[j].includes(endPage1) || pageSplit[j].includes(endPage2) || pageSplit[j].includes(endPage3)) {
					break;
				}
				let reference;
				if (pageSplit[j].includes("  ")){
					reference = parseInt(pageSplit[j].split("  ").slice(-1)[0]);
				}
				else{
					reference = parseInt(pageSplit[j].split(" ").slice(-1)[0]);
				}
				j++;
				const splitEcoPart = pageSplit[j].split(splitNameQuantity);
				const name = splitEcoPart[0];
				let unitPriceHT;
				if (splitEcoPart.length > 1){
					unitPriceHT = parseFloat(splitEcoPart[1].replace(",","."));
				}
				else{
					j++;
					unitPriceHT = parseFloat(pageSplit[j].replace(",","."));
				}
				j++;
				if (isPrice(pageSplit[j])){ // Check if there is a discount
					unitPriceHT -= parseFloat(pageSplit[j].replace(",","."));
					j++;
				}
				const quantity = parseInt(pageSplit[j]);
				j += 2;
				const tvaRate = parseFloat(pageSplit[j].replace(',','.'));
				j++;
				itemsList.push({
					reference,
					name,
					quantity,
					unitPriceHT,
					tvaRate
				});
			}
		}
		return itemsList;
	}

	function add_item_to_restock(): void{
		let t = newRestock.items;
		t.unshift(newItem);
		newRestock.items = t;
		displayedValues = {
			name: 'Nom du produit',
			item_price_calc: 0,
			item_price: 'Prix coûtant TTC',
			item_price_ht: 'Prix coûtant HT',
			amount_of_bundle: 'Nombre de lots',
			amount_per_bundle: 'Nombre de produits par lots',
			bundle_cost_ht: "Prix d'un lot HT",
			tva: '0',
			bundle_cost_ttc: "Prix d'un lot TTC"
		};
		newItem = {
			item_id: '',
			item_name: '',
			amount_of_bundle: 0,
			amount_per_bundle: 0,
			bundle_cost_ht: 0,
			bundle_cost_ttc: 0,
			bundle_cost_float_ttc: 0.0,
			tva: 0
		};
		updateTotalHTandTTC();
	}
</script>

<!-- Popup recap reappro automatique-->
<div
	id="hs-modal-unknown-items"
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
						Articles inconnus
					</h2>
				</div>

				<div class="mt-5">
					<!-- List of unknown items -->
					{#if unknownItems.length > 0}
						<ul class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each unknownItems as item}
								<li class="py-3">
									<div class="flex justify-between items-center">
										<div>
											<p class="text-sm font-medium text-gray-800 dark:text-gray-200">
												Nom : {item.name}
											</p>
											<p class="text-sm text-gray-600 dark:text-gray-400">
												Référence : {item.reference}
											</p>
										</div>
									</div>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="text-sm text-gray-600 dark:text-gray-400">
							Aucun article inconnu trouvé.
						</p>
					{/if}
				</div>

				<div class="mt-5 text-center">
					<button
						type="button"
						class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
						data-hs-overlay="#hs-modal-unknown-items"
					>
						Fermer
					</button>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Popup reappro pdf -->
<div
	id="hs-modal-new-image"
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
						Ajouter une facture PDF
					</h2>
				</div>

				<div class="mt-5">
					<!-- Form -->
					<div class="grid gap-y-4">
						<!-- Form Group -->
						<div>							
							<label for="invoice" class="block text-sm mb-2 dark:text-white">Facture</label>
							<div class="relative">
								<input
									type="file"
									id="invoice"
									name="invoice"
									accept=".pdf"
									bind:this={fileInput}
									on:change={handleFileChange}
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									/>
							</div>

							<button
								type="submit"
								class="mt-4 py-3 px-4 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
								on:click={() => {
									unknownItems = [];
									let items;
									if (newRestock.type == "auchan") {
										items = restockAuchan(pdfText);										
									}
									if (newRestock.type == "auchan_drive"){
										items = restockAuchanDrive(pdfText);
									}									
									if (items != undefined){
										items.forEach(item => {
											itemsApi()
											.getAllItems(undefined, undefined, undefined, undefined, undefined, undefined, {
												withCredentials: true
											})
											.then((res) => {
													let nameItem = undefined;
													let idItem="";
													let amountPerBundle = 0;
													let searchItems = res.data.items ?? [];
													for (let i=0;i<searchItems.length;i++){
														let searchItem = searchItems[i];														
														if (searchItem.ref_bundle == item.reference.toString()){
															nameItem = searchItem.name;
															idItem = searchItem.id;
															if (searchItem.amount_per_bundle != undefined){
																amountPerBundle = searchItem.amount_per_bundle;
															}
															break;
														}
													}
													if (nameItem != undefined){
														newItem.item_id = idItem;
														newItem.item_name = nameItem;
														newItem.amount_of_bundle = item.quantity;
														newItem.amount_per_bundle = amountPerBundle;
														newItem.bundle_cost_ht = Math.round(item.unitPriceHT*100);
														newItem.tva = item.tvaRate*100;
														newItem.bundle_cost_float_ttc = newItem.bundle_cost_ht * (1 + newItem.tva / 10000); 
														newItem.bundle_cost_ttc = Math.round(newItem.bundle_cost_ht * (1 + newItem.tva / 10000));
														add_item_to_restock();		
													}
													else{
														let name = item.name;
														let reference = item.reference;
														unknownItems = [...unknownItems, {name, reference}];
														
													}
												});
									
										});
									}
								}}
								
								data-hs-overlay="#hs-modal-unknown-items">Réappro</button
							>
						</div>
					</div>
					<!-- End Form -->
				</div>
			</div>
		</div>
	</div>
</div>


<div class="max-w-[95%] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto">
	<div class="py-3 px-2 w-1.0 flex m-auto">
		<select
			class="rounded-lg border-transparent appearance-none border border-gray-300 w-96 py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
			placeholder="Type"
			bind:value={newRestock.type}
		>
			{#each restockTypeIterator as [val, name]}
				<option value="{val}">{name}</option>
			{/each}
		</select>
		<div>
			<p class="dark:text-white text-2xl ml-5">
				Total HT : {formatPrice(newRestock.total_cost_ht)}
			</p>
			<p class="dark:text-white text-2xl ml-5">
				Total TTC : {formatPrice(newRestock.total_cost_ttc)}
			</p>
		</div>
		{#if newRestock.type=='auchan' || newRestock.type == 'auchan_drive'}
		<button
			data-hs-overlay="#hs-modal-new-image"
			id="reapproPdf" class="bg-orange-600 hover:bg-orange-700 text-white font-bold py-2 px-4 rounded ml-4">PDF</button>
		{/if}
	</div>
	<div class="flex flex-col">
		<table class="mb-10 min-w-full divide-y divide-gray-200 dark:divide-gray-700 bg-blue-950">
			<thead class="bg-gray-50 dark:bg-blue-600">
				<tr>
					<th scope="col" class="px-12 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nom
						</span>
					</th>
					<th scope="col" class="px-3 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nombre de lots
						</span>
					</th>
					<th scope="col" class="px-3 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nbr produits par lots
						</span>
					</th>
					<th scope="col" class="px-6 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot HT
						</span>
					</th>
					<th scope="col" class="px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							TVA
						</span>
					</th>
					<th scope="col" class="px-6 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot TTC
						</span>
					</th>
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix total HT
						</span>
					</th>
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix total TTC
						</span>
					</th>
					<th scope="col" class="bg-blue-800 px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Ajouter / Supprimer
						</span>
					</th>
				</tr>
			</thead>
			<tr>
				<td class="relative px-12 py-3">
					<div class="flex flex-col">
						<input
							type="text"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent placeholder:text-slate-500"
							placeholder={displayedValues.name}
							on:input={(e) => {
								// @ts-ignore
								searchName = e.target.value.toLowerCase();
								reloadItems();
							}}
							on:change={(e) => {
								// @ts-ignore
								e.target.value = '';
							}}
							on:focusout={() => {
								setTimeout(() => {
									searchName = '';
								}, 200);
							}}
							on:focusin={(e) => {
								// @ts-ignore
								searchName = e.target.value.toLowerCase();
							}}
						/>
					</div>
					<div class="absolute rounded-b-lg bg-slate-100 -translate-y-2 flex flex-col">
						{#if searchName.length > 0}
							{#each items as item}
								<button
									class="p-2"
									on:click={() => {
										displayedValues.name = item.name;
										displayedValues.item_price = formatPrice(item.prices.coutant);
										displayedValues.item_price_ht = formatPrice(
											item.prices.coutant / (1 + (item.last_tva ?? 0) / 10000)
										);
										displayedValues.item_price_calc = item.prices.coutant;
										displayedValues.tva = (item.last_tva ?? 0).toString();
										newItem.tva = item.last_tva ?? 0;
										newItem.item_id = item.id;
										newItem.item_name = item.name;
										searchName = '';
									}}
								>
									{item.name}
								</button>
							{/each}
						{/if}
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nombre de lots"
							min="1"
							max="1000"
							bind:value={newItem.amount_of_bundle}
							on:change={() => {
								if (newItem.amount_of_bundle < 0) {
									newItem.amount_of_bundle = 0;
								}
								updatePrices();
							}}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nombre de produits par lots"
							min="1"
							max="1000"
							bind:value={newItem.amount_per_bundle}
							on:change={() => {
								if (newItem.amount_per_bundle < 0) {
									newItem.amount_per_bundle = 0;
								}
								updatePrices();
							}}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<input
							min="0"
							max="100000"
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder={displayedValues.bundle_cost_ht}
							on:change={(e) => {
								// @ts-ignore
								if (e.target?.value < 0) {
									// @ts-ignore
									e.target.value = 0;
								}
								// @ts-ignore
								newItem.bundle_cost_ht = parsePrice(e.target?.value);
								let r = formatPrice(newItem.bundle_cost_ht);
								displayedValues.bundle_cost_ht = r;
								// @ts-ignore
								e.target.value = r;
								newItem.bundle_cost_ttc = Math.round(
									newItem.bundle_cost_ht * (1 + newItem.tva / 10000)
								);
								newItem.bundle_cost_float_ttc = newItem.bundle_cost_ht * (1 + newItem.tva / 10000);
								displayedValues.bundle_cost_ttc = formatPrice(
									Number((newItem.bundle_cost_ht * (1 + newItem.tva / 10000)).toFixed(0))
								);
							}}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<select
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							id="tva"
							on:change={(e) => {
								// @ts-ignore
								newItem.tva = parseInt(e.target?.value);
								updatePrices();
							}}
							bind:value={displayedValues.tva}
						>
							<option value="0">0%</option>
							<option value="550">5.5%</option>
							<option value="1000">10%</option>
							<option value="2000">20%</option>
						</select>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<input
							min="0"
							max="100000"
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder={displayedValues.bundle_cost_ttc}
							on:change={(e) => {
								newItem.bundle_cost_ht = Math.round(
									// @ts-ignore
									parsePrice(e.target?.value) / (1 + (newItem.tva ?? 0) / 10000)
								);
								// @ts-ignore
								newItem.bundle_cost_ttc = parsePrice(e.target?.value);
								newItem.bundle_cost_float_ttc = 0.0;
								let r = formatPrice(newItem.bundle_cost_ttc);
								displayedValues.bundle_cost_ttc = r;
								displayedValues.bundle_cost_ht = formatPrice(newItem.bundle_cost_ht);
								// @ts-ignore
								e.target.value = r;
							}}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							disabled
							placeholder={displayedValues.item_price_ht}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							disabled
							placeholder={displayedValues.item_price}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<button
							class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
							on:click={add_item_to_restock}
						>
							Ajouter
						</button>
					</div>
				</td></tr
			>
			{#each newRestock.items as item, i}
				<tr>
					<td class="px-12 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.item_name}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<input
								type="number"
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
								min="1"
								max="1000"
								bind:value={item.amount_of_bundle}
							/>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<input
								type="number"
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
								min="1"
								max="1000"
								bind:value={item.amount_per_bundle}
							/>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<input
								type="string"
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
								min="1"
								max="1000"
								value={(item.bundle_cost_ht / 100).toString()}
								on:input={(e) => {
									// @ts-ignore
									item.bundle_cost_ht = parsePrice(e.target?.value);
									item.bundle_cost_ttc = calculateTtc(item);
									updateTotalHTandTTC();
								}}
							/>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<select
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
								id="tva"
								value={item.tva.toString()}
								on:change={(e) => {
									// @ts-ignore
									item.tva = parseInt(e.target?.value);
									item.bundle_cost_ttc = calculateTtc(item);
									updateTotalHTandTTC();
								}}
							>
								<option value="0">0%</option>
								<option value="550">5.5%</option>
								<option value="1000">10%</option>
								<option value="2000">20%</option>
							</select>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<input
								type="string"
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
								min="1"
								max="1000"
								value={(item.bundle_cost_ttc / 100).toString()}
								on:input={(e) => {
									// @ts-ignore
									item.bundle_cost_ttc = parsePrice(e.target?.value);
									item.bundle_cost_ht = calculateHt(item);
									updateTotalHTandTTC();
								}}
							/>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht * item.amount_of_bundle)}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>
									{formatPrice(item.bundle_cost_ttc * item.amount_of_bundle)}
								</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<button
								class="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
								on:click={() => {
									newRestock.items = newRestock.items.filter((_, index) => index !== i);
									updateTotalHTandTTC();
								}}
							>
								Supprimer
							</button>
						</div>
					</td>
				</tr>
			{/each}
		</table>

		<div class="flex p-2 m-8 bg-slate-600 items-center">
			<p class="font-bold text-white text-2xl">
				Ma réappro est irréprochable, et j'en suis responsable :
			</p>
			<input
				id="CHECKBOX"
				class="m-2 mr-auto max-w-lg w-6 h-6"
				type="checkbox"
				bind:checked={sure}
			/>

			{#if sure}
				{#if !selectedEditRestock}
					<button
						on:click={() => applyRestock(RestockState.RestockFinished)}
						class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
					>
						<p class="font-bold text-white text-2xl">Terminer la réappro</p>
					</button>
				{:else}
					<button
						on:click={() => applyEditRestock()}
						class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
					>
						<p class="font-bold text-white text-2xl">Terminer la réappro en attente</p>
					</button>
				{/if}
			{:else if !selectedEditRestock}
				<button
					on:click={() => applyRestock(RestockState.RestockPending)}
					class="bg-gray-800 hover:bg-gray-900 text-white font-bold py-2 px-4 rounded"
				>
					<p class="font-bold text-white text-2xl">Mettre la réappro en attente</p>
				</button>
			{:else}
				<button
					on:click={() => saveEditRestock()}
					class="bg-gray-800 hover:bg-gray-900 text-white font-bold py-2 px-4 rounded"
				>
					<p class="font-bold text-white text-2xl">Sauvegarder la réappro en attente</p>
				</button>
			{/if}
		</div>
	</div>

	<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
		<thead class="bg-gray-50 dark:bg-slate-800">
			<tr>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Date
					</span>
				</th>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Fournisseur
					</span>
				</th>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Conducteur
					</span>
				</th>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						État
					</span>
				</th>
				<th scope="col" class="px-2 py-3">
					<p
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Prix total TTC
					</p>
				</th>
				<th scope="col" class="px-2 py-3">
					<p
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Actions
					</p>
				</th>
			</tr>

			{#if deletingRestock}
				<ConfirmationPopup
					message={confirmationMessage}
					confirm_text="Supprimer"
					cancel_callback={() => {
						deletingRestock = false;
					}}
					confirm_callback={deleteRestockCallback}
				/>
			{/if}
			{#each restocks as restock}
				<tr class={restock.state == RestockState.RestockPending ? 'bg-orange-700' : 'bg-green-700'}>
					<td class="px-6 py-3">
						<div class="text-center flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p><Time timestamp={restock.created_at * 1000} format="YYYY/MM/DD HH:mm:ss" /></p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="text-center flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{restock.type}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="text-center flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{restock.created_by_name}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="text-center flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{restock.state}</p>
							</div>
						</div>
					</td>
					<td class="px-2 py-3">
						<p
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							{formatPrice(restock.total_cost_ttc)}
						</p>
					</td>
					<td class="px-2 py-3">
						<div
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							<button
								class="px-4 py-2 bg-red-600 text-white font-semibold rounded-lg shadow-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-green-400 focus:ring-offset-2 transition ease-in-out duration-200"
								on:click={() => {
									deleteRestockCallback = () => {
										deletingRestock = false;
										deleteRestock(restock.id);
									};
									confirmationMessage =
										'Supprimer la réappro de ' +
										restock.created_by_name +
										' à ' +
										restock.type +
										' ?';
									deletingRestock = true;
								}}
							>
								Supprimer
							</button>
							<button
								class="px-4 py-2 bg-blue-500 text-white font-semibold rounded-lg shadow-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-green-400 focus:ring-offset-2 transition ease-in-out duration-200"
								on:click={() => {
									if (selectedViewRestock == restock) {
										selectedViewRestock = undefined;
									} else {
										selectedViewRestock = restock;
									}
								}}
							>
								{#if selectedViewRestock == restock}
								 	Fermer
								{:else}
									Voir
								{/if}
							</button>
							{#if restock.state == RestockState.RestockPending}
								<button
									class="px-4 py-2 bg-orange-500 text-white font-semibold rounded-lg shadow-md hover:bg-orange-800 focus:outline-none focus:ring-2 focus:ring-green-400 focus:ring-offset-2 transition ease-in-out duration-200"
									on:click={() => {
										if (selectedEditRestock == restock) {
											selectedEditRestock = undefined;
											newRestock.items = [];
											updateTotalHTandTTC();
										} else {
											selectedEditRestock = restock;
											newRestock = structuredClone(selectedEditRestock);
										}
									}}
								>
									{#if selectedEditRestock == restock}
										Fermer sans sauvegarder
									{:else}
										Modifier
									{/if}
								</button>
							{/if}
						</div>
					</td>
				</tr>
			{/each}
		</thead>
	</table>
	<div
		class="px-2 py-4 grid gap-3 md:flex md:justify-between md:items-center border-t border-gray-200 dark:border-gray-700"
	>
		<div>
			<p class="text-sm text-gray-600 dark:text-gray-400">
				<span class="font-semibold text-gray-800 dark:text-gray-200">{restocks.length}</span>
				résultats
			</p>
		</div>

		<div>
			<div class="inline-flex gap-x-2">
				<button
					type="button"
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
					on:click={prevPage}
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
					Page {pageRestock} / {maxPageRestock}
				</p>

				<button
					type="button"
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
					on:click={nextPage}
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
	{#if selectedViewRestock != undefined}
		<table class="mb-10 min-w-full divide-y divide-gray-200 dark:divide-gray-700 bg-blue-950">
			<thead class="bg-gray-50 dark:bg-blue-600">
				<tr>
					<th scope="col" class="px-12 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nom
						</span>
					</th>
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix coûtant HT
						</span>
					</th>
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix coûtant TTC
						</span>
					</th>
					<th scope="col" class="px-3 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nombre de lots
						</span>
					</th>
					<th scope="col" class="px-3 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nbr produits par lots
						</span>
					</th>
					<th scope="col" class="px-6 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot HT
						</span>
					</th>
					<th scope="col" class="px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							TVA
						</span>
					</th>
					<th scope="col" class="px-6 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot TTC
						</span>
					</th>
				</tr>
			</thead>
			{#each selectedViewRestock.items as item}
				<tr>
					<td class="px-12 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.item_name}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht * item.amount_of_bundle)}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>
									{formatPrice(item.bundle_cost_ttc * item.amount_of_bundle)}
								</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.amount_of_bundle}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.amount_per_bundle}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht)}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.tva / 100}%</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ttc)}</p>
							</div>
						</div>
					</td>
				</tr>
			{/each}
		</table>
	{/if}
</div>
