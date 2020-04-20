
# Golang Restapi Dynamodb Serverless Auth0 App

This project includes a sample application that uses Amazon DynamoDB to perform CRUD operations with authentication/ authorization via Auth0

## Objects and Methods

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

## Serverless

Serverless framework helps you develop and deploy your AWS Lambda functions using AWS CloudFormation templates through Infrastructure as Code. The serverless.yaml file is responsible for setting up endpoint routes via Amazon API Gateway, saving persisted data via Amazon DynamoDB tables, and storing secrets via AWS Parameter Store.

## Auth0

## Amazon Simple Email Service (SES)

Amazon Simple Email Service is a cloud-based email sending service designed to send notification emails. Anytime the Artists table is updated, an email will be sent to the client (in this case jaskiratvig@gmail.com). The recipient of these emails will be defined as an AWS Parameter Store secret in serverless.yml. For this service to work, please ensure that the email that sends/receives emails is verified in the AWS console.

## DynamoDB

DynamoDB is a key-value/document No-SQL database that provides single-digit millisecond performance at any scale. Two tables are defined in serverless.yml:
* Artists: Contains the primary key of "ArtistID" and stores an artist's name, a list of songs, the subcategory of their music, and whether the artist is domestic to the United States
* SessionData: Contains the primary key of "ClientID" and stores the session state and information about the logged-in user

## AWS Parameter Store

AWS System Manager Parameter Store provides secure, hierarchical storage for configuration data/secrets management. The values for Domain, ClientID and ClientSecret can be found in the Auth0 Client settings. The values for RedirectURL and LoggedInURL are the endpoints retrieved after ``` sls deploy ``` is called. The following environment variables are defined under serverless.yml:
* Domain: The subdomain of Auth0 used to authenticate the user
* ClientID: The ID of the client
* ClientSecret: A secret of the client
* RedirectURL: The URL Auth0 redirects the user to after they have authenticated
* LoggedInURL: The URL that represents the loggedIn state of the application
* Recipient: The email address used to send all email alerts when the Artists database is updated

## Dependancies

The following dependancies are required to run the project:
* AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
* Amazon DynamoDB Table: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/SampleData.CreateTables.html
* Serverless via ``` npm install -g serverless ```
* Relevant *github.com/aws/* packages via ``` go get {PACKAGE_NAME} ```

## Run Project

To run this project, first install all the dependancies and make sure that the AWS credentials (Access Key, Secret Access Key) have been setup via ``` aws configure ``` as ``` createDynamoDBClient() ``` uses these credentials to connect to the database. Also make sure serverless is installed and then run

``` make && sls deploy ```

``` make ``` executes the MakeFile which converts the handlers into binaries. <br />
``` sls deploy ``` uses the serverless.yaml file to deploy the infrastructure to AWS CloudFormation and should return some endpoints. <br />
To send a request to an endpoint, run
``` curl -d '{"Field": "Value"}' -X {CRUD OPERATION} https://URL/ENDPOINT ```

To run the Auth0 portion of this project, navigate to ``` https://URL/home ```, click on "login" where the user will be redirected to the login endpoint. The user will be required to authenticate using either an Auth0 account or a federated identity provider (Google/Facebook). Then the user will be redirected to the loggedIn endpoint where their name will be displayed on the page.