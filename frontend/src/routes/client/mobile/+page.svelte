<script lang="ts">
	// Check the query params
	const urlParams = new URLSearchParams(window.location.search);
	const returnType = urlParams.get('rt');
	const returnMessage = urlParams.get('rm');
</script>

<div class="container">
	{#if returnType == 'authSuccess'}
		<span class="check-icon">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				version="1.1"
				width="270"
				height="270"
			>
				<path d="M 30,180 90,240 240,30" style="stroke:#050; stroke-width:30; fill:none" />
			</svg></span
		>
		<p class="message">Connexion réalisée!</p>
	{/if}

	{#if returnType == 'authError'}
        <span class="check-icon">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                version="1.1"
                width="270"
                height="270"
                class="text-red-300"
            >
                <path d="M 30,30 240,240 M 30,240 240,30" style="stroke:#500; stroke-width:30; fill:none" />
            </svg>
        </span>
        <p class="message">Il y a eu un problème : {returnMessage}</p>
    
    {/if}
</div>

<style>
	/* Center the content vertically and horizontally */
	.container {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		height: 100vh;
		padding: 1rem;
	}

	.check-icon {
		font-size: 6rem;
		color: #10b981;
	}

	.message {
		font-size: 1.5rem;
		text-align: center;
		margin-top: 1rem;
	}
</style>
