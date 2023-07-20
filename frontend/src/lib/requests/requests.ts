import { Configuration } from '$lib/api';
import { AuthApiFactory, AccountsApiFactory } from '$lib/api';
import { api } from '$lib/config/config';

export const authApi = AuthApiFactory(
	new Configuration({
		basePath: api()
	})
);

export const accountsApi = AccountsApiFactory(
    new Configuration({
        basePath: api()
    })
);