// Create a svelte store here to pass context to the app
import type { Account } from "$lib/api";
import { writable } from "svelte/store";

// Create type for store
export type Store = {
    account: Account | undefined;
};

// Create variable to store store
let s: Store = {} as Store;

// Create writable store
export const store = writable(s);

// Variable for search name for /comptoir/c/transaction
export const searchName = writable('');

// Create function to set store
export const setStore = (store: Store) => {
    s = store;
}
