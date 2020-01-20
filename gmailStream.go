package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

var user string = "me"
var srv *gmail.Service

func GmailStream() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err = gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	CheckLetter()

	return
}

//CheckLetter проверка на новые письма
func CheckLetter() {

	fmt.Println("Check message...")

	/*r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}*/

	//MSGList := Messages{}

	res, _ := srv.Users.Messages.List(user).Q("subject:Make message is:unread").Do()

	tok := res.Messages
	count := res.ResultSizeEstimate

	for count > 0 {
		msg, _ := srv.Users.Messages.Get(user, tok[count-1].Id).Format("full").Do()

		Rawadr := []byte(msg.Payload.Headers[7].Value)

		start := strings.LastIndexAny(string(Rawadr), "<")
		end := strings.LastIndexAny(string(Rawadr), ">")

		for i := len(Rawadr) - 1; i > end; i-- {
			Rawadr = Rawadr[:len(Rawadr)-2]
		}

		for i := 0; i < start; i++ {
			Rawadr = Rawadr[1:]
		}
		Rawadr = Rawadr[1 : len(Rawadr)-1]

		fmt.Println("Message:")

		fmt.Println(msg.Snippet)

		fmt.Println("From: ", string(Rawadr))

		fmt.Println("----------------------------------------------------------------------")
		count--
	}
}
