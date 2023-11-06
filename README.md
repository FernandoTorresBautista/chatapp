# chat application
Webapp to register/login and use a common chat between the users.<br>
You can create your own chatroom or join to an existing one different from the common room.<br>
The messages received with the prefix "/stock=" and followed by the stock_code, it will be managed by a separate process to make the request to the api and send the response to the channel, the message it will be send to all clients in the same room.<br>
Rabbitmq save the last 50 messages.<br>

## Application layout
main.go <br>
root of the application, here we start the application <br>
- app/<br>
    - api/<br>
        This folder allow the code to start your application, usually the router of the api have in their own directory with their respective purpose or version 
    - biz/<br>
        The biz layer works like a bridge, starting the third part application, our owns pkg, clients like mysql, redis, kafka, etc..; in the start function and using their availables functions that we add in the BizHandle interface
    - client/<br>
        This folder is used to load the services like db, kafka, redis, payments gateways ans so on
    - config/<br>
        This folder is used to load the default configuration and the environment variables that we can add to the application
- build/<br>
    This folder is used to save the builds of the application
- deployment/<br>
    This folder is used to save the yaml files to deploy the application easily, usually have the yaml file for the namespace, secrets, rbca, configmap, ingress, service and deployment|pod|statefulset.
- docs/<br>
    This folder is used to save the swagger documentation<br>

The pkg folder allows to have the differents services and modules that are used for ours app/client and the app/biz layer
- pkg
    - clients/
        - mysql/<br>
            This able the application to use mysql, this is used for our app/client/db/mysql module
    - rabbitmq/
        this is the simple create queue and push messages to the queue




## Start application
You need to create a config.yaml file in app/config/ folder with your environment variables or locally 
run the application with the variables before the go run, example: <br>
`DB_TYPE='mysql' 
DB_MYSQL_IP='localhost:3306' 
DB_MYSQL_NAME='db_name' 
DB_MYSQL_USER='db_user' 
DB_MYSQL_PASS='db_pass' 
DB_MYSQL_RETRY='5' 
MIGRATE_DB_USER='db_user' 
MIGRATE_DB_PASS='db_pass' 
MIGRATE_DB=false 
CONTINUE_AFTER_MIGRATE=false go run main.go`

NOTE: is better to include this variables in the environmentwhen you are deploying something, just to not change or add more extra steps.



## Recommendations
For the db use golang-migrate, this [link](https://github.com/golang-migrate/migrate) get you the list of drivers <br>

Installation in ubuntu using WSL<br>
First add the script to migrate for golang:<br>
    curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash<br>
Then run the next commands<br>
    sudo apt update <br>
    sudo apt install migrate <br> 

Now you can use migrate. The "create", "goto", "up" and "down" commands are availables to do what you want to.<br>

The create command use the name_of_script that you want to create:<br> 
- migrate create -ext=.sql [name_of_script] <br>

The goto command use the version of the file created by the create command, the version is the prefix of the file until the first "_" character:<br>
- migrate -path=[scripts_folder] -database="mysql://[user]:[password]"@tcp([ip])/[db_name] goto [version]<br>

The up command update the db to the latest version:<br>
- migrate -path=[scripts_folder] -database="mysql://[user]:[password]"@tcp([ip])/[db_name] up <br>

The down command downgrade a specific counts of migrations steps:<br>
- migrate -path=[scripts_folder] -database="mysql://[user]:[password]"@tcp([ip])/[db_name] down [count_steps]<br>

The others parameters are the scripts_folder, this is the path to the folder where you are creating your scripts. user is the user to your database, password is the password to your database, the ip and bd_name should exist.<br>

---

Rabbitmq used to test 3.8:
Using docker, more simple and fast <br>
- docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
With the defualt user/password
- docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3.8-management

--- 

Swagger used to test the endpoints an the responses that we need<br>
install command: <br>
- go get -v -u github.com/swaggo/swag/cmd/swag

install library http-swagger in the project 
- go get "github.com/swaggo/http-swagger"

run swagger: <br>
Usually you can use just swag(if the command not found specify the path to the executable)
- $HOME/go/bin/swag init

---

