# Tài liệu Hướng dẫn NekoKit 🐱

*Chào mừng bạn đến với tài liệu chính thức của NekoKit – bộ chuyển đổi mã nguồn đơn giản và hiệu quả dành cho Termux.*

## NekoKit là gì?

NekoKit là một trình chuyển đổi mã nguồn (transpiler) gọn nhẹ, giúp biên dịch các tệp kịch bản game viết bằng cú pháp đơn giản (tệp dạng `.nk`) thành mã nguồn Go chuẩn hóa. 

NekoKit cho phép bạn viết mã nguồn bằng các cú pháp rút gọn, dễ hiểu và tự động dịch chúng sang mã nguồn Go để chạy trực tiếp trên môi trường **Termux** hoặc máy tính cá nhân. Đây là công cụ lý tưởng để phát triển các trò chơi giao diện dòng lệnh (text-based game) mà không cần phải xử lý các cú pháp phức tạp của Go.

## Cài đặt

```bash
curl -L https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh
```

Hoặc nếu bạn muốn tự biên dịch từ mã nguồn gốc:

```bash
apt update
apt install git golang
git clone https://github.com/NekoKatoriChan/NekoKit.git
cd NekoKit
go build -o nekokit
```

## Hướng dẫn sử dụng

### Chạy trực tiếp tệp kịch bản

```bash
nekokit myGame.nk
```

Lệnh này sẽ tiến hành biên dịch và chạy tệp kịch bản của bạn ngay lập tức.

### Biên dịch thành chương trình thực thi (Binary)

```bash
nekokit myGame.nk --build
```

Lệnh này tạo ra một tệp thực thi độc lập, cho phép trò chơi chạy trên các thiết bị khác mà không cần cài đặt sẵn NekoKit.

### Biên dịch với tên tệp đầu ra tùy chỉnh

```bash
nekokit myGame.nk --build --output superCoolGame
```

### Cập nhật NekoKit

```bash
nekokit --update
```

Tự động tải về và cập nhật phiên bản NekoKit mới nhất từ kho lưu trữ trực tuyến.

### Chế độ hiển thị chi tiết (Verbose)

```bash
nekokit file.nk -v
```

Hiển thị chi tiết quá trình xử lý và chuyển đổi cú pháp ở hậu trường.

---

## Danh mục các câu lệnh 🐾

### Quản lý biến số

#### give - Khai báo và gán giá trị cho biến

```
give x=10
give y=20+5
give name="whiskers"
give health=maxHealth
```

Lệnh `give` được dùng để khai báo hoặc gán giá trị cho một biến số.
*(Lưu ý kỹ thuật: Bộ biên dịch sẽ tự động thêm cơ chế tránh lỗi "biến được khai báo nhưng không sử dụng" bằng cách chèn dòng gán ẩn `_ = varName` trong mã nguồn Go đầu ra, giúp quá trình biên dịch không bị gián đoạn).*

### Nhập và xuất dữ liệu

#### write & writeln - Hiển thị văn bản ra màn hình

```
write "hello "
write "hooman"
writeln "meow meow!"
```

- `write`: In văn bản ra màn hình nhưng không xuống dòng.
- `writeln`: In văn bản ra màn hình và tự động xuống dòng sau khi kết thúc.

#### read - Nhập dữ liệu từ người dùng

```
read playerName
```

Tạm dừng chương trình để chờ người dùng nhập văn bản từ bàn phím và lưu vào biến số được chỉ định.

#### read -p - Nhập dữ liệu kèm theo gợi ý

```
read -p "Tên của bạn là gì? " playerName
```

Hiển thị một thông điệp gợi ý trước, sau đó chờ người dùng nhập dữ liệu từ bàn phím.

### Vòng lặp và Gọi khối lệnh 🔁

#### gameloop - Vòng lặp chính của trò chơi

```
gameloop start
    read -p "Nhấn 1 để bắt đầu: " aha
    if aha == "1"
        call Game1    
    }
gameloop end
```

