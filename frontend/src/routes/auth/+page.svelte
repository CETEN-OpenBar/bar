<script lang="ts">
    import { api } from "$lib/config/config";
    import { page } from '$app/stores';
	import Error from "$lib/components/error.svelte";

    const noAccountError = $page.url.searchParams.has("noaccount")
</script>

<style>
    .connect-button {
        background-color: #4285F4;
        color: white;
        font-size: 1.25rem;
        padding: 12px 24px;
        border-radius: 8px;
        text-decoration: none;
        transition: background-color 0.2s ease;
    }

    .connect-button:hover {
        background-color: #3367D6;
    }
</style>

<svelte:head>
    <title>OpenBar - Connexion</title>
</svelte:head>

<!-- Connect with Google at the middle of the screen -->
<div class="flex flex-col items-center justify-center h-screen">
    <h1 class="font-semibold text-3xl m-3">OpenBar</h1>
    <img 
        src="/logo.png"
        alt="Logo du Bar"
        class="h-48"
    />
    <div class="flex flex-col items-center justify-center">
        <a href={api() + "/auth/google?r=client/commande"} class="connect-button mt-4">Connexion avec Google</a>
    </div>

    {#if noAccountError}
    <div class="mt-3 flex flex-col items-center bg-red-500 rounded-lg shadow-lg p-6 text-white max-w-[95vw] md:max-w-lg space-y-2 text-center">
        <div class="flex flex-row items-center">
            <!-- error icon -->
            <svg class="w-10 h-10 text-white mr-3" viewBox="0 0 24 24">
                <path
                    fill="currentColor"
                    d="M12 2C6.486 2 2 6.486 2 12s4.486 10 10 10 10-4.486 10-10S17.514 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"
                />
            </svg>
            <h1 class="text-2xl font-bold">Compte introuvable</h1>
        </div>
        <p>
            Ce compte TN.net n'est lié à aucun compte OpenBar.
        </p>
        <p>
            Vous pouvez lier votre compte Google à votre compte OpenBar depuis les bornes.
        </p>
        <p>
            Si le problème persiste ou que vous avez des difficultés pour lier vos comptes, n'hésitez pas à contacter un respo info du bar.
        </p>
    </div>
    {/if}

</div>
