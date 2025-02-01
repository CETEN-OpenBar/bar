<script lang="ts">
	import { page } from '$app/stores';

	// New type for modules with Name, Color and Link
	type Module = {
		name: string;
		color: string;
		link: string;
	};

	// Array of modules
	let modules = [] as Module[];
	let currentPage = 0;

	// Add modules to array (append more when editing modules)
	// TODO: make sure all modules are here

	// Color code :
	// Blue => panel related
	// Green => product related
	// Red => account related

	modules.push({
		name: 'Création de menu',
		color: 'bg-green-600',
		link: '/panel/products/create-menu'
	});

	modules.push({
		name: 'Dé/bloquer un compte',
		color: 'bg-red-600',
		link: '/panel/accounts/restrict'
	});

	modules.push({
		name: 'Changer un surnom',
		color: 'bg-red-600',
		link: '/panel/accounts/nickname'
	});

	modules.push({
		name: 'Réappro',
		color: 'bg-green-600',
		link: '/panel/products/reappro'
	});

	modules.push({
		name: 'Réinitialiser un code pin',
		color: 'bg-red-600',
		link: '/panel/accounts/reset_pin'
	});

	modules.push({
		name: 'Repositioner les catégories',
		color: 'bg-blue-600',
		link: '/panel/categories/reposition'
	});

	modules.push({
		name: 'Trésorerie',
		color: 'bg-blue-600',
		link: '/panel/treso'
	});
	modules.push({
		name: 'Incohérences',
		color: 'bg-green-600',
		link: '/panel/products/incoherants'
	});
	modules.push({
		name: 'Fournisseur',
		color: 'bg-green-600',
		link: '/panel/products/fournisseur'
	});
	modules.push({
		name: 'Course',
		color: 'bg-green-600',
		link: '/panel/products/course'
	});
	modules.push({
		name: 'Comptoir',
		color: 'bg-yellow-600',
		link: '/comptoir/c/transactions'
	});

	// Sort the modules by name
	modules.sort((a, b) => a.name.localeCompare(b.name));
</script>

<div class="grid lg:grid-cols-4 md:grid-cols-3 sm:grid-cols-2 grid-cols-1 gap-4">
	{#each modules as module, i}
		<!-- only display modules that are in the current page -->
		{#if i >= currentPage * 12 && i < (currentPage + 1) * 12}
			<a href={module.link}>
				<div
					class="flex flex-col items-center justify-center h-32 w-64 break-words flex-wrap text-center rounded-lg shadow-lg cursor-pointer {module.color}"
				>
					<p class="text-2xl font-bold text-white p-5">{module.name}</p>
				</div>
			</a>
		{/if}
	{/each}

	<!-- if page is not whole, fill it with invisible modules -->
	{#if currentPage == modules.length % 12 && modules.length % 12 != 0}
		{#each Array(12 - (modules.length % 12)) as _}
			<div
				class="flex flex-col items-center justify-center h-32 w-64 break-words flex-wrap text-center rounded-lg shadow-lg"
			>
				<p class="text-2xl font-bold text-white p-5" />
			</div>
		{/each}
	{/if}
</div>

<!-- arrows for navigation -->
<div id="pagination" class="mt-5">
	<button
		class=" text-white font-bold py-2 px-4 rounded-l {currentPage == 0
			? 'bg-gray-400 pointer-events-none'
			: 'bg-blue-600 hover:bg-blue-700'}"
		on:click={() => {
			if (currentPage > 0) currentPage--;
		}}
	>
		Précédent
	</button>
	<button
		class="text-white font-bold py-2 px-4 rounded-r {(currentPage + 1) * 12 < modules.length
			? 'bg-blue-600 hover:bg-blue-700'
			: 'bg-gray-400 pointer-events-none'}"
		on:click={() => {
			if ((currentPage + 1) * 12 < modules.length) currentPage++;
		}}
	>
		Suivant
	</button>
</div>
