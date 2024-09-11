<script lang="ts">
    import { api } from "$lib/config/config";
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
</div>
