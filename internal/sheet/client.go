package sheet

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	folder      = "folder"
	spreadsheet = "spreadsheet"
)

type Client struct {
	parentID string
	root     string

	driveService  *drive.Service
	sheetsService *sheets.Service

	cachedFileID map[string]string
}

func New(credentialJSON []byte, root string) (*Client, error) {
	ctx := context.TODO()

	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON(credentialJSON))
	if err != nil {
		return nil, err // TODO: wrap error
	}

	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialJSON))
	if err != nil {
		return nil, err // TODO: wrap error
	}

	client := &Client{
		driveService:  driveService,
		sheetsService: sheetsService,
		root:          root,
		cachedFileID:  map[string]string{},
	}

	id, err := client.fetchFileID(folder, root)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	client.parentID = id

	return client, nil
}

// This method will fetch the Google Drive's internal ID for either a folder or a file.
// The kind is the file type (eg/ spreadsheet, folder), see https://developers.google.com/drive/api/guides/mime-types for the full list.
func (c *Client) fetchFileID(kind, name string) (string, error) {
	cacheKey := kind + "||" + name
	if key, ok := c.cachedFileID[cacheKey]; ok {
		return key, nil
	}

	resp, err := c.driveService.Files.
		List().
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Q(fmt.Sprintf("mimeType='%s' and name='%s'", c.mimeType(kind), name)).
		Do()
	if err != nil {
		return "", err // TODO: wrap error
	}

	if len(resp.Files) == 0 {
		return "", &FileNotFoundError{kind: kind, name: name}
	}

	if len(resp.Files) > 1 {
		return "", &DuplicatedFilesError{kind: kind, name: name} // TODO: list the links of folders
	}

	fileID := resp.Files[0].Id
	c.cachedFileID[cacheKey] = fileID

	return fileID, nil
}

func (c *Client) mimeType(kind string) string {
	return "application/vnd.google-apps." + kind
}

func (c *Client) CreateTable(name string, columns []interface{}) error {
	// raise error if the table already exists
	fileID, err := c.fetchFileID(spreadsheet, name)
	if !errors.Is(err, &FileNotFoundError{}) {
		if err == nil && fileID != "" {
			return &DuplicatedTableError{kind: spreadsheet, name: name}
		}
		return err
	}

	// create the table (aka. the spreadsheet in Google Drive)
	file := &drive.File{
		Name:     name,
		MimeType: c.mimeType(spreadsheet),
		Parents:  []string{c.parentID},
	}
	file, err = c.driveService.Files.Create(file).Do()
	if err != nil {
		return err // TODO: wrap error
	}

	// add the header of the columns. Unfortunately, this call cannot be combined with table creation.
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{columns},
	}
	_, err = c.sheetsService.Spreadsheets.Values.Append(file.Id, "Sheet1!A:A", valueRange).ValueInputOption("RAW").Do()
	return err
}

func (c *Client) InsertRows(name string, values map[string][]interface{}) error {
	// find the ID of the table
	fileID, err := c.fetchFileID(spreadsheet, name)
	if err != nil {
		return err // TODO: wrap error
	}
	resp, err := c.sheetsService.Spreadsheets.Values.Get(fileID, "Sheet1!1:1").Do()
	if err != nil {
		return err // TODO: wrap error
	}

	// Say if the table has 3 columns: name, age, city
	// and the input parameter values looks something like { "city": ["Toronto", "New York"], "name": ["Jessica", "Tim"] }
	// we want to morph the rowsToBeAppended to mimic how the rows will look like in excel during an append
	// so it will be [["Jessica", "", "Toronto"], ["Tim", "", "New York"]]
	columnIndexMap := map[string]int{}
	for index, column := range resp.Values[0] {
		columnIndexMap[column.(string)] = index
	}
	rowsToBeAppended := make([][]interface{}, len(values))
	totalNumberOfColumns := len(resp.Values[0])
	for index := range rowsToBeAppended {
		rowsToBeAppended[index] = make([]interface{}, totalNumberOfColumns)
	}
	for columnName, items := range values {
		i := columnIndexMap[columnName]
		for j, item := range items {
			rowsToBeAppended[j][i] = item
		}
	}

	// perform append
	valueRange := &sheets.ValueRange{Values: rowsToBeAppended}
	_, err = c.sheetsService.Spreadsheets.Values.Append(fileID, "Sheet1!A:A", valueRange).ValueInputOption("RAW").Do()
	return err
}

func (c *Client) QueryTable(name string) error {
	// fileID, err := c.fetchFileID(spreadsheet, name)
	// if err != nil {
	// 	return err // TODO: wrap error
	// }
	// resp, err := c.sheetsService.Spreadsheets.Values.Get(fileID, "Sheet1").Do()
	// if err != nil {
	// 	return err // TODO: wrap error
	// }

	// table := newTableFromGoogleSheet(resp.Values)

	return nil
}
