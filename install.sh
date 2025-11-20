clear
	apt update -y
	apt install wget golang -y
	wget https://github.com/NekoKatoriChan/NekoKit/raw/refs/heads/main/test/nekokit
	chmod +x nekokit
	mv nekokit $PREFIX/bin/
	clear
	echo "Installation completed. Usage: nekokit <file>.nk."
	sleep 3
	rm install.sh
