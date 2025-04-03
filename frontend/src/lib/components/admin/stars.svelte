<script lang="ts">
import type { Account } from '$lib/api';
import { starsApi } from '$lib/requests/requests';

export let account: Account;
export let close: () => void;

let amount: number = 0;

function sendStars() {
  if (amount >= 0) {
    starsApi()
      .postStarring(account.id, amount, "staff", { withCredentials: true })
      .then(() => {
        close();
      });
  }
}
</script>

<!-- Popup overlay -->
<button
  id="overlay"
  class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-10 hover:cursor-default"
  on:click={close}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
  <div
    class="relative text-black flex flex-col justify-center items-center gap-4 p-10 bg-white rounded-xl shadow-xl z-20"
  >
    <button
      class="absolute top-0 right-0 p-2 text-xl font-bold m-2 rounded-full transition-all text-black"
      on:click={close}
    >
      <iconify-icon icon="mdi:close" />
    </button>

    <h1 class="text-3xl mb-4">Ajouter des étoiles</h1>

    <div class="flex flex-col gap-4 items-center">
      <div class="w-full relative">
        <input
          type="number"
          required
          class="w-full p-2 border-2 rounded-md"
          bind:value={amount}
          min="0"
          placeholder="Nombre d'étoiles"
        />
        <iconify-icon icon="mdi:star" class="absolute right-9 top-1/2 -translate-y-1/2 text-yellow-500" />
      </div>

      <button
        class="px-4 py-2 bg-blue-500 text-white rounded-md disabled:bg-gray-400"
        disabled={amount < 0}
        on:click={sendStars}
      >
        Valider
      </button>
    </div>
  </div>
</div>
