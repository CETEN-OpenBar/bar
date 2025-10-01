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
		link: '/panel/products/restocks'
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

<style>
	.modules-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
		gap: 1rem;
		justify-content: center;
		justify-items: center;
		padding: 0.5rem;
		width: 100%;
		max-width: 1200px;
		margin: 0 auto;
	}

	.module-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 8rem;
		width: 16rem;
		word-wrap: break-word;
		text-align: center;
		border-radius: 0.5rem;
		box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
		cursor: pointer;
	}

	.module-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: white;
		padding: 1.25rem;
	}

	/* Color classes */
	.bg-green-600 {
		background-color: rgb(22, 163, 74);
	}

	.bg-red-600 {
		background-color: rgb(220, 38, 38);
	}

	.bg-blue-600 {
		background-color: rgb(37, 99, 235);
	}

	.bg-yellow-600 {
		background-color: rgb(202, 138, 4);
	}
</style>

<div class="modules-grid">
	{#each modules as module}
			<a href={module.link}>
				<div
					class="module-card {module.color}"
				>
					<p class="module-title">{module.name}</p>
				</div>
			</a>
		{/each}
</div>
