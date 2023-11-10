export const parsePrice = (price: string) => {
    let splt: Array<string> = [price];
    if (price.includes(',')) {
        splt = price.split(',');
    } else if (price.includes('.')) {
        splt = price.split('.');
    } 

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