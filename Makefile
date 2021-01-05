all: 
	./utils/buildproxy.sh
	./utils/buildinit.sh
	./utils/buildclient.sh
	./utils/buildserver.sh
	./utils/deployserver.sh

proxy: 
	./utils/buildproxy.sh

client: 
	./utils/buildclient.sh

server: 
	./utils/buildserver.sh

init: 
	./utils/buildinit.sh