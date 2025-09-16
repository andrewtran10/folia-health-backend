# Getting Up & Running

## BEFORE YOU GET STARTED

Our goal with this project is to assess your application development & problem solving capability, not necessarily your DevOps acumen. If you are to run into any issues while getting this starter project up and running, feel free to reach out & we will do our best to get you up and running!

## Pre-requisites

- [Docker](https://www.docker.com/get-started/) has been downloaded & installed on your machine

## Starting the application

- Ensure that Docker is running on your machine
- In your terminal, navigate to the location of the project
- Navigate one level down into the `folia-backend-project/laravel-starter` directory. This will be the base directory for the application.
- Create the SQLite database file by running the following: `touch source/database/database.sqlite`
- Run the following command to build & start the application container & NGINX server: `make up`
  - This may take a few minutes the first time you run it.
- Run `make migrate-fresh` to run our example stack of database migrations and seeds
  - This command will wipe the database, generate new database tables as they are defined in migration files, and seed test data rows in that order.
- If the build & start go as expected, you should now have the application running on your local machine at the address `http://localhost:8080`
  - Open this address in a browser to see a hello world message.
  - This project will primarily work with API routes, so `http://localhost:8080/api` will likely be the base URL for many or all of your routes.

### Makefile reference

Here are all of the commands we have defined for you in the `Makefile` in an effort to make Docker less of a hassle for this project:

| command    | Description                                                               |
|------------|---------------------------------------------------------------------------|
| make build         | Build the images defined in the compose file provided             |
| make up            | Starts the docker containers defined in the compose file provided |
| make stop          | Stops the containers defined in the compose file provided         |
| make down          | Removes the containers defined in the compose file provided       |
| make migrate-fresh | Runs `php artisan migrate:fresh --seed` on the application db     |

**For more information on makefiles, see [this link](https://www.gnu.org/software/make/manual/make.html#Introduction)
