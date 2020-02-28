package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Artist struct {
	ArtistID    string
	Name        string
	Songs       []string
	Subcategory string
	Domestic    bool
}

var artists []Artist

func createDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	return dynamodb.New(sess)
}

func (a Artist) String() string {
	return fmt.Sprintf("{%s, %s, %s, %s, %t}", a.ArtistID, a.Name, a.Songs, a.Subcategory, a.Domestic)
}

//Create

func createArtist(name string, songs []string, subcategory string, domestic bool) Artist {
	u, _ := uuid.NewV4()
	artist := Artist{u.String(), name, songs, subcategory, domestic}
	artists = append(artists, artist)
	return artist
}

func createArtistHandler(w http.ResponseWriter, r *http.Request) {
	var a Artist
	err := json.NewDecoder(r.Body).Decode(&a)

	artist := createArtist(a.Name, a.Songs, a.Subcategory, a.Domestic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svc := createDynamoDBClient()

	av, err := dynamodbattribute.MarshalMap(artist)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Artists"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(w, "Added %+v", artist)
}

//Get Artist

func getArtist(name string) Artist {
	for i := 0; i < len(artists); i++ {
		if artists[i].Name == name {
			return artists[i]
		}
	}
	return Artist{}
}

func getArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	svc := createDynamoDBClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Artists"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(vars["artistName"]),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	item := Artist{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.Name == "" {
		fmt.Println("Could not find '" + item.Name)
		return
	}

	fmt.Fprintf(w, "Found artist: ")
	fmt.Fprintf(w, "ArtistID:  %+v ", item.ArtistID)
	fmt.Fprintf(w, "Name: %+v ", item.Name)
	fmt.Fprintf(w, "Subcategory:  %+v ", item.Subcategory)
	fmt.Fprintf(w, "Songs: %+v ", item.Songs)
	fmt.Fprintf(w, "Domestic: %+v ", item.Domestic)
}

//Get All

func getAllArtistsHandler(w http.ResponseWriter, r *http.Request) {
	svc := createDynamoDBClient()

	proj := expression.NamesList(expression.Name("Name"), expression.Name("Subcategory"), expression.Name("Songs"), expression.Name("Domestic"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Artists"),
	}
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	for _, i := range result.Items {
		item := Artist{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Fprintf(w, "Name: %+v ", item.Name)
		fmt.Fprintf(w, "Subcategory: %+v ", item.Subcategory)
		fmt.Fprintf(w, "Songs: %+v ", item.Songs)
		fmt.Fprintf(w, "Domestic: %+v ", item.Domestic)
	}
}

//Delete

func deleteArtist(artistName string) {
	for i := 0; i < len(artists); i++ {
		if artists[i].Name == artistName {
			artists = append(artists[:i], artists[i+1:]...)
		}
	}
}

func deleteArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	svc := createDynamoDBClient()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(vars["artistName"]),
			},
		},
		TableName: aws.String("Artists"),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(w, "Deleted %+v", vars["artistName"])
}

//Edit

func (art Artist) editArtist(artist Artist) Artist {
	deleteArtist(artist.Name)
	art.ArtistID = artist.ArtistID
	art.Name = artist.Name
	art.Songs = artist.Songs
	art.Subcategory = artist.Subcategory
	art.Domestic = artist.Domestic
	artists = append(artists, art)
	return art
}

func editArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var a Artist
	err := json.NewDecoder(r.Body).Decode(&a)
	artist := getArtist(vars["artistName"]).editArtist(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svc := createDynamoDBClient()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":so": {
				SS: aws.StringSlice(artist.Songs),
			},
			":su": {
				S: aws.String(artist.Subcategory),
			},
			":d": {
				BOOL: aws.Bool(artist.Domestic),
			},
		},
		TableName: aws.String("Artists"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(vars["artistName"]),
			},
		},
		UpdateExpression: aws.String("SET Songs = :so, Subcategory = :su, Domestic = :d"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Fprintf(w, "Updated %+v", vars["artistName"])
}

//Main

func main() {
	var songsRL []string
	songsRL = append(songsRL, "Aurora")
	songsRL = append(songsRL, "Valhalla")
	songsRL = append(songsRL, "Core")

	_ = createArtist("RLGrime", songsRL, "Trap", true)

	var songsSlander []string
	songsSlander = append(songsSlander, "All You Need To Know")
	songsSlander = append(songsSlander, "First Time")

	_ = createArtist("Slander", songsSlander, "Future Bass", true)

	r := mux.NewRouter()

	r.HandleFunc("/artists/{artistName}", editArtistHandler).Methods("PUT")
	r.HandleFunc("/artists/{artistName}", getArtistHandler).Methods("GET")
	r.HandleFunc("/artists", getAllArtistsHandler).Methods("GET")
	r.HandleFunc("/artists/{artistName}", deleteArtistHandler).Methods("DELETE")
	r.HandleFunc("/artists", createArtistHandler).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
