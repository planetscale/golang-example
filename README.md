# PlanetScale Golang Quickstart

## Setup

1. Rename the `.env.example` file to `.env`
2. Add your PlanetScale connection string to the DSN environment variable
   property. It should look like
   `DSN=a07l1lwsfgix:pscale_pw_rAkT8COqDTkQTKCCkz4dIPmvVOQfxYkbPuG3F3iEcQQ@tcp(ne6uu6jolloj.us-west-2.psdb.cloud)/go-quickstart?tls=true`
3. Install the required dependencies:

```
go get github.com/gorilla/mux
github.com/joho/godotenv
gorm.io/gorm
gorm.io/driver/mysql
```

We'll use `gorill/mux` as our router for the app. We'll use `joho/godotenv` to
easily load our environment variables We'll use the `gorm` ORM and accompanying
MySQL driver to talk to our PlanetScale database.

4. Run the application via `go run main.go` in your terminal.
5. Navigate to `localhost:8080/seed` to run the migrations and add starter data.

At this point you can navigate to the other routes to see all the products
(`/products`), a single product (`/products/{id}`), all the categories
(`/categories`), or just a single category (`/categories/{id}`).

You're all set!
