version = "n5.1.2"
srcPath = "tmp/$(version)/src"
postCheckout = ""
platform = ""

generate-flags:
	go run internal/cmd/flags/main.go

install-ffmpeg:
	rm -rf $(srcPath)
	mkdir -p $(srcPath)
	# cd $(srcPath) is necessary for windows build since otherwise git doesn't clone in the proper dir
	cd $(srcPath) && git clone https://github.com/FFmpeg/FFmpeg $(srcPath)
	cd $(srcPath) && git checkout $(version) $(postCheckout)
	cd $(srcPath) && ./configure --prefix=.. $(configure)
	cd $(srcPath) && make
	cd $(srcPath) && make install

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-platform-build:
	docker build -t astiav/$(platform) ./testdata/docker/$(platform)

test-platform-run:
	mkdir -p ./testdata/docker/$(platform)/tmp/gocache
	mkdir -p ./testdata/docker/$(platform)/tmp/gomodcache
	docker run -v .:/opt/astiav -v ./testdata/docker/$(platform)/tmp/gocache:/opt/gocache -v ./testdata/docker/$(platform)/tmp/gomodcache:/opt/gomodcache astiav/$(platform)
