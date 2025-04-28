<script lang="ts">
  import { onMount } from 'svelte';
  import type { Starring } from '$lib/api';
  import { starsApi } from '$lib/requests/requests';

  let stars: Starring[] = [];
  let displayed: Starring[] = [];
  let start_date = '';
  let end_date = '';
  let nameFilter = '';

  async function fetchStars() {
    const response = await starsApi().getStarrings(undefined, undefined, start_date, end_date, { withCredentials: true });
    console.log(response.data)
    stars = Array.isArray(response.data.stars) ? response.data.stars : [];
    applyFilter();
  }

  function applyFilter() {
    const term = nameFilter.trim().toLowerCase();
    displayed = term
      ? stars.filter((s) => s.account_name.toLowerCase().includes(term))
      : stars;
  }

  onMount(fetchStars);
</script>

<div class="mb-4 flex space-x-4">
  <label>
    Date début:
    <input
      type="date"
      bind:value={start_date}
      on:change={fetchStars}
      class="border px-2 py-1 rounded"
    />
  </label>
  <label>
    Date fin:
    <input
      type="date"
      bind:value={end_date}
      on:change={fetchStars}
      class="border px-2 py-1 rounded"
    />
  </label>
  <label>
    Compte:
    <input
      type="text"
      placeholder="Nom"
      bind:value={nameFilter}
      on:input={applyFilter}
      class="border px-2 py-1 rounded"
    />
  </label>
</div>

<table class="min-w-full border-collapse">
  <thead class="bg-gray-100">
    <tr>
      <th class="border p-2">Date</th>
      <th class="border p-2">Compte</th>
      <th class="border p-2">Émetteur</th>
      <th class="border p-2">Montant</th>
      <th class="border p-2">Type</th>
      <th class="border p-2">État</th>
    </tr>
  </thead>
  <tbody>
    {#each displayed as star}
      <tr class="hover:bg-gray-50">
        <td class="border p-2">{new Date(star.issued_at).toLocaleString()}</td>
        <td class="border p-2">{star.account_name}</td>
        <td class="border p-2">{star.issued_by_name}</td>
        <td class="border p-2">{star.amount}</td>
        <td class="border p-2">{star.type}</td>
        <td class="border p-2">{star.state}</td>
      </tr>
    {/each}
  </tbody>
</table>
