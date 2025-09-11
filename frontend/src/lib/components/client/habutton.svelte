<script lang="ts">
    export let disabled: boolean = false;
    export let type: "button" | "submit" | "reset" | null | undefined;
    export let onclick: () => void = () => {};
</script>

<div class="HaPay">
    <button class="HaPayButton" on:click={onclick} {disabled} {type}>
        <img src="https://api.helloasso.com/v5/img/logo-ha.svg" alt="helloasso logo" class="HaPayButtonLogo" />
        <div class="HaPayButtonLabel">
            <span> Payer avec </span>
            <svg
                width="73"
                height="14"
                viewBox="0 0 73 14"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    d="M72.9992 8.78692C72.9992 11.7371 71.242 13.6283 68.4005 13.6283C65.5964 13.6283 63.8018 11.9073 63.8018 8.74909C63.8018 5.79888 65.559 3.90771 68.4005 3.90771C71.2046 3.90771 72.9992 5.64759 72.9992 8.78692ZM67.2041 8.74909C67.2041 10.5457 67.5779 11.2265 68.4005 11.2265C69.223 11.2265 69.5969 10.5079 69.5969 8.78692C69.5969 6.99031 69.223 6.30949 68.4005 6.30949C67.5779 6.30949 67.1854 7.04705 67.2041 8.74909Z"
                />
                <path
                    d="M62.978 5.08045L61.8003 6.89597C61.1647 6.47991 60.4356 6.25297 59.6692 6.23406C59.1084 6.23406 58.9214 6.40426 58.9214 6.65011C58.9214 6.9527 59.0149 7.08508 60.716 7.61461C62.4172 8.14413 63.3332 8.88169 63.3332 10.527C63.3332 12.3803 61.576 13.6474 59.1084 13.6474C57.5381 13.6474 56.0986 13.0801 55.1826 12.2101L56.7529 10.4514C57.3885 10.962 58.211 11.3402 59.0336 11.3402C59.6131 11.3402 59.9683 11.1511 59.9683 10.7918C59.9683 10.3568 59.7813 10.2622 58.2484 9.78945C56.5847 9.27883 55.65 8.31434 55.65 6.85814C55.65 5.23174 57.0333 3.92684 59.5383 3.92684C60.8656 3.90793 62.1555 4.36181 62.978 5.08045Z"
                />
                <path
                    d="M54.7358 5.08045L53.5581 6.89597C52.9225 6.47991 52.1934 6.25297 51.427 6.23406C50.8662 6.23406 50.6792 6.40426 50.6792 6.65011C50.6792 6.9527 50.7727 7.08508 52.4738 7.61461C54.175 8.14413 55.091 8.88169 55.091 10.527C55.091 12.3803 53.3338 13.6474 50.8662 13.6474C49.2959 13.6474 47.8564 13.0801 46.9404 12.2101L48.5107 10.4514C49.1463 10.962 49.9689 11.3402 50.7914 11.3402C51.3709 11.3402 51.7261 11.1511 51.7261 10.7918C51.7261 10.3568 51.5391 10.2622 50.0062 9.78945C48.3238 9.27883 47.4078 8.31434 47.4078 6.85814C47.4078 5.23174 48.7911 3.92684 51.2961 3.92684C52.6234 3.90793 53.9133 4.36181 54.7358 5.08045Z"
                />
                <path
                    d="M46.7721 11.4156L46.0991 13.5526C44.9401 13.477 44.1923 13.1555 43.6876 12.3045C43.0333 13.3068 42.0051 13.6283 40.9956 13.6283C39.201 13.6283 38.042 12.418 38.042 10.7537C38.042 8.74909 39.5375 7.65222 42.3603 7.65222H42.9959V7.42528C42.9959 6.51752 42.6968 6.27167 41.706 6.27167C40.9209 6.30949 40.1357 6.4797 39.4067 6.74446L38.6963 4.62636C39.8179 4.17248 41.0143 3.94554 42.2294 3.90771C45.0709 3.90771 46.23 5.00459 46.23 7.23616V10.3566C46.23 10.9996 46.3795 11.2643 46.7721 11.4156ZM43.0146 10.7348V9.39209H42.6594C41.7247 9.39209 41.2947 9.71359 41.2947 10.4133C41.2947 10.9239 41.5752 11.2643 42.0238 11.2643C42.4164 11.2643 42.7903 11.0563 43.0146 10.7348Z"
                />
                <path
                    d="M37.5363 8.78692C37.5363 11.7371 35.7791 13.6283 32.9376 13.6283C30.1335 13.6283 28.3389 11.9073 28.3389 8.74909C28.3389 5.79888 30.0961 3.90771 32.9376 3.90771C35.7417 3.90771 37.5363 5.64759 37.5363 8.78692ZM31.7412 8.74909C31.7412 10.5457 32.1151 11.2265 32.9376 11.2265C33.7601 11.2265 34.134 10.5079 34.134 8.78692C34.134 6.99031 33.7601 6.30949 32.9376 6.30949C32.1151 6.30949 31.7225 7.04705 31.7412 8.74909Z"
                />
                <path
                    d="M23.8154 10.6972V0.692948L27.1243 0.352539V10.527C27.1243 10.8296 27.2551 10.9809 27.5355 10.9809C27.6477 10.9809 27.7786 10.962 27.8907 10.9052L28.4889 13.2881C27.8907 13.4961 27.2738 13.6096 26.6569 13.5907C24.8249 13.6285 23.8154 12.5505 23.8154 10.6972Z"
                />
                <path
                    d="M18.8057 10.6972V0.692948L22.1145 0.352539V10.527C22.1145 10.8296 22.2454 10.9809 22.5071 10.9809C22.6192 10.9809 22.7501 10.962 22.8623 10.9052L23.4418 13.2881C22.8436 13.4961 22.2267 13.6096 21.6098 13.5907C19.8151 13.6285 18.8057 12.5505 18.8057 10.6972Z"
                />
                <path
                    d="M17.9071 9.71359H12.4859C12.6728 11.0185 13.3084 11.2454 14.2805 11.2454C14.9161 11.2454 15.533 10.9807 16.2994 10.4511L17.6454 12.2856C16.6172 13.1555 15.3087 13.6283 13.9627 13.6283C10.6912 13.6283 9.13965 11.5858 9.13965 8.78692C9.13965 6.13929 10.6352 3.90771 13.5888 3.90771C16.2247 3.90771 17.9632 5.60976 17.9632 8.63562C17.9819 8.93821 17.9445 9.39209 17.9071 9.71359ZM14.7291 7.70895C14.7105 6.80119 14.5235 6.04473 13.6823 6.04473C12.9719 6.04473 12.6167 6.46079 12.4859 7.84134H14.7291V7.70895Z"
                />
                <path
                    d="M8.24307 6.61229V13.2692H4.93423V7.21746C4.93423 6.49882 4.7286 6.32862 4.4295 6.32862C4.07431 6.32862 3.70043 6.61229 3.30786 7.21746V13.2503H-0.000976562V0.692948L3.30786 0.352539V5.06154C4.07431 4.24834 4.82207 3.90793 5.83154 3.90793C7.32706 3.90793 8.24307 4.89133 8.24307 6.61229Z"
                />
            </svg>
        </div>
    </button>
    <div class="HaPaySecured">
        <svg width="9" height="10" viewBox="0 0 11 12" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
                d="M3.875 3V4.5H7.625V3C7.625 1.96875 6.78125 1.125 5.75 1.125C4.69531 1.125 3.875 1.96875 3.875 3ZM2.75 4.5V3C2.75 1.35938 4.08594 0 5.75 0C7.39062 0 8.75 1.35938 8.75 3V4.5H9.5C10.3203 4.5 11 5.17969 11 6V10.5C11 11.3438 10.3203 12 9.5 12H2C1.15625 12 0.5 11.3438 0.5 10.5V6C0.5 5.17969 1.15625 4.5 2 4.5H2.75ZM1.625 6V10.5C1.625 10.7109 1.78906 10.875 2 10.875H9.5C9.6875 10.875 9.875 10.7109 9.875 10.5V6C9.875 5.8125 9.6875 5.625 9.5 5.625H2C1.78906 5.625 1.625 5.8125 1.625 6Z"
            />
        </svg>
        <span>Paiement sécurisé</span>
        <img
            src="https://helloassodocumentsprod.blob.core.windows.net/public-documents/bouton_payer_avec_helloasso/logo-visa.svg"
            alt="Logo Visa"
        />
        <img
            src="https://helloassodocumentsprod.blob.core.windows.net/public-documents/bouton_payer_avec_helloasso/logo-mastercard.svg"
            alt="Logo Mastercard"
        />
        <img
            src="https://helloassodocumentsprod.blob.core.windows.net/public-documents/bouton_payer_avec_helloasso/logo-cb.svg"
            alt="Logo CB"
        />
        <img
            src="https://helloassodocumentsprod.blob.core.windows.net/public-documents/bouton_payer_avec_helloasso/logo-pci.svg"
            alt="Logo PCI"
        />
    </div>
