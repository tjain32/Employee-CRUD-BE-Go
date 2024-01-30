# Employee-CRUD-BE-Go

This is an application which is used to perform basic CRUD operations over 'Employee entity'. 

# Technologies used
 1. GoLang 
 2. PostgreSQL 
 3. Docker
 4. Docker compose

# Pre-requisites

1. Install Go from the following link - [Download and install](https://go.dev/doc/install)https://go.dev/doc/install
2. Install PostgreSQL latest version from - https://www.postgresql.org/download/
3. Install PgAdmin to visualise your databases and tables.
4. Install VSCode for the development process (Install exntensions for Golang).
5. Install Docker Desktop to containerize your application.

# Functionalities

<pre>
----------------------------------------------------------------------------------------------------------------------
"/employee"            Methods("GET")            Get all the employees from DB
"/employee/{id}        Methods("GET")            Get an employee as per Id
"/employee"            Methods("POST")           Create a new employee
"/employee/{id}"       Methods("PUT")            Update an employee
"/employee/{id}"       Methods("DELETE")         Delete an employee 
-----------------------------------------------------------------------------------------------------------------------
</pre>

# Instructions to run in local 

1. Clone the code using the URL - https://github.com/tjain32/Employee-CRUD-BE-Go.git
2. Change the DB_URL as per your connection string to PostgreSQL in file main.go
3. Run the command <code>go run .</code>
4. Your application will start running on port: 8000. 

- Author
Tanishqa Jain









