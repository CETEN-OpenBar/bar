<script lang="ts">
	import type { Transaction } from "$lib/api";
    import { transactionsApi } from "$lib/requests/requests";
    import { formatPrice } from "$lib/utils";
    import { onMount } from "svelte";

    export let amount: number = 3;

    let transactions: Array<Transaction> = [];

    onMount(() => {
        transactionsApi()
            .getTransactions(0, amount, undefined, { withCredentials: true })
            .then((res) => {
                if (res.data.transactions instanceof Array) transactions = res.data.transactions;
            });
    });

</script>

{JSON.stringify(transactions)}