Tạo một vòng lặp vô hạn, thường được sử dụng làm luồng xử lý chính trong trò chơi cho đến khi có lệnh thoát.

#### Gọi khối lệnh (Khối lệnh có tên danh định)

```
Game1 start
    dialog "Chào mừng đến với Game1!"
    give health=100
Game1 end
```

Bạn có thể gom nhóm các đoạn mã thành các khối lệnh riêng biệt và dùng lệnh `call Game1` để khởi chạy chúng. Các khối lệnh này hoạt động như các hàm độc lập trong Go với phạm vi biến cục bộ riêng biệt, tránh xung đột dữ liệu với luồng chính.

*Lưu ý: Không khai báo lồng một khối lệnh bên trong một khối lệnh khác để tránh lỗi biên dịch.*

#### callonce - Chỉ gọi một lần duy nhất

```
callonce BossIntro
```

Chạy một khối lệnh được chỉ định, nhưng **chỉ thực thi tối đa một lần duy nhất trong suốt phiên chạy trò chơi**. Các lần gọi tiếp theo sẽ tự động bị bỏ qua thông qua bản đồ theo dõi trạng thái tích hợp.

### Chèn biến số vào chuỗi (String Interpolation)

```
writeln "Xin chào, $name!"
writeln "Bạn đang có $gold đồng vàng."
```

Sử dụng cú pháp `$ten_bien` bên trong chuỗi ký tự để hiển thị trực tiếp giá trị của biến số đó ra màn hình.

### Thiết kế giao diện (UI) 🎮

#### clear - Xóa sạch màn hình

```
clear
```

Xóa toàn bộ nội dung đang hiển thị trên thiết bị đầu cuối bằng chuỗi điều khiển ANSI thích hợp.

#### border - Vẽ đường viền ngăn cách

```
border top
border mid
border bot
```

Vẽ các đường viền định dạng ASCII trên màn hình:
- `top` và `bot`: In đường viền kép dày `═════════════════════════`.
- `mid`: In đường viền đơn mỏng `─────────────────────────`.

#### dialog, prompt, menu - Các hộp thoại chuyên dụng

```
dialog "Nội dung hộp thoại hiển thị ở đây"
prompt "Nhập lựa chọn của bạn: "
menu fight, defend, run away
```

- `dialog`: In văn bản nằm trong một khung hộp thoại được định dạng sẵn.
- `prompt`: Thêm ký tự trỏ nhập liệu ` > ` phía trước văn bản gợi ý.
- `menu`: Nhận vào các tùy chọn ngăn cách bởi dấu phẩy và tự động hiển thị chúng dưới dạng danh sách được đánh số, ví dụ: `(1) fight, (2) defend`.

### Theo dõi và Tính toán chỉ số 📊

#### stat - Hiển thị thông số định dạng

```
stat "Chỉ số máu" $playerHP
stat level $currentLevel
```

Định dạng và hiển thị thông số theo cấu trúc chuẩn, ví dụ: `Chỉ số máu: 100`.

#### score & level - Thiết lập và tích lũy giá trị

```
score points 10
level bossStage 5
```

- `score`: Cộng dồn giá trị được chỉ định vào biến số (`+=`). Nếu biến chưa tồn tại, hệ thống sẽ tự động khởi tạo.
- `level`: Thiết lập hoặc cập nhật trực tiếp giá trị của biến (`=`).

#### damage & heal - Trừ và cộng chỉ số

```
damage playerHP 10
heal playerHP 25
```

Thực hiện nhanh các phép toán trừ hoặc cộng một lượng giá trị vào biến số tương ứng.

#### reset - Đặt lại giá trị về mặc định

```
reset score
```

Đặt lại giá trị của biến số được chỉ định về mức `0`.

### Random - Tạo số ngẫu nhiên 🎲

```
random roll 6
```

