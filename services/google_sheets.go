package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/GodCratos/google_sheets/configs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GoogleSheetsWriteDataInSheet() error {
	srv, err := GoogleSheetsConnectToSheet()
	if err != nil {
		return err
	}
	arrayOrders := []interface{}{""}
	page := 1
	rangeIndex := 3
	var vr sheets.ValueRange
	for arrayOrders != nil {
		arrayOrders, err := RetailGetOrdersByPages(page)
		if err != nil {
			return err
		}
		for _, value := range arrayOrders {
			structForGet := RetailStructGenerationForGoogleSheets(value.(map[string]interface{}))
			vr.Values = append(vr.Values, structForGet)
		}
		if page%1 == 0 {
			writeRange := fmt.Sprintf("%s!A%v", configs.GoogleSheetsGetSheetsName(), rangeIndex)
			_, err = srv.Spreadsheets.Values.Update(configs.GoogleSheetsGetSheetsID(), writeRange, &vr).ValueInputOption("RAW").Do()
			if err != nil {
				log.Println("Unable to retrieve data from sheet. ", err)
				return err
			}
			rangeIndex = rangeIndex + 10000
			vr.Values = nil
		}
		page++
	}
	return nil
}

func GoogleSheetsConnectToSheet() (*sheets.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Println("Unable to read client secret file: ", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Println("Unable to parse client secret file to config: ", err)
		return nil, err
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("Unable to retrieve Sheets client: ", err)
		return nil, err
	}
	return srv, nil
}

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
