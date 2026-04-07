pkg update -y
pkg upgrade -y
pkg install golang wget
wget https://github.com/NekoKatoriChan/NekoKit/raw/refs/heads/stable/test/nekokit -O /data/data/com.termux/files/usr/bin/nekokit
chmod +x /data/data/com.termux/files/usr/bin/nekokit
clear
echo "Installation completed."
rm -f install.sh
