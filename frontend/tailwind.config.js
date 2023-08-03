/** @type {import('tailwindcss').Config} */
export default {
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		'node_modules/preline/dist/*.js',
		'node_modules/preline/dist/*.js'
	],
	theme: {
		extend: {}
	},
	plugins: [require('preline/plugin')]
};
