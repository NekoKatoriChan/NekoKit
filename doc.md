
# nekokit scribbles 🐱

*nya nya~ welcome to da official nekokit docs! made wif 100% premium catnip and fresh tuna! prrr~*

## da heck is dis??

nekokit is a smol transpiler dat takes cute little game scripts (`.nk` files) and turns dem into big, strict, adult go code! it's liek teaching a kitten to meow in hooman language! 

yu write silly smol commands, and nekokit translates dem into boring go code dat actually runs in yur **Termux** scratching post! purrfect for making text games wifout da grumpy hooman syntax! *swish!*

## gimmee! (instawwation)

```bash
curl -L https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh
```

or if yu iz a brave hunter chasing da red dot and wanna build from da source:

```bash
apt update
apt install git golang
git clone https://github.com/NekoKatoriChan/NekoKit.git
cd NekoKit
go build -o nekokit
```

## how to play (usage)

### run da thing (no wait, jus go!)

```bash
nekokit myGame.nk
```

dis runs yur script right meow! *pounce!*

### building a toy (to share wif frens!)

```bash
nekokit myGame.nk --build
```

makes a standalone binary! nao yur game can run wifout nekokit installed. very sneaky cat!

### give it a cute name (cuz "myGame" is boring)

```bash
nekokit myGame.nk --build --output superCoolGame
```

### update (gitting fresh tuna!)

```bash
nekokit --update
```
pulls da freshest version from da interwebs. automatic catnip delivery! *mlem*

### verbose (sniffing ALL da details)
```bash
nekokit file.nk -v
```
see exactly how da cat translates yur code behind da scenes! *sniff sniff sniff!*

---

## cat language dictionary 🐾

### toys n' keeping dem (variables)

#### give - polite giving!

```
give x=10
give y=20+5
give name="whiskers"
give health=maxHealth
```

da `give` command is liek offering a gift to a variable! much nicer dan jus `=` all alone!
*(nerd note: da cat compiler automatically protects yu from "unused variable" panics by adding `_ = varName`! we keeps da Go compiler well-fed and happy! nya!)*

### meowing at hoomans (in/out)

#### write & writeln - meow loudly!

```
write "hello "
write "hooman"
writeln "meow meow!"
```

`write` outputs on da same line (liek cats walking in a straight line). `writeln` gives it a fresh new line!

#### read - listening to hooman

```
read playerName
```

waits for da hooman to type sumthing and stores it in da bowl.

#### read -p - read wif a prompt 

```
read -p "what's yur name? " playerName
```

asks a question first, den waits for input. so polite! *purr*

### chasing tails & calling frens! 🔁

#### gameloop - chasing tail forever!

```
gameloop start
    read -p "tap 1 to start: " aha
    if aha == "1" {
        call Game1    
    }
gameloop end
```

runs forever until yu stop it! purrfect for games dat keep going and going!

#### calling yur frens! (named bloks)

```
Game1 start
    dialog "welcome to Game1!"
    give health=100
Game1 end
```

yu can make yur own bloks of code and `call Game1` to run dem! dey have deir own clean litter boxes (variable scopes), so dey won't mess wif yur main toys! do NOT mix da kibbles!

💣 cute note from hooman: don't use block inside block, they will explode!

#### callonce - VIP cats ONLY!

```
callonce BossIntro
```

calls a blok, but **only once per game**! if yu try to call it again, da cat just ignores yu and licks a paw. powered by a super-fast tracking map! *mew!*

### magic $ treats (string interpolation)

```
writeln "hello, $name!"
writeln "yu have $gold coins"
```

da `$variable` syntax gets replaced wif actual values! it's liek hiding a treat inside a toy~

### fancy UI tricks 🎮

#### clear - swiping everything off da table!

```
clear
```

makes everything disappear! *crash!* heavily optimized wif native ANSI sequences for Termux!

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

### stats and counting kibbles 📊

#### stat - showing off yur stuff

```
stat "health" $playerHP
stat level $currentLevel
```

formats and prints stats liek: `health: 100` or `level: 5`.

#### score & level - math made easy!

```
score points 10
level bossStage 5
```

first time creates da variable, after dat `score` ADDS to it (`+=`), while `level` SETS it (`=`). smart kitty! 

#### damage & heal - bites and naps

```
damage playerHP 10
heal playerHP 25
```

subtracts or adds to a value! ouchie! and den nom nom healing herbs~ 

#### reset - dropping it in da litter box

```
reset score
```

back to zero! fresh start liek a new day~

### random - knocking random things over 🎲

```
random roll 6
```

generates random number from 0 to 5 (max-1). modern Go seeds it automatically, so no need to wake da cat up to do it!

### burying n' digging up toys (files) 💾

#### load - dig up a toy

```
load ~/saveGame.txt
load myData ~/custom_path/file.txt
```

reads a file! if yu jus say `load ~/my-save.txt`, da cat will automatically name da variable `my_save` (replacing dashes wif underscores so Go doesn't hiss at yu!). Or yu can explicitly name it wif `load myData [path]`!

#### save - bury a toy for later

```
save gameData ~/saveGame.txt
```

writes yur variable contents to a file. hiding it from da dog!

#### peek - sniffing for files

```
peek ~/saveGame.txt {
    writeln "File found! *happy meow*"
} else {
    writeln "No file here, nya!"
}
```

da cat sniffs to see if a file exists before yu try to open it!

#### create - touch a files

```
create fish.txt
```

give some fish! but's it's empty
### choosing which bowl to eat from 🔀

#### if/else - decisions decisions

```
if playerHP <= 0 {
  writeln "game over nya..."
  susu
} else {
  writeln "keep fighting!"
}
```

standard if/else blocks! 

#### susu - nap time!

```
susu
```
ends everything! time for a big sleepy nap~ 😴 (saying bye-bye to da program!)

#### run - doing termux magics

```
run clear && ls
```

runs standard shell commands! strictly mapped to `sh` so it plays purrfectly wif Termux environments! no quotes required, da nerd cat already escaped dem for yu! 😹

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

### how it works internally (da braincells)

1. sniffs yur `.nk` file.
2. transpiles it into go code using cute pattern matching.
3. either runs it wif `go run` or builds it wif `go build`.
4. cleans up temporary files liek a good kitty cleaning its paws~

### da transpiler features:

- **Termux Native**: Optimized paths (`~/` expands to Termux `$HOME`), standard `sh` execution, and raw ANSI clearing. loves da terminal box!
- **Strict Go Compliance**: Bypasses da annoying "declared and not used" errors by generating `_ = varName` safety nets! 
- **Separate Bowls (Block Scoping)**: Named game bloks act as isolated Go functions. Main variables and blok variables don't mix! Keep da food bowls separate!
- **Callonce Tracker**: Uses a lightning-fast `map[string]bool` to remember what VIP bloks have been sniffed already.

## license & credits ✨

made wif love by cats, for cats (and cat-loving hoomans!)

*meow meow, happy coding!* 🐾
