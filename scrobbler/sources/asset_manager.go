package sources

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var NotFoundError = errors.New("no asset found")

const BaseURL = "https://discordapp.com/api/oauth2/applications/%APP_ID%/assets"

type AssetManager struct {
	token     string
	appID     string
	assets    []DiscordAsset
	assetsMux *sync.Mutex
	client    *http.Client
}

type DiscordAsset struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

func NewAssetManager(token, appid string) (*AssetManager, error) {
	ret := &AssetManager{
		token:     token,
		appID:     appid,
		client:    &http.Client{},
		assetsMux: &sync.Mutex{},
	}

	_, err := ret.GetAllAssets()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (am *AssetManager) GetAssetViaID(id string) (*DiscordAsset, error) {
	assets, err := am.GetAllAssets()
	if err != nil {
		return nil, err
	}

	for _, v := range assets {
		if v.ID == id {
			return &v, nil
		}
	}

	return nil, errors.New("no ID found")
}

func (am *AssetManager) GetAssetViaName(name string) (*DiscordAsset, error) {

	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, ".", "", -1)
	name = strings.ToLower(name)

	assets, err := am.GetAllAssets()
	if err != nil {
		return nil, err
	}

	for _, v := range assets {
		if v.Name == name {
			return &v, nil
		}
	}

	return nil, NotFoundError
}

func (am *AssetManager) GetAssetsWithName(name string) (*[]DiscordAsset, error) {

	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, ".", "", -1)
	name = strings.ToLower(name)

	assets, err := am.GetAllAssets()
	if err != nil {
		return nil, err
	}

	var ret []DiscordAsset
	for _, v := range assets {
		if v.Name == name {
			ret = append(ret, v)
		}
	}

	return &ret, nil
}

func (am *AssetManager) GetAssetsOfType(Type int) (*[]DiscordAsset, error) {
	assets, err := am.GetAllAssets()
	if err != nil {
		return nil, err
	}

	var ret []DiscordAsset
	for _, v := range assets {
		if v.Type == Type {
			ret = append(ret, v)
		}
	}

	return &ret, nil

}

func (am *AssetManager) RemoveAssetViaName(name string) error {

	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, ".", "", -1)
	name = strings.ToLower(name)

	assets, err := am.GetAllAssets()
	if err != nil {
		return err
	}

	for _, v := range assets {
		if v.Name == name {
			return am.RemoveAsset(v.ID)
		}
	}

	return NotFoundError
}

func (am *AssetManager) RemoveAsset(id string) error {
	req, err := am.makeRequest("DELETE", am.getBaseURL()+"/"+id, nil)
	if err != nil {
		return err
	}

	_, err = am.executeRequest(req)

	assets, err := am.GetAllAssets()
	if err != nil {
		return err
	}

	am.assetsMux.Lock()
	for k, v := range assets {
		if v.ID == id {
			am.assets = append(assets[:k], assets[k+1:]...)
		}
	}
	am.assetsMux.Unlock()

	return err

}

func (am *AssetManager) AddAsset(name, image string, Type int) (*DiscordAsset, error) {

	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, ".", "", -1)
	name = strings.ToLower(name)

	body := fmt.Sprintf("{\"name\": %q, \"image\": %q, \"type\": %d}", name, image, Type)
	req, err := am.makeRequest("POST", am.getBaseURL(), bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	jsonArr, err := am.executeRequestBody(req)
	if err != nil {
		return nil, err
	}

	asset := &DiscordAsset{}
	if err := json.Unmarshal(jsonArr, asset); err != nil {
		return nil, err
	}

	am.assetsMux.Lock()
	am.assets = append(am.assets, *asset)
	am.assetsMux.Unlock()
	return asset, nil
}

func (am *AssetManager) GetAllAssets() ([]DiscordAsset, error) {
	if am.assets != nil && len(am.assets) > 0 {
		return am.assets, nil
	}

	req, err := am.makeRequest("GET", am.getBaseURL(), nil)
	if err != nil {
		return nil, err
	}

	jsonArr, err := am.executeRequestBody(req)
	if err != nil {
		return nil, err
	}

	assets := &[]DiscordAsset{}
	if err := json.Unmarshal(jsonArr, assets); err != nil {
		return nil, err
	}

	am.assetsMux.Lock()
	am.assets = *assets
	am.assetsMux.Unlock()
	return *assets, nil
}

func (am *AssetManager) SetAppName(name string) error {
	body := fmt.Sprintf(`{"name":%q,"description":"","icon":"","tester_name":"","ask_to_join":"false","spectate":"false","cover_image":"","asset_name":""}`, name)
	req, err := am.makeRequest("PUT", "https://discordapp.com/api/oauth2/applications/"+am.appID, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")

	_, err = am.executeRequestBody(req)
	if err != nil {
		return err
	}

	return nil
}

func (am *AssetManager) executeRequestBody(req *http.Request) ([]byte, error) {
	resp, err := am.executeRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	jsonArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonArr, nil
}

func (am *AssetManager) executeRequest(req *http.Request) (*http.Response, error) {
	resp, err := am.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := am.validResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (am *AssetManager) validResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 209 {
		return fmt.Errorf("status code not 200, instead : %d - %q", resp.StatusCode, resp.Status)
	}

	return nil
}

func (am *AssetManager) makeRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	am.addHeaders(req)

	return req, nil
}

func (am *AssetManager) ClearCache() error {
	am.assetsMux.Lock()
	am.assets = nil
	am.assetsMux.Unlock()
	_, err := am.GetAllAssets()
	return err
}

func (am *AssetManager) addHeaders(r *http.Request) {
	r.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	r.Header.Add("authorization", am.token)
	r.Header.Add("origin", "discordapp.com")
}

func (am *AssetManager) getBaseURL() string {
	return strings.Replace(BaseURL, "%APP_ID%", am.appID, 1)
}
