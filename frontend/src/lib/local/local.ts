import { local_api } from '$lib/config/config';


export const open_door = (card_id: string, password: string) => {
    const data = {
        card_id: card_id,
        card_pin: password
    };
    fetch(local_api() + '/porte', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
}

export const open_caisse = (card_id: string, password: string) => {
    const data = {
        card_id: card_id,
        card_pin: password
    };
    fetch(local_api() + '/caisse', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
}

export const open_ventilo = (card_id: string, password: string) => {
    const data = {
        card_id: card_id,
        card_pin: password
    };
    fetch(local_api() + '/ventilo', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
}