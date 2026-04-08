# NekoKit

a silly little transpiler for making text-based games without the boring stuff. write simple commands, get go code. purrfect for termux! (currently take a nap with catnip, may not code in someday) (⁠ㆁ⁠ω⁠ㆁ⁠)

## install

```bash
curl -L https://github.com/NekoKatoriChan/NekoKit/raw/refs/heads/stable/install | sh
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

[doc](doc.md)
---

*made with love and tuna* 🐟
