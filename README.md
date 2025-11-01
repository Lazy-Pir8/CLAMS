Club Attendance & Leave Management – Go Backend **(MVP version)**
This project is a backend API for managing student attendance and leave requests, with built-in admin controls and role-based access.
It uses Go + Gin framework, using JWT authentication and a PostgreSQL database (via GORM).


Features Overview - 



Student, Admin, Warden, and Faculty roles (configurable at registration)

JWT-based authentication, with different permissions for roles *(Not Yet Completed)*

Leave request flow: apply, approve/reject

Daily attendance marking for each student *(manual for now)*

Attendance stats endpoint (percentage, present/absent count)

Clean REST API (easy for Postman)

Setup Instructions
Clone this repo
(or just place the code in a folder)

Spin up Postgres
You can use Docker or install locally.
Default database: postgres on port 5432, user: postgres, password: mysecretpassword.

With Docker:

text
docker run --name my-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
Install dependencies

text
go mod download
Run the server

text
go run main.go
API Endpoints
All POST/Protected endpoints require JWT token in Authorization: Bearer <token> header after login.

Authentication
POST /register
Register a user (provide: name, email, password, role)

json
{
    "name": "harsh",
    "email": "harsh@mail.com",
    "password": "secret",
    "role": "student"
}
POST /login
Logs in and returns a JWT token.

Attendance
POST /attendance/mark
Mark a student as present/absent on a date

json
{
  "student_id": 1,
  "date": "2025-11-02",
  "present": true
}
GET /attendance/stats/:student_id
Returns number of present days, total days, attendance %
Sample response:

json
{
  "student_id": 1,
  "present_days": 76,
  "total_days": 90,
  "attendance_percentage": 84.4
}
Leave
POST /leaves/apply
Request leave (students)

json
{
  "student_id": 1,
  "reason": "medical",
  "start_date": "2025-11-05",
  "end_date": "2025-11-08"
}
PATCH /leaves/:id/status
Update/approve/reject leave (admin, warden)

json
{ "status": "approved" }
Admin-only APIs *(yet to be configured)*
GET /users – List all registered users

Later: add role-limited controls, filtering, etc.

- How to Test
Use Postman or curl for API testing.

Register and login as different roles.

Copy your JWT token to “Authorization: Bearer <token>” in the header for protected routes.

Known Shortcuts & Things to Improve
Passwords stored as plain text — in production, use bcrypt.

JWT tokens are not refreshable (24-hour expiry).

Date handling is currently strings; *(I will modify them with more suitable data types)*

No frontend (pure API backend).

NOTE- *This is the MVP version of my project, i will modify it with time and improve many things which I know I have messed up in this version.*
