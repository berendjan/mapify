# in app/js, enable tailwindcss
npx tailwindcss -i ./tailwind.css -o ../go/static/tailwind.css --watch

# in app/go, enable hot reloading
./hotreload.sh

