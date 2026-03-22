# REST API (v1)

Base path: `/api/v1`

## Interactive docs
- Scalar reference UI: `/docs`
- OpenAPI schema JSON: `/docs/openapi.json`
- The schema is generated from the backend DTO and route metadata so request and response models stay aligned with the API implementation.

## Auth
- `POST /auth/signup`
- `POST /auth/login`
- `GET /auth/me`

## Students
- `POST /students`
- `GET /students`
- `GET /students/{id}`
- `PUT /students/{id}`
- `DELETE /students/{id}`
- `GET /students/{id}/enrollments`

## Teachers
- `POST /teachers`
- `GET /teachers`
- `GET /teachers/{id}`
- `PUT /teachers/{id}`
- `DELETE /teachers/{id}`

## Class/Course
- `POST /class-courses`
- `GET /class-courses`
- `GET /class-courses/{id}`
- `PUT /class-courses/{id}`
- `DELETE /class-courses/{id}`

## Optional Fees
- `POST /class-courses/{id}/optional-fees`
- `GET /class-courses/{id}/optional-fees`
- `PUT /optional-fees/{id}`
- `DELETE /optional-fees/{id}`

## Enrollments
- `POST /enrollments`
- `GET /enrollments`
- `GET /enrollments/{id}`
- `PUT /enrollments/{id}`
- `GET /students/{id}/enrollments`

## Payments
- `POST /payments`
- `GET /payments`
- `GET /payments/{id}`
- `GET /enrollments/{id}/payments`

## Expenses
- `POST /expenses`
- `GET /expenses`
- `GET /expenses/{id}`
- `PUT /expenses/{id}`

## Receipts
- `GET /receipts`
- `GET /receipts/{key}`
  - numeric `key` resolves by ID
  - string `key` resolves by `receipt_no`

## Reports
- `GET /reports/dashboard`
- `GET /reports/students`
- `GET /reports/teachers`
- `GET /reports/class-courses`
- `GET /reports/gross`
- `GET /reports/transactions`
- `GET /reports/performance`

## Settings
- `GET /settings`
- `PUT /settings`

## Auth requirement
- All endpoints except `/auth/signup`, `/auth/login`, `/auth/me` (with token), and `/healthz` require bearer authentication.
