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
    <div class="flex flex-col items-center justify-center">
        <a href={api() + "/auth/google?r=admin"} class="connect-button mt-4">Connection avec Google</a>
    </div>
</div>
