# Learn how to integrate PlanetScale with a sample Go application

This sample application demonstrates how to connect to a PlanetScale MySQL database, create and run migrations, seed the database, and display the data.

For the full tutorial, see the [Go PlanetScale documentation](https://docs.planetscale.com/tutorials/connect-go-app).

## Set up the Go app

1. Clone the starter Go application:

```bash
git clone git@github.com:planetscale/golang-example.git
```

2. Navigate into the folder:

```bash
cd golang-example
```

3. Copy the `.env.example` file into `.env`:

```bash
cp .env.example .env
``` 

## Set up the database

1. Sign up for a [free PlanetScale account](https://app.planetscale.com/sign-up).
2. Create a new database. A default branch, `main`, will be created for you.

## Connect to the Go app

1. On the database overview page in the PlanetScale dashboard, click "**Connect**".
2. Click "**New password**".
3. In the "**Connect to**" dropdown, select Go.
4. Copy the connection string.
5. Open your `.env` file and paste the connection string in as the value for `DSN`. You're now connected!

## Run migrations and seeder

1. Start the Go app:

```bash
go run .
```

2. Navigate to [`localhost:8080/seed`](http://localhost:8080/seed) to run the migrations and the seeder.
 
3. View the product and category data as follows:

- Get all products &mdash; [`localhost:8080/products`](http://localhost:8080/products)
- Get all categories &mdash; [`localhost:8080/categories`](http://localhost:8080/categories)
- Get a single product &mdash; [`localhost:8080/product/{id}`](http://localhost:8080/products/1)
- Get a single category &mdash; [`localhost:8080/category/{id}`](http://localhost:8080/categories/1)
