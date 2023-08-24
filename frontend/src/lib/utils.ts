export const parsePrice = (price: string) => {
    let splt: Array<string> = [price];
    if (price.includes(',')) {
        splt = price.split(',');
    } else if (price.includes('.')) {
        splt = price.split('.');
    }

    const cents = parseInt(splt.length>1 ? splt[1] : '00');
    const euros = parseInt(splt[0])

    if (isNaN(cents) || isNaN(euros)) throw new Error('Invalid price');
    if (cents >= 100 || cents < 0 || euros < 0) throw new Error('Invalid price');
    
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