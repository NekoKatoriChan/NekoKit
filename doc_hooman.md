# NekoKit
This document contains the basics you need to get started.

## What's this?
NekoKit is a small transpiler that converts a custom scripting language (.nk) into Go code. It is designed for building simple terminal-based games quickly, with a cleaner and more expressive syntax than raw Go.
The tool parses user-defined scripts and generates valid Go code, which can be executed directly or compiled into standalone binaries. It supports variables, input/output, control flow, file operations, and basic UI elements for terminal applications.

## Installation

```bash
curl -L https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh
```

Build from sources:

```bash
apt update
apt install git golang
git clone https://github.com/NekoKatoriChan/NekoKit.git
cd NekoKit
go build -o nekokit
```

## Usage

### Run the code

```bash
nekokit myGame.nk
```
Run the code instant.

### Build the code

```bash
nekokit myGame.nk --build
```

Builds the script into a standalone binary that can run without NekoKit.

### Set output name

```bash
nekokit myGame.nk --build --output superCoolGame
```


### Update (WIP)

```bash
nekokit --update
```

### Verbose
```bash
nekokit file.nk -v
```

Shows detailed transpilation process.
---


#### Give

```
give x=10
give y=20+5
give name="whiskers"
give health=maxHealth
```

Set variable.
#### write & writeln 

```
write "hello "
write "hooman"
writeln "meow meow!"
```

Print output.
#### read 

```
read playerName
```

It will wait for an input.

#### read -p - read with a prompt

```
read -p "what's yur name? " playerName
```

asks a question first, then wait for input.


#### gameloop

```
gameloop start
    read -p "tap 1 to start: " aha
    if aha == "1" {
        call Game1    
    }
gameloop end
```

It will run forever until you stop.

#### Block call

```
Game1 start
    dialog "welcome to Game1!"
    give health=100
Game1 end
```


Note: Blocks cannot be nested.

#### callonce 

```
callonce BossIntro
```

Call block once

### String interpolation

```
writeln "hello, $name!"
writeln "you have $gold coins"
```


#### clear 

```
clear
```

wipe everything.

#### border - pretty boxes to sit in

```
border top
border mid
border bot
```

draws ascii art borders! cuz if it fits, i sits!
- `top` and `bot`: `═════════════════════════` (thicc edge!)
- `mid`: `─────────────────────────` (thin divider!)

#### dialog, prompt, menu - fancy meows!

```
dialog "da cat says: meow meow!"
prompt "choose yur destiny: "
menu fight, defend, run away
```

- `dialog` prints fancy text inside a cinematic box!
- `prompt` puts a cute ` > ` arrow before yur text!
- `menu` takes a comma-separated list and automatically turns it into a numbered list `(1) fight, (2) defend`! 


#### stat 

```
stat "health" $playerHP
stat level $currentLevel
```

formats and prints stats like: `health: 100` or `level: 5`.

#### score & level

```
score points 10
level bossStage 5
```

#### damage & heal

```
damage playerHP 10
heal playerHP 25
```


#### reset 

```
reset score
```


### random 

```
random roll 6
```

generates random number from 0 to 5 (max-1). modern Go seeds it automatically.

### burying n' digging up toys (files) 💾

#### load

```
load ~/saveGame.txt
load myData ~/custom_path/file.txt
```

#### save

```
save gameData ~/saveGame.txt
```


#### peek

```
peek ~/saveGame.txt {
    writeln "File found! *happy meow*"
} else {
    writeln "No file here, nya!"
}
```

This will searching a file if it exist.

#### create - touch a files

```
create fish.txt
```

Create an empty file.

#### if/else 

```
if playerHP <= 0 {
  writeln "game over..."
  susu
} else {
  writeln "keep fighting!"
}
```


#### susu - Exit 0

```
susu
```

#### run - doing termux magics

```
run clear && ls
```

#### web request 

```
meow request [options] <url>

Options:
  hiss <METHOD>        → HTTP method (GET, POST, PUT, DELETE, etc)
  lick <HEADER>        → Add header "Key: Value"
  spit <DATA>          → Request body (JSON, form data, etc)
  grab <FILE>          → Download response to file
  sniff <HEADER>       → Read-only header check (like -i)
  yowl                 → Verbose output (show headers, body, etc)
  purr                 → Silent mode (no output except response)
  scratch <N>          → Retry N times on failure
  tail                 → Follow redirects (301, 302, etc)

```

curl but cat
---

## purrfect kitty quest! (example game) 🎯

here's a smol adventure game to show it all working together:

```
clear
border top
writeln "🐱 KITTY QUEST 🐱"
border bot

read -p "enter yur cat name: " playerName
give playerHP=100
score gold 0

gameloop start
  clear
  border top
  stat "name" $playerName
  stat "health" $playerHP  
  stat "gold" $gold
  border bot
  writeln ""
  
  dialog "yu encounter a wild mouse!"
  
  menu fight, run away
  prompt "choice: "
  read choice
  
  if choice == "1" {
    random dmg 15
    damage playerHP $dmg
    writeln ""
    dialog "mouse attacks! took $dmg damage! *hiss!*"
    
    if playerHP <= 0 {
      writeln ""
      writeln "💀 GAME OVER 💀"
      susu
    }
    
    score gold 10
    writeln "yu won! found 10 gold~ *purr*"
  } else {
    writeln ""
    dialog "yu ran away safely! coward hooman..."
  }
  
  writeln ""
  prompt "press enter to continue..."
  read dummy
gameloop end
```

---

## nerdy cat stuff 🤓

### how it works

1. Take your `.nk` file.
2. Transpile into go code.
3. Run it with go run or go build.
4. Clean temporary file.

### da transpiler features:

- **Termux Native**: Optimized paths (`~/` expands to Termux `$HOME`), standard `sh` execution, and raw ANSI clearing.
- **Strict Go Compliance**: Bypasses da annoying "declared and not used" errors by generating `_ = varName` safety nets! 
- **Separate Bowls (Block Scoping)**: Named game bloks act as isolated Go functions. Main variables and blok variables don't mix! Keep da food bowls separate!
- **Callonce Tracker**: Uses a lightning-fast `map[string]bool` to remember what VIP bloks have been sniffed already.

## license & credits ✨

made wif love by cats, for cats (and cat-loving hoomans!)

Happy coding.
