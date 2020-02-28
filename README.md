# golang-restapi-dynamodb-serverless

## 1. helloWorld

This project includes a sample application that uses Amazon DynamoDB to perform CRUD operations. A unit test file is also included to ensure functionality of the program.

### Objects and Methods

The struct being used for this project is an *Artist* with the following properties:

```
type Artist struct {
	ArtistID    string
	Name        string
	Songs       []string
	Subcategory string
	Domestic    bool
}
```

The Artist struct is used to perform the following CRUD operations:

* createArtist: 
  * Endpoint: /artists
  * Request: POST request 
  * Input: Client passes in the new object's properties via JSON body
  * Output: A new object is created in the database
* getArtist:
  * Endpoint: /artists/{artistName}
  * Request: GET request 
  * Input: Client passes in the name object to be retreived via endpoint field
  * Output: All the fields of the object are outputted to the console
* getAllArtists:
  * Endpoint: /artists
  * Request: GET request 
  * Input: None
  * Output: All the objects are outputted to the console
* deleteArtist:
  * Endpoint: /artists/{artistName}
  * Request: GET request 
  * Input: Client passes in the name of the object to be deleted via endpoint field
  * Output: Object is deleted from the database
* editArtist:
  * Endpoint: /artists/{artistName}
  * Request: PUT request 
  * Input: Client passes in the new object's properties via JSON body and the name of the object to be editted via JSON body
  * Output: All the fields of the object are outputted to the console and the values of the object are editted in the database

### Run Project

To run this project, first install all the dependancies and make sure that the AWS credentials (Access Key, Secret Access Key) have been setup via ``` aws configure ``` . Then run 

``` go run main.go ```

Open up Postman and send a CRUD request to Port 8080 passing in the respective JSON body.

## 2. golang-restapi-dynamodb-serverless

This project deploys the helloWorld application from above to Lambda functions using Serverless.

### Serverless

One addition to the helloWorld project is the serverless.yaml file, which is responsible for deploying Amazon API Gateway as Infrastructure as Code.

### Code Organization

As opposed to the helloWorld project, the only line inside the main function for each golang file is

``` lambda.Start(Handler) ```

This command executes the respective handler based on the CRUD request and endpoint.

### Run Project

To run this project, make sure serverless is installed and run

``` make && sls deploy ```

``` make ``` executes the MakeFile which converts the handlers into binaries. <br />
``` sls deploy ``` uses the serverless.yaml file to deploy the infrastructure to AWS CloudFormation and should return some endpoints. <br />
To send a request to an endpoint, run
``` curl -d '{"Field": "Value"}' -X {CRUD OPERATION} https://URL/ENDPOINT ```
