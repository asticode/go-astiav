version = "n5.1.2"
srcPath = "tmp/$(version)/src"
postCheckout = ""

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
