// Reads /config.json and store it locally in a variable to redistribute it to the app, we're using VueJS
import axios from "axios";

// Create type for config
export type Config = {
    apiws: string;
    api: string;
    local_api: string;
    local_token: string;
};

// Create variable to store config
let c: Config | null = null;

export const api = () => {
    if (c == null) {
        throw new Error("Config not loaded");
    }
    return c.api;
};

export const apiws = () => {
    if (c == null) {
        throw new Error("Config not loaded");
    }
    return c.apiws;
};

export const local_api = () => {
    if (c == null) {
        throw new Error("Config not loaded");
    }
    return c.local_api;
}

export const local_token = () => {
    if (c == null) {
        throw new Error("Config not loaded");
    }
    return c.local_token;
};

export const loadConfig = () => {
    return new Promise((resolve, reject) => {
        if (c != null) {
            resolve(c);
        }
    
        // Do not render on server
        if (typeof window === "undefined") {
            return null;
        }
        
        // Use axios get and wait for it to finish
        axios.get("/config.json").then((response) => {
            c = response.data;
            resolve(c);
        }).catch((error) => {
            reject(error);
        });
    });
};
