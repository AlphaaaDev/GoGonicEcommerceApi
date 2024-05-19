# GoGonicEcommerceApi

E-Commerce API app implementation. Written in Golang using gin-gonic web framework and GORM.
## Dependencies

#### [GoLang](https://go.dev/)

#### [Docker Desktop](https://www.docker.com/products/docker-desktop/) (optional)
## Environment Variables

To run this project, you will need to edit the following environment variables in your .env file. You can start by editing .env.example and renaming it to .env

`DB_USER=YOUR_DB_USER`

`DB_PASSWORD=YOUR_DB_PASSWORD`

`DB_NAME=YOUR_DB_NAME`

`DB_HOST=DB_ADDRESS`
## Installation
#### Download or clone the project
```bash
git clone https://github.com/AlphaaaDev/GoGonicEcommerceApi
cd GoGonicEcommerceApi
```

#### Deploy Database (If hosted locally)

```bash
docker-compose up
```
#### Install Go project dependencies
```bash
go get github.com/AlphaaaDev/GoGonicEcommerceApi
```
#### Create and seed Database (Only the first time you run the API)
```bash
go run main.go create seed
```
#### Run the API
```bash
go run main.go
```
## API Reference

**!!!PLACEHOLDER!!!**

#### Get all items

```http
  GET /api/items
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `api_key` | `string` | **Required**. Your API key |

**!!!PLACEHOLDER!!!**
## Features

- Authentication / Authorization
- JWT middleware for authentication
- Database seed
- ~~Paging with Limit and Offset using GORM (Golang ORM framework)~~
- CRUD operations on products, categories, orders, fetching products page
- Orders, guest users may place an order
## Project structure:
- **models**: Mvc, it is our domain data.
- **dtos**: it contains our serializers, they will create the response to be sent as json. They also take care of validating the input
- **controllers**: the mvC, they receive the request from the user, they ask the services to perform an action for them on the database.
- **seeds**: contains the file that seeds the database.
- **static**: a folder that will be generated when you create a product or tag or category with images
- **services**: contains some business logic for each model, and for authorization
- **middlewares**: contains middlewares (golang functions) that are triggered before the controller action, for example, a middleware which reads the request looking for the Jwt token and trying to authenticate the user before forwarding the request to the corresponding controller action
## Gruppo

- Avignano Luca
- Ciceri Matteo
- Cocco Alessandro
- Mazzola Christian
