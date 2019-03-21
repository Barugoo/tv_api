**TESTING**
----
  Database scheme can be found in /mysql folder. 
  To make tests run:
  * Install Docker
  * Go into /mysql folder
  * Run following command: 
  `docker run -p 3306:3306 --name tv_api_test -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=api -d mysql`
  * Container name: tv_api_test 
  * Allocated port: 3306 
  * Wait a bit, then you can run tests
  


 
**CREATE**
----
  Creates a new record of TV

* **URL**

  /api/tv/new

* **Method:**

  `POST`
  
* **URL Params**

  None

* **Data Params**

  * **Required:** 
    **Content:** `{ id=[integer]; id > 0, brand=[string]; brand nullable, manufacturer=[string]; manufacturer str length >= 3, model=[string]; model str lenght >= 2, year=[integer]; year >= 2010 }`


* **Success Response Example:**

  * **Code:** 200 <br />
    **Content:** `{ "status" : "success", "msg" : "created record with id:3" }` 
 
* **Error Response Example:**

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ "status" : "error", "msg" : "duplicate key entry, id:3" }`



**READ**
----
  Returns specified TV record

* **URL**

  /api/tv/:id

* **Method:**

  `GET`
  
* **URL Params**

  None

* **Data Params**

  None


* **Success Response Example:**

  * **Code:** 200 <br />
    **Content:** `{ "id" : 1, "brand" : "Bravia", "manufacturer" : "Sony", "model" : "HX929", "year" : 2011 }`

 
* **Error Response Example:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ "status" : "error", "msg" : "record with id:666 is not found" }`



**UPDATE**
----
  Updates specified record of TV

* **URL**

  /api/tv/:id

* **Method:**

  `PUT`
  
* **URL Params**

  None

* **Data Params**

   * **Required:** 
    **Content:** `{ id=[integer]; id > 0, brand=[string]; brand nullable, manufacturer=[string]; manufacturer str length >= 3, model=[string]; model str lenght >= 2, year=[integer]; year >= 2010 }`

* **Success Response Example:**

  * **Code:** 200 <br />
    **Content:** `{ "status" : "success", "msg" : "updated record with id:2" }` 

 
* **Error Response Example:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ "status" : "error", "msg" : "nothing to update: record with id:666 is not found" }`



**DELETE**
----
  Deletes specified record of TV

* **URL**

  /api/tv/:id

* **Method:**

  `DELETE`
  
* **URL Params**

  None

* **Data Params**

  None

* **Success Response Example:**

  * **Code:** 200 <br />
    **Content:** `{ "status" : "success", "msg" : "deleted record with id:1" }` 

 
* **Error Response Example:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ "status" : "error", "msg" : "nothing to delete: record with id:1 is not found" }`    

