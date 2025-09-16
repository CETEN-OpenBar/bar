<script lang="ts" context="module">
	function isDigit(chr: string): boolean{
		return chr >= '0' && chr <= '9';
	}

	function isPrice(str: string): boolean{
		const len = str.length;
		return len >= 4 && isDigit(str[len-1]) && isDigit(str[len-2]) && str[len-3] == ',' && isDigit(str[len-4]);
	}

    export function restockAuchanDrive(input: string): Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> {
		const splitPage = "Référence   Caractéristiques produit   Prix U. (HT) € Remises U. (HT) €   Qte.   Prix total (HT) € Taux TVA % Cagnotte WAAOH! Prix total (TTC) €  "
		const splitLine = "   ";
		const endPage1 = "Votre commande :";
		const endPage2 = "Taux TVA";	
		const endPage3 = "Votre facture :";
		const pages = input.split(splitPage);
		const itemsList: Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> = [];
		for (let i=1;i<pages.length;i++){
			const pageSplit = pages[i].split(splitLine);
			for (let j = 0; j < pageSplit.length; ) {
				if (pageSplit[j].includes(endPage1) || pageSplit[j].includes(endPage2) || pageSplit[j].includes(endPage3)) {				
					break;
				}
				const reference = parseInt(pageSplit[j].split(" ").slice(-1)[0]);
				j++;
				let name;
				let unitPriceHT;
				let splitName = pageSplit[j].split(" ");
				if (isPrice(splitName[splitName.length-1])){
					name = splitName.slice(0, -1).join(" ");
					unitPriceHT = parseFloat(splitName[splitName.length-1].replace(",","."));
				}
				else{
					name = pageSplit[j];
					j++;
					unitPriceHT = parseFloat(pageSplit[j].replace(",","."));
				}
				j++;
				if (isPrice(pageSplit[j])){ // Check if there is a discount
					unitPriceHT -= parseFloat(pageSplit[j].replace(",","."));
					j++;
				}
				const quantity = parseInt(pageSplit[j]);
				j+=2;
				const tvaRate = parseFloat(pageSplit[j].replace(",","."));
				j++;
				if (isPrice(pageSplit[j])){ // Check if there is a cagnotte WHAOOH					
					j++;
				}				
				itemsList.push({
					reference,
					name,
					quantity,
					unitPriceHT,
					tvaRate
				});

			}
		}
		return itemsList;
	}

	export function restockAuchan(input: string): Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> {
		const splitPage = "Référence   Caractéristiques produit   Prix U. (HT) € Remises U. (HT) €   Qte.   Prix total (HT) € Taux TVA % Prix total (TTC) €  ";
		const splitLine = "   ";
		const splitNameQuantity = "  Dont éco-participation :  ";
		const endPage1 = "Votre commande :";
		const endPage2 = "TVA déjà collectée  Mode de paiement";
		const endPage3 = "Taux TVA (%)";
		const pages = input.split(splitPage);
		const itemsList: Array<{ reference: number; name: string; quantity: number; unitPriceHT: number; tvaRate: number }> = [];
		for (let i=1;i<pages.length;i++){
			const pageSplit = pages[i].split(splitLine);
			for (let j = 0; j < pageSplit.length;) {
				if (pageSplit[j].includes(endPage1) || pageSplit[j].includes(endPage2) || pageSplit[j].includes(endPage3)) {
					break;
				}
				let reference;
				if (pageSplit[j].includes("  ")){
					reference = parseInt(pageSplit[j].split("  ").slice(-1)[0]);
				}
				else{
					reference = parseInt(pageSplit[j].split(" ").slice(-1)[0]);
				}
				j++;
				const splitEcoPart = pageSplit[j].split(splitNameQuantity);
				const name = splitEcoPart[0];
				let unitPriceHT;
				if (splitEcoPart.length > 1){
					unitPriceHT = parseFloat(splitEcoPart[1].replace(",","."));
				}
				else{
					j++;
					unitPriceHT = parseFloat(pageSplit[j].replace(",","."));
				}
				j++;
				if (isPrice(pageSplit[j])){ // Check if there is a discount
					unitPriceHT -= parseFloat(pageSplit[j].replace(",","."));
					j++;
				}
				const quantity = parseInt(pageSplit[j]);
				j += 2;
				const tvaRate = parseFloat(pageSplit[j].replace(',','.'));
				j++;
				itemsList.push({
					reference,
					name,
					quantity,
					unitPriceHT,
					tvaRate
				});
			}
		}
		return itemsList;
	}
</script>