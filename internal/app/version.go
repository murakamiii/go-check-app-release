package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"golang.org/x/net/html"
	"errors"
)

// App ...
type App struct {
	Client *http.Client
}

// GetiOSVersion ...
func (app *App) GetiOSVersion() (string, error) {
	resp, err := app.Client.Get("https://itunes.apple.com/lookup?id=944884603&country=JP")
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
func (app *App) GetAndroidVersion() (string, error) {
	resp, err := app.Client.Get("https://play.google.com/store/apps/details?id=com.coconala.android.portal&hl=ja")
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
			version = n.Parent.NextSibling.FirstChild.FirstChild.FirstChild.Data
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