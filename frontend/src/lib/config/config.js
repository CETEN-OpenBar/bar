// Reads /config.json and store it locally in a variable to redistribute it to the app, we're using VueJS
import axios from "axios";

let c = null;

export const api = () => {
    return c.api;
};

export const loadConfig = async () => {
    if (c != null) {
        return c;
    }
    
    // Use axios get and wait for it to finish
    await axios.get("/config.json").then((response) => {
        c = response.data;
    });

    return c;
};
