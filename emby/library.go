package emby

import (
	"fmt"
	"net/http"
)

type LibraryItems struct {
	MovieCount      int `json:"MovieCount"`
	SeriesCount     int `json:"SeriesCount"`
	EpisodeCount    int `json:"EpisodeCount"`
	GameCount       int `json:"GameCount"`
	ArtistCount     int `json:"ArtistCount"`
	ProgramCount    int `json:"ProgramCount"`
	GameSystemCount int `json:"GameSystemCount"`
	TrailerCount    int `json:"TrailerCount"`
	SongCount       int `json:"SongCount"`
	AlbumCount      int `json:"AlbumCount"`
	MusicVideoCount int `json:"MusicVideoCount"`
	BoxSetCount     int `json:"BoxSetCount"`
	BookCount       int `json:"BookCount"`
	ItemCount       int `json:"ItemCount"`
}

type MediaFolder struct {
	Id             string `json:"Id"`
	Name           string `json:"Name"`
	CollectionType string `json:"CollectionType"`
}

type MediaFoldersResponse struct {
	Items []MediaFolder `json:"Items"`
}

func (c *Client) GetLibraryItems() (*LibraryItems, error) {
	return c.getLibraryCounts("")
}

// GetLibraryItemsForParent returns counts scoped to a single library
// (parent collection). When parentId is empty this is identical to
// GetLibraryItems and returns the server-wide totals.
func (c *Client) GetLibraryItemsForParent(parentId string) (*LibraryItems, error) {
	return c.getLibraryCounts(parentId)
}

func (c *Client) getLibraryCounts(parentId string) (*LibraryItems, error) {
	url := fmt.Sprintf("%s/Items/Counts", c.URL)
	if parentId != "" {
		url = fmt.Sprintf("%s?ParentId=%s", url, parentId)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c.ctx)

	res := LibraryItems{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetMediaFolders() ([]MediaFolder, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Library/MediaFolders", c.URL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c.ctx)

	res := MediaFoldersResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res.Items, nil
}
