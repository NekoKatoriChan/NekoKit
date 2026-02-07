# nekokit documentation 🐱

*nya nya~ welcome to the official nekokit docs, made with tuna*

## what the hell is this?

nekokit is a silly little transpiler that turns cute game script files (`.nk` files) into proper go code! it's like teaching a kitten to speak human language, but for programming~ 

basically you write simple commands, and nekokit translates them into boring adult go code that actually runs. purrfect for making text-based games without all the fussy syntax!

## installation 

```bash
curl -fsSL https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh
```

or if you're feeling brave and want to build from source like a real adventurer cat:

```bash
git clone https://github.com/NekoKatoriChan/NekoKit.git
cd nekokit
go build -o nekokit
```

## usage

### running a script (no build, just go go go!)

```bash
nekokit myGame.nk
```

this runs your script immediately! like pressing the turbo button on a cat

### building an executable (for when you wanna share with friends!)

```bash
nekokit myGame.nk --build
```

makes a standalone binary! now your game can run without nekokit installed. very cat

### custom output name (because "myGame" is boring)

```bash
nekokit myGame.nk --build --output superCoolGame
```

### updating nekokit (get the freshest fish!)

```bash
nekokit --update
```

pulls the latest version from the interwebs and installs it. automatic catnip delivery! 

## language reference 

### variables and assignment

#### give - assigning values (so polite!)

```
give x=10
give y=20+5
give name="whiskers"
give health=maxHealth
```

the `give` command is like offering a gift to a variable! much nicer than just `=` all alone~

supports:
- numbers: `give score=100`
- expressions: `give total=score+bonus`
- strings: `give greeting="hello!"`
- variables: `give currentHP=maxHP`

### input/output (talking to humans!)

#### write - print without newline

```
write "hello "
write "world"
```

outputs: `hello world` (all on same line, like cats walking in a line)

#### writeln - print with newline

```
writeln "meow meow"
writeln "nya nya"
```

each one gets its own line!

#### read - get input from user

```
read playerName
```

waits for the hooman to type something and stores it.

#### read -p - read with a prompt 

```
read -p "what's your name? " playerName
```

asks a question first, then waits for input.

### string interpolation (magic $dollar signs!)

```
writeln "hello, $name!"
writeln "you have $gold coins"
```

the `$variable` syntax gets replaced with actual values! it's like mail merge but for cats~

### game-specific commands 

#### gameloop - the main game loop!

```
gameloop start
  clear
  writeln "still playing..."
  writeln "press ctrl+c to stop nya~"
gameloop end
```

runs forever until you stop it! perfect for games that keep going and going like an energetic kitten~

#### clear - clear the screen

```
clear
```

makes everything disappear! *poof!* 

#### border - draw pretty boxes

```
border simple
border double  
border thick
```

draws ascii art borders! makes your game look all professional-like:
- `simple`: `┌─┐` style (clean and minimal)
- `double`: `╔═╗` style (fancy double lines!)
- `thick`: `▓▓▓` style (EXTRA THICC)

#### dialog - character speech

```
dialog "the cat says: meow meow!"
dialog "hero: i must save the kingdom!"
```

adds `>>>` before text to make it look like important dialog. very dramatic! 🎭

#### prompt - visual input prompt

```
prompt "choose your destiny: "
```

shows ` > ` before the text. makes choices feel official~

#### menu - quick menu display

```
menu options
```

displays a simple numbered menu:
```
[1] Option 1
[2] Option 2  
[3] Option 3
```

(note: you still need to handle the input yourself! this just shows the pretty menu)

### stats and game values 📊

#### stat - display a stat nicely

```
stat "health" $playerHP
stat level $currentLevel
```

formats and prints stats like: `health: 100` or `level: 5`

works with quoted or unquoted names!

#### score - add to score (with auto-declaration!)

```
score points 10
score points 5
```

first time creates the variable, after that it adds to it. smart kitty! 

