Simple Golang CRUD using Mux, Gorm and MySQL

Step to install:
- Clone this repository
- Create database `simple_mux`
- Run application and database migration `go run main.go`

Testing the API:
- Create products `http://localhost:3000/product`
- Display products `http://localhost:3000/products`
- Display single product `http://localhost:3000/product/{id}`
- Update product `http://localhost:3000/product/{id}`
- Delete product `http://localhost:3000/product/{id}`

