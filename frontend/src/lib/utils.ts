export const formatPrice = (price: number) => {
    // Price is in cents, so we divide by 100 to get dollars
    return `${(price / 100).toFixed(2)} €`;
}

export const formatDate = (date: number): string => {
    // date is in unix seconds
    const d = new Date(date * 1000);
    return d.toLocaleDateString('fr-FR');
}