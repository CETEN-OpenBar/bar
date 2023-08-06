<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';

	accountsApi()
		.getAccountAdmin({ withCredentials: true })
		.then((res) => {
			if (!res.data.can_restore) {
				goto('/auth');
			}
		})
		.catch(() => {
			goto('/auth');
		});
</script>

<slot/>