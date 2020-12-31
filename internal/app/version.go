package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// App ...
type App struct {
	Client *http.Client
}

// GetVersions return {"ios": "x.y.z", "android": "a.b.c"} or error
func (app *App) GetVersions(iosID string, androidID string, appStoreCache bool) (map[string]string, error) {
	v := map[string]string{}

	if len(iosID) > 0 {
		iosv, err := app.GetiOSVersion(iosID, appStoreCache)
		if err != nil {
			return v, err
		}
		v["ios"] = iosv
	}

	if len(androidID) > 0 {
		androidv, err := app.GetAndroidVersion(androidID)
		if err != nil {
			return v, err
		}
		v["android"] = androidv
	}

	return v, nil
}

// GetiOSVersion ...
func (app *App) GetiOSVersion(id string, appStoreCache bool) (string, error) {
	url := fmt.Sprintf("https://itunes.apple.com/lookup?id=%s&country=JP", id)
	if !appStoreCache {
		datetime := strconv.FormatInt(time.Now().Unix(), 10) // キャッシュ対策
		url += fmt.Sprintf("&d=%s", datetime)
	}

	resp, err := app.Client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})["version"].(string), nil
}

// GetAndroidVersion do scraping html
/**
↓を期待している
<div class="hAyfc">
	<div class="BgcNfc">現在のバージョン</div>
	<span class="htlgb">
		<div class="IQ1z0d">
			<span class="htlgb">4.14.2</span>
		</div>
	</span>
</div>
**/
func (app *App) GetAndroidVersion(id string) (string, error) {
	url := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=ja", id)
	resp, err := app.Client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	version := ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && n.Data == "現在のバージョン" { // 日本語決めうち
			// ノードの親の次の子の子の子
			version = emptyOrData(firstChildOrNil(firstChildOrNil(firstChildOrNil(nextSiblingOrNil(n.Parent)))))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	if len(version) > 0 {
		return version, nil
	}
	return "", errors.New("can not find the android version Info")
}

func firstChildOrNil(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}
	return node.FirstChild
}

func nextSiblingOrNil(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}
	return node.NextSibling
}

func emptyOrData(node *html.Node) string {
	if node == nil {
		return ""
	}
	return node.Data
}
