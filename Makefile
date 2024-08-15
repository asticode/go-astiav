version=n7.0
srcPath=tmp/$(version)/src
patchPath=
platform=

generate-flags:
	go run internal/cmd/flags/main.go

install-ffmpeg:
	rm -rf $(srcPath)
	mkdir -p $(srcPath)
	cd $(srcPath) && git clone https://github.com/FFmpeg/FFmpeg .
	cd $(srcPath) && git checkout $(version)
ifneq "" "$(patchPath)"
	cd $(srcPath) && git apply $(patchPath)
endif
	cd $(srcPath) && ./configure --prefix=.. $(configure)
	cd $(srcPath) && make
	cd $(srcPath) && make install

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-platform-build:
	docker build -t astiav/$(platform) ./internal/test/$(platform)

test-platform-run:
	mkdir -p ./internal/test/$(platform)/tmp/gocache
	mkdir -p ./internal/test/$(platform)/tmp/gomodcache
	docker run -v .:/opt/astiav -v ./internal/test/$(platform)/tmp/gocache:/opt/gocache -v ./internal/test/$(platform)/tmp/gomodcache:/opt/gomodcache astiav/$(platform)
