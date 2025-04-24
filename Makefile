live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

live/server:
	go run github.com/air-verse/air@v1.61.7 \
		--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

live/sync_assets:
	go run github.com/air-verse/air@v1.61.7 \
		--build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "static" \
		--build.include_ext "js,css"

live:
	make -j3 live/templ live/server live/sync_assets
