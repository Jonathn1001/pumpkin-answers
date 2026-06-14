# Pumpkin Answers

Repo này là bài nộp cho thử thách của **Pumpkin**, gồm **2 phần** nằm trong 2 thư mục con. README này chỉ mô tả sơ qua — chi tiết của mỗi phần xem trong README riêng của thư mục đó.

> 🔗 **Live demo (engineer-challenge):** https://pumpkin-answers-production.up.railway.app — UI ở `/`, API ở `/api`, health ở `/healthz`.

| Thư mục | Phần | Nội dung |
| --- | --- | --- |
| [`engineer-challenge/`](./engineer-challenge/README.md) | Bài coding | Nền tảng cấu hình & xử lý claims bảo hiểm đa tenant |
| [`logical-questions/`](./logical-questions/README.md) | Câu hỏi tư duy | Trả lời 7 câu hỏi logic/cá nhân |

---

## 1. `engineer-challenge/` — Bài coding

Một nền tảng **đa tenant** để **cấu hình** cách xử lý claims bảo hiểm cho từng tenant, và **chạy** claims qua một decision engine điều khiển hoàn toàn bằng config. Mỗi tenant có một cấu hình được **đánh version** trải trên 6 khía cạnh (branding, loại claim, luồng duyệt, thông báo, SLA, custom fields) — chính config mà admin UI chỉnh sửa cũng là thứ engine thực thi.

**🔗 Live demo:** [pumpkin-answers-production.up.railway.app](https://pumpkin-answers-production.up.railway.app) — deploy thật trên Railway (1 service Go + Postgres), seed sẵn 3 tenant mẫu.

- **Backend:** Go 1.25 + Gin + GORM/Postgres
- **Frontend:** React 19 + Ant Design 6 + Vite
- Có sẵn 3 tenant mẫu (SafeGuard, HealthFirst, GovHealth) seed lúc khởi động.

→ Cách chạy, kiến trúc, API và demo: xem [`engineer-challenge/README.md`](./engineer-challenge/README.md).

## 2. `logical-questions/` — Câu hỏi tư duy

Phần trả lời **7 câu hỏi** logic/cá nhân (viết bằng tiếng Việt), mỗi câu được trả lời một cách chân thật và ngắn gọn — từ ký ức tuổi thơ đến cách viết hướng dẫn đổ xăng nhanh nhất.

→ Toàn bộ câu hỏi và câu trả lời: xem [`logical-questions/README.md`](./logical-questions/README.md).
