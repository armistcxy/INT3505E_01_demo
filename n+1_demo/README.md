### Vấn đề mẫu cần giải quyết
Giả sử đang làm một project clone Facebook chẳng hạn, có 2 bảng là `users` và `posts` với quan hệ là 1 user có thể có nhiều post (1 - N)

**Cách làm 1**: List toàn bộ n user rồi lấy post của từng user (mỗi lần lấy post của một user là một query) 
-> Cách làm này có vấn đề khi tốn 1 query để lấy users + n query để lấy post của từng user 
-> đây chính là vấn đề n + 1 query, điều này vừa phí tài nguyên mạng dành cho việc truy vấn dữ liệu, vừa tốn thời gian khi phải chờ tận N query nữa mới có đủ dữ liệu các post

Endpoint cho cách làm 1 được đặt tại endpoint `/api/v1/users/posts`

Cách giải quyết cho vấn đề n + 1 query này được trình bày thông qua cách làm 2 

**Cách làm 2**: Sử dụng kĩ thuật `Eager Loading`, đơn giản là load sẵn các dữ liệu liên quan rồi chỉ dùng thêm 1 truy vấn nữa 

Endpoint cho cách làm 2 được đặt tại endpoint `/api/v2/users/posts`

 