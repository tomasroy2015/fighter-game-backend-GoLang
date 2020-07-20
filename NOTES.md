# Notes

*Database setup, structure change, tests setup, etc.*
1) restructred the project by my way as i could not understand some point
2) seeder added for automigration
3) main.go modified
4) all code modified in api folder
5) restructured code
 --GO-ASSIGNMENT   
    --api
      --controllers
      --middlewares
      --models
      --responses
      --seed
    --server.go
    --test
 --main.go

//instruction to run the project
Database configuration are set in .env file
step1) create database named zenport_db for production and zenport_db_test for testing
step2) run the project by writting  -> go run main.go
step3) run the test -> go test -v --run ./api/tests/...

Docker integration is not set up because docker kernel is not working properly in my windows PC.

