<script lang="ts">
	import type {TransactionItem} from '$lib/api';
	import { api } from '$lib/config/config';
	import { itemsApi, transactionsApi } from '$lib/requests/requests';
	import { onDestroy, onMount } from 'svelte';
	import { dragscroll } from '@svelte-put/dragscroll';
	import Autodisconnect from '../random/autodisconnect.svelte';
	import Item from '../borne/items.svelte';

	let transaction_items: Array<TransactionItem> = [];
	let amount_left: Array<number> = [];
    let name: string = "";
        
    $: filtered_items = name !== "" ?
        transaction_items.filter((item) => {
            return item.item_name.toLocaleLowerCase().includes(name.toLocaleLowerCase())
        }) : transaction_items
    
	onMount(() => {
		reloadTransactionsItems();
	});

	onDestroy(() => {
	});

	function reloadTransactionsItems() {
		transactionsApi()
			.getTransactionsItems(undefined, { withCredentials: true })
			.then((res) => {
				if (!(res.data instanceof Array)) return;
				transaction_items = res.data;
                 //stock amount left of each item
                transaction_items.forEach((transaction_item, index) => {
                    itemsApi()
                    .getItem(transaction_item.item_id, { withCredentials: true })
                    .then((res) => {
                        amount_left[index] = res.data.amount_left; //stock amount left of each item
                    });
            });
			});
        
	}

    
</script>

<!-- Good looking dropdown for transaction -->
<div class="w-full flex justify-center">
    <!-- show clearly if there's a scrollbar -->
    <div use:dragscroll class="overflow-auto max-h-[80vh] mt-5">
        <table class="table-fixed min-w-[60vw] text-center text-xl border-separate border-spacing-y-1 max-h-fit">
            <thead class="bg-slate-500 sticky top-0">
                <th class="w-[48%] p-2 flex flex-row justify-items items-center">
                    <div>Nom</div>
                    <input class="mx-5 rounded-md p-2 text-black" placeholder="rechercher" bind:value={name}/>
                </th>
                <th class="w-[13%]">Commandés</th>
                <th class="w-[13%]">Terminés</th>
                <th class="w-[13%]">Restants</th>
                <th class="w-[13%]">En stock</th>

            </thead>
    
            <tbody class="overflow-scroll max-h-fit">
                {#each filtered_items as item, index}
                    <tr class="bg-slate-800">
                        <td class="p-4">
                            <div class="flex flex-row items-center space-x-5">
                                <img
                                    src={api() + item.picture_uri}
                                    alt="item"
                                    class="w-10 h-10 rounded-2xl self-center"
                                />
                                <div class="">{item.item_name}</div>
                            </div>
                        </td>
                        <td>{item.item_amount}</td>
                        <td>{item.item_already_done}</td>
                        <td>{item.item_amount - item.item_already_done}</td>
                        <td>{amount_left[index]}</td>




                    </tr>
                {/each}
    
            </tbody>
        </table>
    </div>
    
</div>
