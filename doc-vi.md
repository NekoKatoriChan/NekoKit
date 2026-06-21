# Sổ tay NekoKit 🐱

*meo meo~ chào mừng sen đã đến với tài liệu chính thức của NekoKit! Được viết bằng 100% cỏ mèo thượng hạng và cá ngừ tươi ngon! hừ hừ~*

## oát đờ phắc

NekoKit là một bộ dịch mã siêu nhỏ (transpiler) giúp biến các đoạn mã kịch bản game siêu đáng yêu (file `.nk`) thành mã nguồn Go người lớn, nghiêm túc và chuẩn chỉ! Giống như dạy một bé mèo con nói tiếng người vậy đó! 

Bạn chỉ cần viết các câu lệnh ngắn gọn, dễ thương, và NekoKit sẽ dịch chúng sang ngôn ngữ Go "tẻ nhạt" để chạy trực tiếp trên **bàn cào móng Termux** của bạn! Quá hoàn hảo để tạo các trò chơi văn bản (text game) mà không cần phải nhức đầu với mớ cú pháp phức tạp của loài người! *vẫy đuôi!*

## Cho xin đi! (Cài đặt)

```bash
curl -L https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh
```

Hoặc nếu bạn là một thợ săn dũng cảm thích đuổi theo chấm đỏ laser và muốn tự tay biên dịch từ nguồn:

```bash
apt update
apt install git golang
git clone https://github.com/NekoKatoriChan/NekoKit.git
cd NekoKit
go build -o nekokit
```

## Cách chơi (Sử dụng)

### Chạy trực tiếp (Không cần chờ, đi thôi!)

```bash
nekokit myGame.nk
```

Lệnh này sẽ chạy kịch bản của bạn ngay và luôn! *vồ lấy!*

### Chế tạo đồ chơi (Để chia sẻ với bạn bè!)

```bash
nekokit myGame.nk --build
```

Lệnh này sẽ tạo ra một file thực thi độc lập! Giờ thì game của bạn có thể chạy ở bất cứ đâu mà không cần cài sẵn NekoKit. Mèo thật là tinh quái!

### Đặt cho nó một cái tên dễ thương (Vì "myGame" nghe chán vl)

```bash
nekokit myGame.nk --build --output superCoolGame
```

### Cập nhật (Nhận cá ngừ tươi mới nhất!)

```bash
nekokit --update
```
Tự động tải phiên bản mới nhất từ internet về máy. Hệ thống giao cỏ mèo tự động hoạt động! *mlem* (không chắc chắn chạy, đừng thử)

### Verbose (Ngửi sạch mọi chi tiết)
```bash
nekokit file.nk -v
```
Xem chính xác cách Boss mèo dịch mã nguồn của bạn ở hậu trường như thế nào! *ngửi ngửi ngửi!*

---

## Từ điển tiếng mèo 🐾

### Đồ chơi và cách cất giữ (Biến số)

#### give - Lịch sử thế!

```
give x=10
give y=20+5
give name="whiskers"
give health=maxHealth
```

Lệnh `give` giống như việc bạn tặng một món quà cho biến số vậy! Lịch sự hơn nhiều so với việc chỉ dùng mỗi dấu `=` lạnh lùng!
*(Góc học thuật: Trình biên dịch mèo sẽ tự động bảo vệ bạn khỏi lỗi "biến được khai báo nhưng không sử dụng" trong Go bằng cách thêm các tấm lưới bảo vệ `_ = varName`! Chúng tôi luôn giữ cho trình biên dịch Go được ăn no và vui vẻ! nya!)*

### Meo meo với loài người (Nhập/Xuất)

#### write & writeln - Meo thật to!

```
write "hello "
write "hooman"
writeln "meow meow!"
```

`write` sẽ in văn bản trên cùng một dòng (như mèo đi thẳng hàng). `writeln` sẽ xuống dòng mới sau khi in xong!

#### read - Lắng nghe sen nói

```
read playerName
```

Chờ loài sen gõ phím và lưu câu trả lời vào bát ăn (biến số).

#### read -p - Vừa hỏi vừa nghe