Sinh ra một số nguyên ngẫu nhiên trong khoảng từ `0` đến `giá trị tối đa trừ 1` (ví dụ: `random roll 6` sinh số từ `0` đến `5`). Trình biên dịch sử dụng cơ chế tạo số ngẫu nhiên tự động của Go hiện đại mà không cần gieo hạt (seed) thủ công.

### Quản lý tệp tin (Đọc & Ghi dữ liệu) 💾

#### load - Đọc dữ liệu từ tệp tin

```
load ~/saveGame.txt
load myData ~/custom_path/file.txt
```

Đọc nội dung của tệp tin mục tiêu và lưu vào biến. Nếu chỉ dùng cấu trúc `load ~/path/to/file.txt`, tên biến sẽ tự động được đặt theo tên của tệp tin (các ký tự không hợp lệ như dấu gạch ngang sẽ được chuyển thành dấu gạch dưới). Bạn cũng có thể đặt tên biến cụ thể bằng cú pháp `load <ten_bien> <duong_dan>`.

#### save - Ghi dữ liệu ra tệp tin

```
save gameData ~/saveGame.txt
```

Ghi nội dung của biến số được chỉ định vào tệp tin đích để lưu trữ dữ liệu lâu dài.

#### peek - Kiểm tra sự tồn tại của tệp tin

```
peek ~/saveGame.txt {
    writeln "Tìm thấy tệp lưu trữ."
} else {
    writeln "Không tìm thấy tệp lưu trữ, khởi tạo dữ liệu mới."
}
```

Kiểm tra xem một tệp tin có tồn tại trong hệ thống hay không trước khi thực hiện các thao tác xử lý dữ liệu tiếp theo.

#### create - Tạo một tệp tin trống

```
create fish.txt
```

Tạo một tệp tin trống tại đường dẫn được chỉ định.

### Các cấu trúc điều khiển rẽ nhánh 🔀

#### if/else và else if - Điều kiện rẽ nhánh nhiều trường hợp

NekoKit hỗ trợ cấu trúc điều kiện rẽ nhánh tiêu chuẩn và các khối điều kiện mở rộng liên tiếp bằng `else if`:

```
if playerHP <= 0
  writeln "Trò chơi kết thúc."
  susu
} else if playerHP < 20
  writeln "Cảnh báo: Lượng máu của bạn đang ở mức nguy hiểm!"
} else
  writeln "Tiếp tục chiến đấu!"
}
```

#### Toán tử logic (and / or) 🍬

Hỗ trợ kiểm tra đồng thời nhiều biểu thức điều kiện trong cùng một câu lệnh bằng toán tử logic `and` (và) hoặc `or` (hoặc):

```
if playerHP > 50 and gold > 100
  writeln "Chỉ số nhân vật đang ở trạng thái tối ưu!"
}

if choice == "1" or choice == "y"
  writeln "Tiến hành di chuyển..."
}
```

*(Lưu ý kỹ thuật: Trình biên dịch tự động chuyển đổi `and` thành `&&` và `or` thành `||` trong mã nguồn Go, đồng thời bỏ qua các cụm từ trùng khớp nằm bên trong chuỗi ký tự được bọc bởi dấu ngoặc kép để tránh lỗi thay thế sai dữ liệu).*

#### susu - Thoát chương trình

```
susu
```
Dừng thực thi và thoát chương trình ngay lập tức (tương ứng với lệnh kết thúc tiến trình).

#### run - Thực thi câu lệnh hệ thống từ Termux

```
run clear && ls
```

Thực thi trực tiếp các câu lệnh shell thông qua trình bao `sh` tiêu chuẩn của Termux mà không cần bọc trong dấu ngoặc kép.

#### meow request - Thực hiện yêu cầu HTTP

