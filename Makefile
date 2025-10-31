templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

server:
	go run github.com/air-verse/air@v1.61.7 \
		--build.cmd "go build -o tmp/bin/main" \
		--build.bin "tmp/bin/main" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

tailwind:
	npx @tailwindcss/cli -i "./static/css/input.css" -o "./static/css/output.css" --minify --watch

sync_static:
	go run github.com/air-verse/air@v1.61.7 \
		--build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "static" \
		--build.include_ext "js,css"

dev:
	make -j4 templ server tailwind sync_static