```
read -p "cá ngừ đâu? " playerName
```

Đưa ra câu hỏi trước, sau đó mới chờ nhập dữ liệu. Thật lịch sự làm sao! *gừ gừ*

### Đuổi đuôi và rủ bạn đi chơi! 🔁

#### gameloop - Đuổi theo cái đuôi mãi mãi!

```
gameloop start
    read -p "Ấn phím 1 để bắt đầu: " aha
    if aha == "1"
        call Game1    
    }
gameloop end
```

Lặp lại mãi mãi cho đến khi bạn bắt dừng lại! Hoàn hảo cho các trò chơi cần hoạt động liên tục!

#### Gọi bạn đi chơi! (Các khối code được đặt tên)

```
Game1 start
    dialog "Chào mừng đến với Game1!"
    give health=100
Game1 end
```

Bạn có thể tự tạo các khối code của riêng mình và dùng lệnh `call Game1` để chạy chúng! Chúng có khay cát vệ sinh riêng (phạm vi biến số độc lập), nên không lo bị lẫn lộn đồ chơi với bên ngoài đâu! KHÔNG được trộn lẫn các loại hạt với nhau!

💣 Lưu ý nhỏ từ 1 con sen: Đừng gọi khối code này bên trong một khối code khác, chúng sẽ phát nổ đấy!

#### callonce - Dành riêng cho Boss VIP!

```
callonce BossIntro
```

Gọi một khối code, nhưng **chỉ duy nhất một lần trong suốt trò chơi**! Nếu bạn cố gọi lại lần nữa, Boss mèo sẽ lờ bạn đi và tự liếm chân. Tính năng này được quản lý bằng một sơ đồ theo dõi siêu nhanh! *mew!*

### Phép thuật $ (Chèn biến vào chuỗi)

```
writeln "Xin chào, $name!"
writeln "Bạn đang có $gold đồng vàng đó."
```

Cú pháp `$ten_bien` sẽ được thay thế bằng giá trị thực tế của biến! Giống như giấu bánh thưởng vào trong đồ chơi vậy~

### Giao diện vjp pro max

#### clear - Gạt phăng mọi thứ khỏi bàn!

```
clear
```

Làm mọi thứ biến mất trong một nốt nhạc! *loảng xoảng!* Được tối ưu hóa bằng mã ANSI thuần cho màn hình Termux!

#### border - Chiếc hộp xinh xắn để chui vào

```
border top
border mid
border bot
```

Vẽ các đường viền đẹp mắt! Vì cứ chỗ nào vừa là mình chui vào nằm thôi!
- `top` và `bot`: `═════════════════════════` (Cạnh dày!)
- `mid`: `─────────────────────────` (Vạch chia mỏng!)

#### dialog, prompt, menu - Meo meo kiểu quý tộc!

```
dialog "Boss mèo phán: meo meo!"
prompt "Hãy chọn số phận của bạn: "
menu fight, defend, run away
```

- `dialog` in văn bản bên trong một chiếc hộp thoại điện ảnh sang trọng!
- `prompt` thêm một dấu mũi tên ` > ` siêu dễ thương trước văn bản của bạn!
- `menu` nhận một danh sách phân tách bằng dấu phẩy và tự động biến nó thành danh sách đánh số dạng `(1) fight, (2) defend`!

### Chỉ số và đếm hạt 📊

#### stat - Khoe chỉ số

```
stat "Máu" $playerHP
stat level $currentLevel
```

Định dạng và in các chỉ số như: `Máu: 100` hoặc `level: 5`.

#### score & level - Làm toán siêu dễ!

```
score points 10
level bossStage 5
```

Lần đầu tiên sử dụng sẽ tự tạo biến, những lần sau lệnh `score` sẽ CỘNG THÊM vào biến (`+=`), còn lệnh `level` sẽ ĐẶT THAY THẾ giá trị biến (`=`). Boss mèo thật thông minh!

#### damage & heal - Cắn một phát và đi ngủ

```
damage playerHP 10
heal playerHP 25
```

