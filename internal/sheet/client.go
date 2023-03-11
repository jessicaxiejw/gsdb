package sheet

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	defaultFolderName = "gsdb" // TODO: configurable folder name

	folder      = "application/vnd.google-apps.folder"
	spreadsheet = "application/vnd.google-apps.spreadsheet"
)

type Client struct {
	parentID string

	driveService  *drive.Service
	sheetsService *sheets.Service
}

func New(credentialJSON []byte) (*Client, error) {
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
	}

	id, err := client.fetchFileID(folder, defaultFolderName)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	if id == "" {
		return nil, fmt.Errorf("cannot find a shared folder named %s. Did you forget to share the folder with %s", defaultFolderName, "") // TODO: wrap error, add the email of the service account
	}
	client.parentID = id

	return client, nil
}

// This method will fetch the Google Drive's internal ID for either a folder or a file.
// The kind is the file type (eg/ spreadsheet, folder), see https://developers.google.com/drive/api/guides/mime-types for the full list.
func (c *Client) fetchFileID(kind, name string) (string, error) {
	resp, err := c.driveService.Files.
		List().
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Q(fmt.Sprintf("sharedWithMe and mimeType='%s' and name='%s'", c.mimeType(kind), name)).
		Do()
	if err != nil {
		return "", err // TODO: wrap error
	}

	if len(resp.Files) == 0 {
		return "", fmt.Errorf("cannnot find %s named %s", kind, name)
	}

	if len(resp.Files) > 1 {
		return "", fmt.Errorf("there are more than 1 %s shared with this service account that is named %s.\nPlease ensure there are only 1 %s named %s", kind, defaultFolderName, kind, defaultFolderName) // TODO: wrap error, list the links of folders
	}

	return resp.Files[0].Id, nil
}

func (c *Client) mimeType(kind string) string {
	return "application/vnd.google-apps." + kind
}

func (c *Client) CreateTable(title string, columns []interface{}) error {
	file := &drive.File{
		Name:     title,
		MimeType: c.mimeType(spreadsheet),
		Parents:  []string{c.parentID},
	}
	file, err := c.driveService.Files.Create(file).Do()
	if err != nil {
		return err // TODO: wrap error
	}

	columns = append([]interface{}{"id"}, columns...)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{columns},
	}
	_, err = c.sheetsService.Spreadsheets.Values.Append(file.Id, "Sheet1!A:A", valueRange).ValueInputOption("RAW").Do()
	return err
}

func (c *Client) InsertRows(name string, columnNames []string, values [][]interface{}) error {
	return nil
}
