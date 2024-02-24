/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "../go/**/*.{html,js,templ}"
    ],
    theme: {
        extend: {},
    },
    plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
