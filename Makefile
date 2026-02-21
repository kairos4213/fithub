templ:
	go tool templ generate --watch --proxy="http://localhost:8080" --open-browser=false

server:
	go tool air \
		--build.cmd "go build -o ./tmp/bin/main ." \
		--build.entrypoint "./tmp/bin/main" \
		--build.delay "1000" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

tailwind:
	npx --yes @tailwindcss/cli \
		-i "./static/css/input.css" \
		-o "./static/css/output.css" \
		--minify --watch=always

sync_static:
	go tool air \
		--build.cmd "go tool templ generate --notify-proxy" \
		--build.entrypoint true \
		--build.delay "1000" \
		--build.exclude_dir "" \
		--build.include_dir "static/css,static/js" \
		--build.include_ext "js,css"

dev:
	make -j4 tailwind templ server sync_static

build:
	go tool templ generate
	npx --yes @tailwindcss/cli -i "./static/css/input.css" -o "./static/css/output.css" --minify
	CGO_ENABLED=0 GOOS=linux go build -o fithub .
