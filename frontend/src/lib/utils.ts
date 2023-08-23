export const parsePrice = (price: string) => {
    let splt: Array<string> = [price];
    if (price.includes(',')) {
        splt = price.split(',');
    } else if (price.includes('.')) {
        splt = price.split('.');
    }
    
    let res = 0;
    if (price.length > 1) res = parseInt(price[0]) * 100 + parseInt(price[1]);
    else res = parseInt(price[0]) * 100;
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