import { Fournisseur, RestockType } from "./api";

export const parsePrice = (price: string) => {
    let splt: Array<string> = [price];
    if (price.includes(',')) {
        splt = price.split(',');
    } else if (price.includes('.')) {
        splt = price.split('.');
    } 

    if (splt.length>1 && splt[1].length === 1) splt[1] += '0';
    const cents = parseInt(splt.length>1 ? splt[1] : '00');
    const euros = parseInt(splt[0])

    if (isNaN(cents) || isNaN(euros)) throw new Error('Invalid price - 1');
    if (cents >= 100 || cents < 0 || euros < 0) throw new Error('Invalid price - 2');
    
    let res = 0;
    if (splt.length > 1) res = euros * 100 + cents;
    else res = euros * 100;
    return res;
}

export const formatPrice = (price: number) => {
    // Price is in cents, so we divide by 100 to get dollars
    return `${(price / 100).toFixed(2)} €`;
}

export const formatDate = (date: number): string => {
    // date is in unix seconds
    const d = new Date(date * 1000);
    return d.toLocaleDateString('fr-FR');
}

export const formatDateTime = (date: number): string => {
    // date is in unix seconds
    const d = new Date(date * 1000);
    // We need day and hour
    return `${d.toLocaleDateString('fr-FR')} ${d.toLocaleTimeString('fr-FR')}`;
}

export const file2Base64 = (file: File): Promise<string> => {
    return new Promise<string>((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result?.toString() || '');
        reader.onerror = (error) => reject(error);
    });
};

export const time2Utc = (time: number): number => {
    // Time is in second, we need to get the delta between local time and UTC time
    const d = new Date();
    const delta = d.getTimezoneOffset() * 60;
    return time + delta;
}

export const enumIterator = (e: any): Iterable<[value: string, name: string]> => {
    return Object.entries(Fournisseur).map(([key, value]) => {
        console.log("Key : " + key + " val : " + value)
        return [value, formatEnumName(key)]
    });
}

export const formatEnumName = (name: string): string => {
    // Replace '_' with spaces and capitalize every word
    return name.replaceAll("_", " ")
        .toLowerCase()
        .split(' ')
        .map(word => word.charAt(0).toUpperCase() + word.substring(1))
        .join(' ');
}

export const fournisseurIterator = enumIterator(Fournisseur);
export const restockTypeIterator = enumIterator(RestockType);
