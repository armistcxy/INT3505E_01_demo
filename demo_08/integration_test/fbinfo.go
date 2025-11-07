package integrationtest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type UserInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
}

// GetUserInfo retrieves user information from Facebook Graph API.
func GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	const getUserInfoEndpoint = "https://graph.facebook.com/v22.0/me?fields=id,name,picture"

	u, err := url.Parse(getUserInfoEndpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("access_token", token)
	u.RawQuery = q.Encode()

	// TODO: Call Facebook Graph API to retrieve user information.
	httpReq, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	info := UserInfo{}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	info.PictureURL = gjson.GetBytes(body, "picture.data.url").String()

	return &info, nil
}
