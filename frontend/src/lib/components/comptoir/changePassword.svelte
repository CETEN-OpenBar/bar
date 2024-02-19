<script lang="ts">
	import { accountsApi } from '$lib/requests/requests';
	import Error from '../error.svelte';
	import Password from '../password.svelte';
	import Success from '../success.svelte';

	export let onEnd = () => {};

	let state = '';
	let timeout: number;

	let oldPassword = '';
	let newPassword = '';

	function setNewPassword(old: string, new_: string) {
		if (new_ == '') {
			state = '';
			onEnd();
			return;
		}

		accountsApi()
			.patchAccountPassword(
				{
					old_password: old,
					new_password: new_
				},
				{
					withCredentials: true
				}
			)
			.then(() => {
				state = 'success';
			})
			.catch(() => {
				state = 'error';
			})
			.finally(() => {
				timeout = setTimeout(() => {
					state = '';
					onEnd();
				}, 3000);
			});
	}

	function nextStep(password: string) {
		oldPassword = password;
		state = 'new';
	}

	function nextStep2(password: string) {
		newPassword = password;
		setNewPassword(oldPassword, newPassword);
	}
</script>

<Password custom_text="Entrez votre ancien mot de passe" callback={nextStep} />
{#if state == 'new'}
	<Password custom_text="Entrez votre nouveau mot de passe" callback={nextStep2} />
{/if}
{#if state == 'success'}
	<Success message="Votre mot de passe a bien été changé" />
{/if}
{#if state == 'error'}
	<Error error="Une erreur est survenue" />
{/if}
