version = "n5.1.2"
srcPath = "tmp/$(version)/src"
postCheckout = ""

generate-flags:
	go run internal/cmd/flags/main.go

install-ffmpeg:
	rm -rf $(srcPath)
	mkdir -p $(srcPath)	
	# cd $(srcPath) prepend to the next command is necessary for windows build since otherwise git doesn't clone in the proper dir
	git clone --depth 1 --branch $(version) https://github.com/FFmpeg/FFmpeg $(srcPath)
	cd $(srcPath) && ./configure --prefix=.. $(configure)
	cd $(srcPath) && make
	cd $(srcPath) && make install