Trừ bớt hoặc cộng thêm vào một chỉ số! Đau quá! Rồi sau đó nhai nhai cỏ thảo dược để hồi phục thôi~

#### reset - Vứt vào khay cát vệ sinh

```
reset score
```

Đưa chỉ số về số 0! Khởi đầu ngày mới thanh tịnh sạch sẽ~

### random - Gạt đổ đồ vật ngẫu nhiên 🎲

```
random roll 6
```

Tạo ra một số ngẫu nhiên từ 0 đến 5 (giá trị tối đa trừ 1). Go hiện đại tự động gieo hạt ngẫu nhiên, nên không cần đánh thức Boss dậy để làm việc này đâu!

### Chôn giấu và đào đồ chơi (File dữ liệu) 💾

#### load - Đào đồ chơi lên

```
load ~/saveGame.txt
load myData ~/custom_path/file.txt
```

Đọc dữ liệu từ một file! Nếu bạn chỉ viết `load ~/my-save.txt`, Boss mèo sẽ tự động đặt tên biến là `my_save` (tự đổi dấu gạch ngang thành gạch dưới để Go không nổi giận!). Hoặc bạn có thể tự đặt tên biến rõ ràng bằng lệnh `load myData [đường_dẫn]`!

#### save - Chôn đồ chơi để dành

```
save gameData ~/saveGame.txt
```

Ghi nội dung biến số vào một file. Giấu kỹ kẻo bị mấy chú cún phát hiện!

#### peek - Ngửi tìm file

```
peek ~/saveGame.txt {
    writeln "Tìm thấy file rồi! *meo meo vui sướng*"
} else {
    writeln "Không có file nào ở đây cả, nya!"
}
```

Boss mèo sẽ ngửi xem file có tồn tại hay không trước khi bạn cố gắng mở nó!

#### create - Cào ra một file mới

```
create fish.txt
```

Cho bạn một con cá! Nhưng bên trong nó rỗng tuếch.

### Chọn bát ăn 🔀

#### if/else và else if - Đưa ra lựa chọn

Hỗ trợ các khối if/else tiêu chuẩn, giờ đây đã có thể liên kết nhiều bát ăn lại với nhau bằng `else if`!

```
if playerHP <= 0
  writeln "Trò chơi kết thúc rồi, nya..."
  susu
} else if playerHP < 20
  writeln "Nguy hiểm! Mau ăn cỏ mèo đi thôi!"
} else
  writeln "Tiếp tục chiến đấu nào!"
}
```

#### Bánh thưởng logic (and / or) 🍬

Bạn có thể kiểm tra nhiều điều kiện cùng lúc trong một câu lệnh bằng bánh thưởng logic `and` và `or`!

```
if playerHP > 50 and gold > 100
  writeln "Bạn là một chú mèo giàu có và khỏe mạnh!"
}

if choice == "1" or choice == "y"
  writeln "Vồ về phía trước!"
}
```

*(Góc học thuật: Trình biên dịch mèo tự động dịch `and` thành `&&` và `or` thành `||`, đồng thời giữ an toàn cho văn bản nằm trong dấu ngoặc kép để không vô tình đổi tên một bé mèo tên là "brandy" hoặc "orlando"!)*

#### susu - Đi ngủ trưa!

```
susu
```
Chấm dứt mọi thứ! Đến giờ đi ngủ một giấc thật ngon rồi~ 😴 (Thoát chương trình!)

#### run - Phép thuật Termux

```
run clear && ls
```

Chạy các lệnh shell tiêu chuẩn! Được liên kết chặt chẽ với `sh` để hoạt động hoàn hảo trên môi trường Termux! Không cần dùng dấu ngoặc kép, Boss mèo thông thái đã tự động xử lý các ký tự đặc biệt cho bạn rồi! 😹

#### web request (Yêu cầu mạng)

