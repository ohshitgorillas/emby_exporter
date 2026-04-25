package emby

import (
	"fmt"
	"net/http"
)

type PackageInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Overview    string `json:"overview"`
	GUID        string `json:"guid"`
	Versions    []struct {
		Version string `json:"version"`
	} `json:"versions"`
}

func (c *Client) GetPackageUpdates() ([]PackageInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Packages/Updates", c.URL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c.ctx)

	res := []PackageInfo{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
