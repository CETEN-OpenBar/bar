export const formatPrice = (price: number) => {
    // Price is in cents, so we divide by 100 to get dollars
    return `${(price / 100).toFixed(2)} €`;
}