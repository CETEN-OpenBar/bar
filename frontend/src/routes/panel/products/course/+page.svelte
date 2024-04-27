<script lang="ts">
	import { api } from '$lib/config/config';
	import { CourseApi } from '$lib/requests/requests';
	import type { CourseItem, Item, RestockType } from '$lib/api';

	let items: CourseItem[] = [];
    let fournisseur = "promocash";
    function reloadCourse() {
        CourseApi()
            .getCourse(fournisseur, {
                withCredentials: true
            })
            .then((res) => {
                if (res.data.items != null) {
                    items = res.data.items;
                } else {
                    items = [];
                }
            });
    }
    reloadCourse()
</script>

<div class="relative mt-4 w-96 md:mt-0">
    <!-- filter by state -->
    <select
        id="fournisseur"
        name="fournisseur"
        class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
        required
        aria-describedby="text-error"
        on:change={(e) => {
            // @ts-ignore
            let val = e.target?.value;
            fournisseur = val;
            reloadCourse();
        }}
    >
        <option value="promocash">Promocash</option>
        <option value="auchan_drive">Auchan drive</option>
        <option value="auchan">Auchan</option>
        <option value="viennoiserie">Boulangerie Benoist</option>
    </select>
</div>

<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
	<thead class="bg-gray-50 dark:bg-slate-800">
        <tr>
            <th scope="col" class="px-6 py-3">
                <span
                    class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
                >
                    Nom
                </span>
            </th>
            <th scope="col" class="px-6 py-3">
                <span
                    class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
                >
                    Nombre à acheter
                </span>
            </th>
        </tr>
    </thead>
    <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
        <td>
        {#each items as item}
            <p class="py-3 px-2 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400">
            {item.item.name}
            </p>
        {/each}
        </td>
        <td>
            {#each items as item}
                <p class="py-3 px-2 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400">
                {item.amountToBuy}
                </p>
            {/each}
        </td>
    </tbody>
</table>

