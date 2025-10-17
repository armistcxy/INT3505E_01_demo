### Extensibility
Extensibility thể hiện qua việc định nghĩa tài nguyên và sử dụng HTTP chuẩn trong REST API. Không cố định hành vi bằng endpoint cụ thể như /createUser hay /deleteUser mà chỉ cần định nghĩa tài nguyên như /users và dùng các HTTP method chuẩn (GET, POST, DELETE, v.v) -> Tức là khi ta muốn mở rộng hệ thống thêm nhiều loại resource mới thì chỉ cần thêm qua các endpoint mới như /payments, /reviews, một cách khác là mở rộng field JSON

Quan trọng nhất của extensibility còn là khả năng thiết kế API sao cho DỄ MỞ RỘNG và KHÔNG PHÁ VỠ CÁI CŨ


