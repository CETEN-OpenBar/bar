// Reads /config.json and store it locally in a variable to redistribute it to the app, we're using VueJS
import axios from "axios";

// Create type for config
export type Config = {
    api: string;
};

// Create variable to store config
let c: Config | null = null;

export const api = () => {
    if (c == null) {
        throw new Error("Config not loaded");
    }
    return c.api;
};

export const loadConfig = async () => {
    if (c != null) {
        return c;
    }

    // Do not render on server
    if (typeof window === "undefined") {
        return null;
    }
    
    // Use axios get and wait for it to finish
    await axios.get("/config.json").then((response) => {
        c = response.data;
    });

    return c;
};
