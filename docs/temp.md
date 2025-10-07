Отлично! Вот корректная схема, где **используются все технологии** и роли не путаются.

# Каркас

* **API Gateway** (edge): TLS, JWT/OIDC-проверка, rate limit, маршрутизация REST для клиентов; внутрь — REST/gRPC.
* **Внутри**: синхронные запросы по **gRPC**, события — через **Kafka**, фоновые задачи — через **RabbitMQ**.

# Сервисы и хранилища

| Сервис                              | Протокол с gateway            | Хранилище                                                      | Зачем                                                           |
| ----------------------------------- | ----------------------------- | -------------------------------------------------------------- | --------------------------------------------------------------- |
| **Auth Service**                    | REST (вход), gRPC (внутрь)    | **Postgres** (пользователи), **Redis** (сессии/OTP/блокировки) | Логин/регистрация, выдача токенов, лимиты и сессии в Redis      |
| **User Profile Service**            | gRPC                          | **Postgres**                                                   | Профили клиента, контакты                                       |
| **Accounts & Cards Service** (ядро) | gRPC                          | **Postgres**                                                   | Счета, карты, баланс, лимиты; источник истины                   |
| **Payments/Transfers Service**      | gRPC                          | **Postgres** + outbox                                          | Приём платежей (sync), запись в БД и публикация событий в Kafka |
| **User Docs Service**               | REST/gRPC                     | **MongoDB**                                                    | Паспорт/КYC-сканы, договоры, выписки-PDF, любые бинарники/JSON  |
| **Notifications Service**           | gRPC (команды) → **RabbitMQ** | —                                                              | Кладёт задачи “email/SMS/push” в очереди, воркеры потребляют    |
| **Reporting/Audit/Fraud**           | — (консъюмеры)                | читает **Kafka**, пишет в **Postgres**/**Mongo** по нужде      | Аналитика, аудит, антифрод, репликации/витрины                  |

# Как используются Kafka и RabbitMQ (оба)

* **Kafka = события домена** (журнал, можно переигрывать):
  Топики: `accounts.events`, `payments.events`, `cards.events`, `auth.events`.
  Публикация через **outbox-паттерн** из сервисов на Postgres (надёжно, без потерь).
* **RabbitMQ = рабочие очереди** (таски с ack/ретраями):
  Очереди: `notify.email`, `notify.sms`, `docs.ocr`, `reports.generate`.
  Команды приходят из сервисов/гейтвея → RMQ → воркеры выполняют.

# Потоки (коротко)

**1) Вход в приложение**
Client → **Gateway** (REST) → Auth (REST/gRPC).
Auth читает пользователей из **Postgres**, хранит OTP/сессии в **Redis**. Успех → `auth.events` в **Kafka**.

**2) Перевод денег**
Client → Gateway (REST) → Transfers (gRPC).
Transfers пишет в **Postgres**, в outbox → публикует `payment.initiated` в **Kafka**.
Fraud/Reporting читают Kafka. Notifications кладёт “письмо клиенту” в **RabbitMQ** → воркер отправляет.

**3) Загрузка документов/KYC**
Client → Gateway (REST) → User Docs (REST/gRPC) → файлы и метаданные в **Mongo**.
Готовность выписки → задача `reports.generate` в **RabbitMQ** → воркер генерит PDF в Mongo и шлёт событие в **Kafka**.

# Где Redis ещё полезен

* Rate limit/блокировки в gateway.
* **Idempotency-Key** для POST-команд.
* Кэш справочников (BIN/IIN, валюты) и короткие сессии BFF.

# Минимальные правила

* **Gateway** не потребляет очереди/топики (только синхронный периметр).
* Внутренние вызовы — **gRPC**; внешнее публичное API — **REST**.
* Все изменения в Postgres → **outbox → Kafka** (надёжные события).
* Долгие/повторяемые задачи — **RabbitMQ**.
