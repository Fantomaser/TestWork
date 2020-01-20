package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

type Message struct {
	Id       string `json: "id"`
	LabelIds string `json: "labelIds"`
	Raw      []byte `json: "raw"`
}

/*type Messages struct {
	Messages           string `json: "messages"`
	NextPageToken      string `json: "nextPageToken"`
	ResultSizeEstimate []byte `json: "resultSizeEstimate"`
}*/

//GmailStream ... Call gmail stream
func GmailStream() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
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

	res, _ := srv.Users.Messages.List(user).Q("subject:Make message").Do()

	tok := res.Messages
	i := res.ResultSizeEstimate
	//Msg := MSGList[]

	for _, l := range res.Messages {
		fmt.Println(l, " -- ", l.Id)
		fmt.Println(l.Raw)
	}

	for i > 0 {
		//tok[i].Header.Clone().Get()
		fmt.Println(i, " ----  ", tok[i-1].Header.Get(tok[i-1].Id), "\n")
		msg, _ := srv.Users.Messages.Get(user, tok[i-1].Id).Format("raw").Do()

		//dec := new(mime.WordDecoder)
		///var RawURLEncoding = URLEncoding.WithPadding(NoPadding)
		fileData, err := base64.StdEncoding.DecodeString(msg.Raw)
		fileData2, err := base64.StdEncoding.DecodeString(string(fileData))
		if err != nil {
			fmt.Println( /*"problem", err*/ )
		}
		fmt.Println(string(fileData))
		fmt.Println("=================")
		fmt.Println(string(fileData2))
		fmt.Println("=================================")

		fmt.Println("=================")
		str := "bJGfmMjrfapK/hQDUFcm9pQ2TZWBNx6LTvIkDapInKS4aS38e09vALBXXHUzK/VXaygzf0jxotc+RO7f7LixC/4Xr5gUGsZbLeZ88kgNsPub+7tDhehetMgsg+T2qoZg673mcXOKOvRkAd4T/o7GG2d8MfsGfOzTEd8kCcp32AZ8GaTbFkBm1ZrAVZ/tIA1eS/ZrgOvCZf8bGEOeoc1UT9oxAgcaXQ9IKnTzNIuGcRedbw98GHCUWMHdivfn4mhL4z+v+96elzsR82s3s3abwNBJeAKuraKPBgsnby0q/kbtJv3eaMiX/x1OGsWpEifgcCKe1dHNHMRvjt27E5vmVA=="
		fileData3, err := base64.StdEncoding.DecodeString(str)
		fmt.Println(string(fileData3))
		fmt.Println("=================================")

		fmt.Println("\n\n\n\n", msg.Header.Get(tok[i-1].Id), "\n\n\n\n")

		//fmt.Println(msg.Raw)
		i--
	}
	return
}
