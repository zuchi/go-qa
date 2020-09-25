# Questions and Answers

### Description
The idea of this api is receiving questions and answers from the users. However, each question there will be 
just one response. If the user think that his answer is better than other they can substitute it.

### What Do I need to run?
 - You will need to have installed the MongoDB (I use Mongo Version 4.4) without a certificate configurated.
 - You will need to Golang Version 1.15 installed, with the GoPath configured in your environment
 - You need to *make* command installed to perform some CLI operations
 - If you wouldn't like to use make command to run the API, you must configure some environment variables that I will describe bellow: </br>
 ```
 - MONGO_URL: This Variable configure the URL to connect on MongoDB. You need to use something like it, for example: mongodb://localhost:27017
 - MONGO_COLLECTION: This Variable configure what is the name of collection that you'd like to put into mongodb. Example: BairesDev
 - SERVER_URL: This Variable configure what is the Address and port that the server are going to running. Example: localhost:3000 or :3000 </br> 
```
- If you are in OS that have make cli, just edit the file, configure the environment variables there, save and in your prompt type make run after all services dependences are online.

### What I used of develop it:
I used Golang version 1.15 with Go Modules enabled;
I used MongoDB version 4.4 in a Docker container;
I used MacOs Operation System (Catalina version);
I used the Goland Jetbrains IDE;
I used Postman to consume this API.

 
### What are command that have in Makefile
Here I will describe some commands that you can run using Makefile interface
 ```
make build: this command will compile golang project
make run: this command will compile and run the project
make test: to see unit test coverage
```


#### Documentation
Here, I will describe the some points that I consider important to share.

##### Endpoint
Describing all 3 endpoint of API:

- ``GET: <url>/question/``: This will bring to us all questions and answers that we have in questions collections. The payload returned by api is:
```json
[
    {
        "id": "5f6c121031d86158915d8b83",
        "question": "What is the best language in the world?",
        "user": "Zuchi, Jederson",
        "createdDate": "0000-12-31T20:53:32-03:06",
        "updatedDate": "2020-09-24T00:29:39.319-03:00",
        "answer": "golang is the best language of the entire world."
    },
    {
        "id": "5f6c1321880b9da8ba7c6063",
        "question": "The rust is a good language",
        "user": "Zuchi, Jederson",
        "createdDate": "2020-09-24T00:31:40.929-03:00",
        "updatedDate": "2020-09-24T00:32:35.621-03:00",
        "answer": "Yes, it's good but go is better"
    }
]
```
The status code that can be return here are: 200 (StatusOK) or 500 (InternalServerError) <br/><br/>


``POST: <url>/question/``: In this Url we can perform some insert operations. The Json that need to be send is a single json object described bellow:
```json
{
    "question": "The rust is a good language",
    "user": "Zuchi, Jederson",
    "answer": "Yes, it's good but go is better"
}
```
The only attribute that is **required** in Json is the **question** attribute. The **others attributes are optional**. <br/>
The Status code that can be return here are: 201 (Created), 400 (BadRequest) or 500 (InternalServerError) <br/><br/> 


``PUT: <url>question/<id>``: In this Url we can perform some update operations. The Json that need to be sent is a single json object described bellow:
```
{
    "question": "The rust is a good language",
    "user": "Zuchi, Jederson",
    "createdDate": "2020-09-24T00:29:39.319-03:00"
    "answer": "Yes, it's good but go is better"
}
```
The only attribute that is **required** in Json is the **question** attribute. However, I recommend you that you send all Attributes describe above, the only exception is the answer attribute that we can send or not.
The Status code that can be return here are: 204 (NoContent), 400 (BadRequest) or 500 (InternalServerError) <br/><br/>

##### Database Collection
In the database collection I Have just one collection called question stored on mongoDb

##### Future works, If I will:
- In put endpoint I could make other struct to improve the attributes validation.
- In mongodb repository I will need make some improvement in Unit test.
- I'd like to execute some integration test and also some end-to-end test
- Also We can use an API-Key. However, just to demonstrate proposed, I thought that is not important in this moment.
- Last but not least, I'd like to make a front end using VueJs in version 3.0