```
meow request [options] <url>

Tùy chọn:
  hiss <METHOD>        → Phương thức HTTP (GET, POST, PUT, DELETE,...)
  lick <HEADER>        → Thêm tiêu đề "Key: Value"
  spit <DATA>          → Thân yêu cầu (JSON, form data,...)
  grab <FILE>          → Tải phản hồi về file
  sniff <HEADER>       → Chỉ kiểm tra tiêu đề phản hồi (giống -i)
  yowl                 → Xem chi tiết (hiển thị tiêu đề, nội dung,...)
  purr                 → Chế độ im lặng (chỉ hiển thị phản hồi)
  scratch <N>          → Thử lại N lần nếu thất bại
  tail                 → Đi theo chuyển hướng (301, 302,...)
```

Giống như lệnh curl nhưng phiên bản mèo.

---

## Hành trình của mèo con! (Trò chơi mẫu) 🎯

Dưới đây là một trò chơi phiêu lưu nhỏ để minh họa cách mọi thứ phối hợp hoạt động:

```
clear
border top
writeln "🐱 HÀNH TRÌNH CỦA MÈO CON 🐱"
border bot

read -p "Hãy đặt tên cho bé mèo: " playerName
give playerHP=100
score gold 0

gameloop start
  clear
  border top
  stat "Tên" $playerName
  stat "Máu" $playerHP  
  stat "Vàng" $gold
  border bot
  writeln ""
  
  dialog "Bạn bắt gặp một chú chuột nhắt hung dữ!"
  
  menu Chiến đấu, Chạy trốn
  prompt "Lựa chọn của bạn: "
  read choice
  
  if choice == "1" 
    random dmg 15
    damage playerHP $dmg
    writeln ""
    dialog "Chuột nhắt tấn công! Bạn mất $dmg máu! *khè khè!*"
    
    if playerHP <= 0 
      writeln ""
      writeln "💀 TRÒ CHƠI KẾT THÚC 💀"
      susu
    }
    
    score gold 10
    writeln "Bạn đã thắng! Nhận được 10 vàng~ *gừ gừ*"
  } else if choice == "2"
    writeln ""
    dialog "Bạn đã chạy trốn an toàn! Đồ nhát gan..."
  } else
    writeln ""
    dialog "Lựa chọn không hợp lệ! Nút bấm đó có mùi lạ lắm..."
  }
  
  writeln ""
  prompt "Nhấn enter để tiếp tục hành trình..."
  read dummy
gameloop end
```

"Sao nó chạy được? tôi cũng chẳng biết." — 1 con sen đã nói vậy

---

## Công nghệ mèo sâu xa 🤓

### Cách thức hoạt động bên trong (Tế bào não hoạt động)

1. Kiểm tra file `.nk` của bạn.
2. Dịch nó sang mã Go bằng cách so khớp mẫu cực kỳ thông minh.
3. Chạy file bằng `go run` hoặc biên dịch bằng `go build`.
4. Dọn dẹp các file tạm sạch sẽ như cách mèo liếm sạch chân sau khi ăn xong~

### Các tính năng nổi bật của bộ dịch:

- **Tối ưu cho Termux**: Đường dẫn được tối ưu hóa (tự động mở rộng `~/` thành `$HOME` của Termux), thực thi bằng `sh` chuẩn và xóa màn hình bằng ANSI. Boss cực kỳ thích nằm trong chiếc hộp terminal này!
- **Tuân thủ Go nghiêm ngặt**: Loại bỏ hoàn toàn các lỗi phiền phức liên quan đến "biến khai báo nhưng không sử dụng" bằng cách tự sinh mã gán ẩn `_ = varName`!
- **Bát ăn riêng biệt (Phạm vi hàm)**: Các khối code đặt tên hoạt động hoàn toàn độc lập như các hàm Go riêng biệt. Biến toàn cục và biến cục bộ không bị lẫn lộn! Đảm bảo bát ăn của mèo nào mèo nấy ăn!
- **Theo dõi Callonce**: Sử dụng cấu trúc dữ liệu `map[string]bool` siêu tốc để ghi nhớ những khối code VIP nào đã được gọi rồi.

## Bản quyền & Đóng góp ✨

Được hoàn thiện bằng tình yêu thương từ loài mèo, dành riêng cho loài mèo (và những sen yêu mèo!)

*meo meo, chúc sen viết code vui vẻ!* 🐾
