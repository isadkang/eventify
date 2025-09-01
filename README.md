# Eventify API 🎉

Eventify adalah aplikasi backend berbasis **Golang + Gin + PostgreSQL** untuk mengelola **event**, **tickets**, dan **quizzes**.  
Mendukung dua role utama:  
- **User (peserta)**: bisa register, login, join event, melihat tiket, mengerjakan quiz.  
- **Admin**: bisa mengelola user, event, tiket, dan quiz.

## 🚀 Tech Stack
- [Golang](https://go.dev/) (Gin Framework)
- [PostgreSQL](https://www.postgresql.org/)
- [Railway](https://railway.app/) (deployment)

## 📦 Setup Project

1. **Clone repo**
```bash
   git clone https://github.com/username/eventify.git
   cd eventify
```

2. **Copy `.env` file**

```bash
   cp .env.example .env
```

Sesuaikan value:

```env
   DATABASE_URL=postgres://user:password@localhost:5432/eventify?sslmode=disable
   JWT_SECRET=supersecret
```

3. **Install dependencies**

```bash
   go mod tidy
```

4. **Run migration (contoh pakai goose)**

```bash
   goose up
```

5. **Run server**

```bash
   go run main.go
```

Server akan jalan di `http://localhost:8080`.

## 🔑 Authentication Flow

* **Register** → dapat membuat akun baru.
* **Login** → akan mengembalikan **JWT Token**.
* Gunakan token tersebut di **Authorization Header**:

```
  Authorization: Bearer <token>
```

## 📖 API Endpoints

### Public

#### `POST /api/register`

**Request**

```json
{
  "name": "Irsyad",
  "email": "irsyad@mail.com",
  "password": "secret"
}
```

**Response**

```json
{
  "message": "user registered"
}
```

---

#### `POST /api/login`

**Request**

```json
{
  "email": "irsyad@mail.com",
  "password": "secret"
}
```

**Response**

```json
{
  "token": "<jwt-token>"
}
```

---

### User (Auth required)

#### `GET /api/me`

**Response**

```json
{
  "id": 1,
  "name": "Irsyad",
  "email": "irsyad@mail.com",
  "role": "peserta"
}
```


#### `POST /api/events/:id/join`

Join ke sebuah event.

**Response**

```json
{
  "message": "joined event",
  "ticket_id": 5
}
```

#### `GET /api/tickets`

**Response**

```json
[
  {
    "id": 5,
    "user_id": 1,
    "username": "Irsyad",
    "event_id": 2,
    "event_name": "Tech Conference",
    "status": "pending",
    "date": "2025-09-01T10:00:00Z"
  }
]
```

#### `GET /api/events/:id/quizzes`

**Response**

```json
[
  {
    "id": 1,
    "event_id": 2,
    "question": "Siapa penemu Golang?",
    "options": ["Rob Pike", "James Gosling", "Linus Torvalds"],
    "answer_key": "Rob Pike"
  }
]
```

#### `POST /api/events/:id/quizzes/submit`

**Request**

```json
{
  "answers": [
    {"quiz_id": 1, "answer": "Rob Pike"}
  ]
}
```

**Response**

```json
{
  "message": "quiz submitted",
  "score": 100
}
```

#### `GET /api/quizzes/submissions/me`

**Response**

```json
[
  {
    "id": 1,
    "user_id": 1,
    "username": "Irsyad",
    "event_title": "Tech Conference",
    "event_id": 2,
    "score": 100,
    "status": "passed",
    "submitted_at": "2025-09-01T10:00:00Z"
  }
]
```

### Admin (Auth: role = admin)

#### `GET /api/admin/dashboard`

Dashboard statistik (jumlah user, event, ticket, quiz, dll).
**Response (contoh)**

```json
{
  "users": 10,
  "events": 3,
  "tickets": 25,
  "quizzes": 12
}
```

#### `GET /api/admin/users`

Daftar semua user.

#### `GET /api/admin/users/:id`

Detail user berdasarkan ID.

#### `GET /api/admin/events`

List event.

#### `POST /api/admin/events`

**Request**

```json
{
  "title": "Tech Conference",
  "description": "Event teknologi terkini",
  "date": "2025-09-20T10:00:00Z",
  "location": "Jakarta",
  "quota": 100
}
```

**Response**

```json
{
  "message": "event created",
  "data": {
    "id": 2,
    "title": "Tech Conference",
    "description": "Event teknologi terkini",
    "date": "2025-09-20T10:00:00Z",
    "location": "Jakarta",
    "quota": 100
  }
}
```


#### `PUT /api/admin/tickets/:id/approve`

**Response**

```json
{
  "message": "ticket approved"
}
```

#### `PUT /api/admin/tickets/:id/reject`

**Response**

```json
{
  "message": "ticket rejected"
}
```


#### `POST /api/admin/events/:id/quizzes`

Tambah quiz ke event.

#### `GET /api/admin/events/:id/quizzes/submissions`

Lihat semua submission quiz untuk event tertentu.


## 📝 Notes

* Semua endpoint **User** perlu JWT dengan role `peserta`.
* Semua endpoint **Admin** perlu JWT dengan role `admin`.
* Model JSON response sudah disesuaikan dengan struktur di `models/`.


## 🤔 Kenapa Strukturnya Begini?

* **Separation of concerns**: User & Admin punya route group + middleware masing-masing.
* **Models terpisah**: bikin gampang maintenance & mapping DB.
* **JWT Auth**: lebih aman & scalable untuk multi-role.
* **PostgreSQL**: powerful relational DB untuk handle relasi user-event-ticket-quiz.