</div>

<style>
	.HaPay {
		width: fit-content;
		display: -webkit-box;
		display: -ms-flexbox;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		-webkit-box-pack: center;
		-ms-flex-pack: center;
	}

	.HaPay * {
		font-family: 'Open Sans', 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans',
			Arial, sans-serif;
		transition: all 0.3s ease-out;
	}

	.HaPayButton {
		align-items: stretch;
		-webkit-box-pack: stretch;
		-ms-flex-pack: stretch;
		background: none;
		border: none;
		display: -webkit-box;
		display: -ms-flexbox;
		display: flex;
		padding: 0;
		border-radius: 8px;
	}

	.HaPayButton:hover {
		cursor: pointer;
	}

	.HaPayButton:not(:disabled):focus {
		box-shadow: 0 0 0 0.25rem rgba(73, 211, 138, 0.25);
		-webkit-box-shadow: 0 0 0 0.25rem rgba(73, 211, 138, 0.25);
	}

	.HaPayButton:not(:disabled):hover .HaPayButtonLabel,
	.HaPayButton:not(:disabled):focus .HaPayButtonLabel {
		background-color: #483dbe;
	}

	.HaPayButton:not(:disabled):hover .HaPayButtonLogo,
	.HaPayButton:not(:disabled):focus .HaPayButtonLogo,
	.HaPayButton:not(:disabled):hover .HaPayButtonLabel,
	.HaPayButton:not(:disabled):focus .HaPayButtonLabel {
		border: 1px solid #483dbe;
	}

	.HaPayButton:disabled {
		cursor: not-allowed;
	}

	.HaPayButton:disabled .HaPayButtonLogo,
	.HaPayButton:disabled .HaPayButtonLabel {
		border: 1px solid #d1d6de;
	}

	.HaPayButtonLogo {
		background-color: #ffffff;
		border: 1px solid #4c40cf;
		border-top-left-radius: 8px;
		border-bottom-left-radius: 8px;
		padding: 4px 12px;
		width: 56px;
	}

	.HaPayButtonLabel {
		align-items: center;
		-webkit-box-pack: center;
		-ms-flex-pack: center;
		justify-content: space-between;
		column-gap: 5px;
		background-color: #4c40cf;
		border: 1px solid #4c40cf;
		border-top-right-radius: 8px;
		border-bottom-right-radius: 8px;
		color: #ffffff;
		font-size: 16px;
		font-weight: 800;
		display: -webkit-box;
		display: -ms-flexbox;
		display: flex;
		padding: 16px;
	}

	.HaPayButton:disabled .HaPayButtonLabel {
		background-color: #d1d6de;
		color: #505870;
	}

	.HaPaySecured {
		align-items: center;
		-webkit-box-pack: center;
		-ms-flex-pack: center;
		justify-content: space-between;
		display: -webkit-box;
		display: -ms-flexbox;
		display: flex;
		column-gap: 5px;
		padding: 8px 16px;
		font-size: 12px;
		font-weight: 600;
		color: #2e2f5e;
	}

	.HaPay svg {
		fill: currentColor;
	}
</style>
