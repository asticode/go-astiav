version = "n4.4.1"
srcPath = "tmp/$(version)/src"

generate-flags:
	go run internal/cmd/flags/main.go

install-ffmpeg:
	mkdir -p $(srcPath)
	git clone https://github.com/FFmpeg/FFmpeg $(srcPath)
	cd $(srcPath) && git checkout $(version)
	cd $(srcPath) && ./configure --prefix=../.. $(configure)
	cd $(srcPath) && make
	cd $(srcPath) && make install
