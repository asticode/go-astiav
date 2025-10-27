version=
srcPath=
patchPath=

install-ffmpeg:
	rm -rf $(srcPath)
	mkdir -p $(srcPath)
	cd $(srcPath) && git clone https://github.com/FFmpeg/FFmpeg .
	cd $(srcPath) && git checkout $(version)
ifneq "" "$(patchPath)"
	cd $(srcPath) && git apply $(patchPath)
endif
	cd $(srcPath) && ./configure --prefix=..
	cd $(srcPath) && make
	cd $(srcPath) && make install