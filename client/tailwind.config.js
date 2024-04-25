
/** @type {import('tailwindcss').Config} */
const defaultTheme = require('tailwindcss/defaultTheme');
const { iconsPlugin, getIconCollections } = require("@egoist/tailwindcss-icons");

export default {
    content: [
        './src/**/*.{html,js,svelte,ts}',
        './index.html',
    ],
    theme: {
        extend: {
            colors: {
                background: '#0D1117',
                primary: '#7371FF',
                light: '#8C8E92',
                secondary: '#22F1BF',
                secodaryLight: '#FAFAFA',
                cardBorder: '#756D6D'
            },
            fontFamily: {
                sans: ['Inter var', ...defaultTheme.fontFamily.sans],
            },
            screens: {
                '3xl': '1700px'
            }
        },
    },
    plugins: [
        iconsPlugin({
            collections: getIconCollections(["mdi", "lucide", "logos"])
        })
    ],
}

