<script lang="ts">
	import { api } from '$lib/config/config';
	import { CourseApi } from '$lib/requests/requests';
	import type { CourseItem, Item } from '$lib/api';

	let items: CourseItem[] = [];
	CourseApi()
		.getCourse('promocash', {
			withCredentials: true
		})
		.then((res) => {
			items = res.data.items;
		});
</script>

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
        {#each items as item}
        <td class="h-px">
            <p class="py-3 px-2 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400">
            {item.item.name}
            </p>
        </td>
        <td class="h-px">
            <p class="py-3 px-2 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400">
            {item.amountToBuy}
            </p>
        </td>
        {/each}
    </tbody>
</table>

