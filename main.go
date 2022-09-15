package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var sheetColumns = []string{"A", "B", "C"}

type Config struct {
	RootDirectoryId string `envconfig:"root_directory_id" required:"true"`
	SpreadsheetId   string `envconfig:"spreadsheet_id" required:"true"`
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
		"\n%v\nauthorization code: ", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func makeTableRange(start int, i ...int) string {
	if len(i) == 0 {
		return fmt.Sprintf("%s%s", sheetColumns[start-1], fmt.Sprint(start))
	}

	end := i[0]
	cols := 3

	if len(i) > 1 {
		cols = i[1]
	}
	return fmt.Sprintf("%s%s:%s%s", sheetColumns[0], fmt.Sprint(start), sheetColumns[cols-1], fmt.Sprint(end))
}

func makeRow(vals ...interface{}) []interface{} {
	var iVals []interface{}
	for _, val := range vals {
		iVals = append(iVals, val)
	}
	return iVals
}

func filesRecurSpreadSheet(filesList *drive.FilesListCall, files []*drive.File, vals *[][]interface{}, dir string) {
	for _, f := range files {
		filePath := fmt.Sprintf("%s%s", dir, f.Name)
		*vals = append(*vals, makeRow(f.Name, filePath, f.WebViewLink))
		localQ := fmt.Sprintf("'%s' in parents", f.Id)
		localFiles, _ := filesList.Q(localQ).Fields("nextPageToken, files").Do()
		if len(localFiles.Files) > 0 {
			filesRecurSpreadSheet(filesList, localFiles.Files, vals, fmt.Sprintf("%s%s", filePath, "/"))
		}
	}
}

func filesRecurPrint(filesList *drive.FilesListCall, files []*drive.File, tw *tabwriter.Writer, dir string) {
	for _, f := range files {
		filePath := fmt.Sprintf("%s%s", dir, f.Name)
		rowString := fmt.Sprintf("%s\t%s\t%s", f.Name, filePath, f.WebViewLink)
		fmt.Fprintln(tw, rowString)
		localQ := fmt.Sprintf("'%s' in parents", f.Id)
		localFiles, _ := filesList.Q(localQ).Fields("nextPageToken, files").Do()
		if len(localFiles.Files) > 0 {
			filesRecurPrint(filesList, localFiles.Files, tw, fmt.Sprintf("%s%s", filePath, "/"))
		}
	}
}

func main() {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to read .env: %s", err.Error())
	}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Unable to read config: %s", err.Error())
	}

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	filesList := driveSrv.Files.List()

	qFilter := fmt.Sprintf("'%s' in parents", cfg.RootDirectoryId)

	r, err := filesList.
		Q(qFilter).
		Fields("nextPageToken, files").Do()
	if err != nil {
		log.Fatalf("filesList: %s", err.Error())
	}
	if len(r.Files) == 0 {
		log.Print("Root directory is empty!")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Имя\tПуть\tURL")
	filesRecurPrint(filesList, r.Files, w, "/")
	w.Flush()

	// Integration of
	// Spreadsheets
	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	_, err = sheetsSrv.Spreadsheets.Values.Clear(cfg.SpreadsheetId, makeTableRange(1, 1000), &sheets.ClearValuesRequest{}).Do()

	valueInputOption := "RAW"

	var values [][]interface{}
	values = append(values, makeRow("Имя", "Путь", "URL"))

	filesRecurSpreadSheet(filesList, r.Files, &values, "/")

	data := []*sheets.ValueRange{}
	data = append(data, &sheets.ValueRange{
		Range:  makeTableRange(1, len(values)),
		Values: values,
	})

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,
		Data:             data,
	}

	_, err = sheetsSrv.Spreadsheets.Values.BatchUpdate(cfg.SpreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	fmt.Print("Done. Spreadsheet is successfully updated")
}