```
meow request [options] <url>

Các tùy chọn:
  hiss <METHOD>        → Phương thức HTTP (GET, POST, PUT, DELETE,...)
  lick <HEADER>        → Thêm tiêu đề HTTP dạng "Key: Value"
  spit <DATA>          → Nội dung yêu cầu (JSON, form data,...)
  grab <FILE>          → Tải nội dung phản hồi về tệp tin chỉ định
  sniff <HEADER>       → Chỉ kiểm tra tiêu đề phản hồi (tương tự tùy chọn -i của curl)
  yowl                 → Chế độ hiển thị chi tiết thông tin phản hồi
  purr                 → Chế độ im lặng (không in thông tin ngoại trừ phản hồi thực tế)
  scratch <N>          → Số lần thử lại tối đa khi yêu cầu thất bại
  tail                 → Tự động chuyển hướng (Follow redirects)
```

Công cụ thực hiện yêu cầu HTTP tích hợp bên trong trình chuyển đổi.

---

## Trò chơi mẫu: Kitty Quest 🎯

Dưới đây là một kịch bản trò chơi phiêu lưu ngắn minh họa cách kết hợp các cú pháp của NekoKit lại với nhau:

```
clear
border top
writeln "🐱 KITTY QUEST 🐱"
border bot

read -p "Nhập tên nhân vật của bạn: " playerName
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
  
  dialog "Bạn bắt gặp một sinh vật thù địch trên đường di chuyển!"
  
  menu Tấn công, Bỏ chạy
  prompt "Lựa chọn của bạn: "
  read choice
  
  if choice == "1" 
    random dmg 15
    damage playerHP $dmg
    writeln ""
    dialog "Bạn bị tấn công và chịu $dmg sát thương!"
    
    if playerHP <= 0 
      writeln ""
      writeln "💀 GAME OVER 💀"
      susu
    }
    
    score gold 10
    writeln "Bạn đã chiến thắng và nhận được 10 vàng!"
  } else if choice == "2"
    writeln ""
    dialog "Bạn đã rút lui an toàn về khu vực chuẩn bị."
  } else
    writeln ""
    dialog "Lựa chọn không hợp lệ. Vui lòng thao tác lại."
  }
  
  writeln ""
  prompt "Nhấn Enter để tiếp tục hành trình..."
  read dummy
gameloop end
```

---

## Chi tiết kỹ thuật hệ thống 🤓

### Quy trình hoạt động nội bộ

1. Phân tích cú pháp tệp nguồn đầu vào dạng `.nk`.
2. Chuyển đổi mã nguồn sang định dạng ngôn ngữ Go bằng phương thức so khớp mẫu chuỗi.
3. Chạy trực tiếp chương trình bằng lệnh `go run` hoặc biên dịch thành tệp thực thi bằng `go build`.
4. Tự động dọn dẹp các tệp tin trung gian và thư mục tạm sau khi hoàn tất.

### Các đặc tính nổi bật của bộ dịch:

- **Tương thích cao với Termux**: Tối ưu hóa việc phân giải đường dẫn cục bộ (ví dụ: tự động giải nghĩa đường dẫn chứa ký tự rút gọn `~/` thành thư mục `$HOME` thực tế trên Termux), hỗ trợ ANSI để điều khiển giao diện dòng lệnh.
- **Tuân thủ tiêu chuẩn an toàn của Go**: Tránh hoàn toàn các lỗi biên dịch nghiêm ngặt của Go liên quan đến các biến được khai báo nhưng không sử dụng bằng cơ chế gán biến ẩn an toàn tự động.
- **Cách ly phạm vi biến số (Scope)**: Các khối lệnh độc lập được thiết kế tương đương với các hàm Go riêng biệt, giúp dữ liệu nội bộ của khối lệnh được bảo mật và không bị ảnh hưởng bởi biến số từ các luồng chạy khác.
- **Hệ thống theo dõi một lần (Callonce)**: Quản lý tối ưu tiến trình chạy của các phân đoạn mã đặc biệt thông qua cấu trúc dữ liệu bản đồ tra cứu nhanh dạng `map[string]bool`.

## Bản quyền & Đóng góp ✨

Phát triển bởi cộng đồng nguồn mở dành riêng cho các lập trình viên yêu thích sự tối giản và tiện ích trên nền tảng di động Termux.