#### level - set level value

```
level playerLevel 1
level playerLevel 5
```

similar to score but uses `=` instead of `+=` for setting exact values~

#### damage - subtract from health

```
damage playerHP 10
damage bossHP $attackPower
```

reduces a value. ouch! 

#### heal - restore health  

```
heal playerHP 25
heal playerHP $potionPower
```

adds back health! nom nom healing herbs~ 

#### reset - set to zero

```
reset score
```

back to nothing! fresh start like a new day~

### inventory system 

#### inventory add - put item in inventory

```
inventory bag add "sword"
inventory bag add $foundItem
```

adds items to your collection!

#### inventory remove - take item out

```
inventory bag remove "sword"
```

removes items.

### random numbers 🎲

```
random roll 6
random damage 20
```

generates random number from 0 to (max-1). perfect for dice rolls and critical hits!

requires the number to be a max value, so `random roll 6` gives you 0-5 (like a d6!)

### file operations 

#### load - read file contents

```
load ~/saveGame.txt
load config.ini
```

reads file into a variable named after the file (without extension). 

`load saveGame.txt` creates variable `saveGame` with file contents!

#### save - write to file  

```
save gameData ~/saveGame.txt
save config settings.ini
```

writes variable contents to a file. persistence!

### control flow 🔀

#### if/else - decisions decisions

```
if playerHP <= 0 {
  writeln "game over nya..."
  susu
}

if hasKey == true {
  writeln "door unlocked!"
} else {
  writeln "you need a key!"
}
```

standard if/else blocks! just like go but triggered by our silly syntax~

#### susu - exit program

```
susu
```

ends everything! time for a nap~ 😴


## example game 🎯

here's a tiny adventure game to show it all working together:

```
clear
border double

writeln "🐱 KITTY QUEST 🐱"
writeln ""

read -p "enter your cat name: " playerName
give playerHP=100
give gold=0

gameloop start
  clear
  writeln "========== ADVENTURE =========="
  stat "name" $playerName
  stat "health" $playerHP  
  stat "gold" $gold
  writeln "==============================="
  writeln ""
  
  dialog "you encounter a wild mouse!"
  
  writeln "[1] fight"
  writeln "[2] run away"
  prompt "choice: "
  read choice
  
  if choice == "1" {
    random damage 15
    damage playerHP $damage
    writeln ""
    dialog "mouse attacks! took $damage damage!"
    
    if playerHP <= 0 {
      writeln ""
      writeln "💀 GAME OVER 💀"
      susu
    }
    
    score gold 10
    writeln "you won! found 10 gold~"
  } else {
    writeln ""
    dialog "you ran away safely! coward..."
  }
  
  writeln ""
  prompt "press enter to continue..."
  read dummy
gameloop end
```

## technical details for nerds 🤓

### how it works internally

1. reads your `.nk` file
2. transpiles it into go code using pattern matching
3. either runs it with `go run` or builds it with `go build`
4. cleans up temporary files like a good kitty cleaning its paws~

### the transpiler does:

- string interpolation (`$var` → `fmt.Sprintf`)
- import detection (only imports what you use!)
- variable declaration tracking (uses `:=` first time, `=` after)
- indent management for proper go formatting
- path expansion (`~/` → actual home directory)

### files in this project

- `main.go` - CLI argument parsing and file handling
- `transpile.go` - the big transpiler brain! 🧠
- `run.go` - running and building utilities

## contributing 

found a bug? want to add a feature? 

purrfect! just remember:
- you are a cat
- keep it simple and silly
- write cat-style comments (lowercase, playful, helpful!)
- test your changes with actual `.nk` files

## why does this exist? 🤔

why not?

plus it's fun to have a programming language that says "nya~" 

## license & credits ✨

made with love by cats, for cats (and cat-loving humans!)

*meow meow, happy coding!* 🐾

---

*this documentation was written at 3am fueled by catnip tea and the desire to make programming more adorable*
