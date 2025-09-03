import { CashMovementsApiFactory, Configuration } from '$lib/api';
import {
	RestocksApiFactory,
	AuthApiFactory,
	AccountsApiFactory,
	ItemsApiFactory,
	DeletedApiFactory,
	RefillsApiFactory,
	CarouselApiFactory,
	CategoriesApiFactory,
	TransactionsApiFactory,
	CourseApiFactory,
	StarsApiFactory
} from '$lib/api';
import { api, local_token } from '$lib/config/config';

export const authApi = () => {
	return AuthApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const accountsApi = () => {
	return AccountsApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const itemsApi = () => {
	return ItemsApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const deletedApi = () => {
	return DeletedApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const refillsApi = () => {
	return RefillsApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const carouselApi = () => {
	return CarouselApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const categoriesApi = () => {
	return CategoriesApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const transactionsApi = () => {
	return TransactionsApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const restocksApi = () => {
	return RestocksApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const cashMovementsApi = () => {
	return CashMovementsApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const CourseApi = () => {
	return CourseApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};

export const starsApi = () => {
	return StarsApiFactory(
		new Configuration({
			basePath: api(),
			apiKey: (name: string) => {
				if (name == 'X-Local-Token') {
					return local_token();
				} else {
					return '';
				}
			}
		})
	);
};
