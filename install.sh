pkg update -y
pkg upgrade -y
pkg install golang wget
rm $PATH/nekokit
wget https://github.com/NekoKatoriChan/NekoKit/blob/main/test/nekokit
mv nekokit $PATH/nekokit
clear
echo "Installation completed."
rm -f install.sh
