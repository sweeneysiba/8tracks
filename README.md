# 8tracks
# can import the postman collection
    `https://www.getpostman.com/collections/16109a0010a8aac0c678`

# can run the project using command  go run main.go in 8tracks directory

# steps for runnning the project in local 
    make user you have mysql installed in the local system
    login to the mysql terminal and create a database with name 8tracks or else can configure the database config in file found in 
    `config/database.go`

    `dbConfig := DBConfig{
		Host:     "localhost",
		Port:     3306,          // port where the database is connected
		User:     "root",        // your UserName
		Password: "root",        // your password if any
		DBName:   "8tracks",     // database name  where you want tables to be created
	}`
  replace the {{URL}} with `http://127.0.0.1:8080`

    Method        Routes for API                                     Description Of API

    POST      {{URL}}/songs/create          -->   Create a songs | only created songs can be added to the list 
    GET       {{URL}}/songs                 -->   Fetch the Songs List 
    DELETE    {{URL}}/songs/delete/:id      -->   Delete a songs | pass id fetched from  `/song` API  instead of :id to delete that song
    GET       {{URL}}/songs/get/:id         -->   Get a songs details | pass id fetched from  `/song` API  instead of :id to get details of that song
    PUT       {{URL}}/songs/update/:id      -->   Update a songs details | pass id fetched from  `/song` API  instead of :id to update details of that     song
    
    POST      {{URL}}/list/create           -->   Create a List | only created List can be added to the list
    GET       {{URL}}/list                  -->   Fetch the List List
    DELETE    {{URL}}/list/delete/:id       -->   Delete a List | pass id fetched from `/list` API instead of :id to delete that song
    GET       {{URL}}/list/get/:id          -->   Get a List details | pass id fetched from `/list` API instead of :id to get details of that song
    PUT       {{URL}}/list/update/:id       -->   Update a List details | pass id fetched from `/list` API instead of :id to update details of that song
    
    PATCH     {{URL}}/list/updateLike/:id   -->   Update the like of list | pass id fetched from `/list` API instead of :id to update details of that song
    PATCH     {{URL}}/list/addPlay/:id      -->   Update the plays of list | pass id fetched from `/list` API instead of :id to update details of that     song
    GET       {{URL}}/explore/:tags         -->   Get all the List you want to hear by popularity of list  | instead of `:tags` pass you keyword   according 
    
                                                                                                                            to mood `tag1+tag2`  
    GET       {{URL}}/explore               -->   Get all the List you want to hear |