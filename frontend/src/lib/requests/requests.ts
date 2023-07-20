import { Configuration } from '$lib/api';
import { AuthApiFactory, AccountsApiFactory, ItemsApiFactory, DeletedApiFactory, RefillsApiFactory, CarouselApiFactory, CategoriesApiFactory, TransactionsApiFactory } from '$lib/api';
import { api, local_token } from '$lib/config/config';

export const authApi = AuthApiFactory(
	new Configuration({
		basePath: api(),
        apiKey: ((name: string) => { console.log(name); if (name == 'X-Local-Token') { return local_token() } else { return '' } })
	})
);

export const accountsApi = AccountsApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);

export const itemsApi = ItemsApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);

export const deletedApi = DeletedApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);

export const refillsApi = RefillsApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);

export const carouselApi = CarouselApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);

export const categoriesApi = CategoriesApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);

export const transactionsApi = TransactionsApiFactory(
    new Configuration({
        basePath: api(),
        apiKey: ((name: string) => { if (name == 'X-Local-Token') { return local_token() } else { return '' } })
    })
);
