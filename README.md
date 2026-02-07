# NekoKit

a silly little transpiler for making text-based games without the boring stuff. write simple commands, get go code. purrfect for termux!

## install

```bash
curl -fsSL https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh
```

## quick start

```bash
# run a game script
nekokit game.nk

# build an executable  
nekokit game.nk --build

# update nekokit
nekokit --update
```

## example

```
clear
writeln "🐱 hello hooman!"
read -p "what's your name? " name
writeln "nice to meet you, $name!"
```


---

*made with love and tuna* 🐟
