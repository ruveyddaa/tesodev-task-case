# TESODEV Product API

This is a simple RESTful API for managing products, developed as part of the TESODEV Backend Internship Case.

---

##  Features

- Create, read, update, and delete products (CRUD)
- Retrieve product by ID
- Search by product name
- Filter and sort by price

---

## Tech Stack

- Golang
- Echo Web Framework
- MongoDB (via Docker)

---

## Running the Project

1. Start MongoDB with Docker:
```bash
docker-compose up -d
```

2. Run the server:
```bash
go run main.go
```

---

##  API Endpoints

- `GET /products`
- `GET /products/:id`
- `POST /products`
- `PUT /products/:id`
- `PATCH /products/:id`
- `DELETE /products/:id`
- `GET /products/search?name=pen&minPrice=10&maxPrice=50&sort=asc`

---

##  Author

Developed by Rüveyda Baran.
