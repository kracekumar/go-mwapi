package mwapi

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
)

const (
	username = "gomwapi"
	password = "gomwapibot"
)

func TestMWAPiStruct(*testing.T) {
	commons_url := url.URL{
		Scheme: "https",
		Host:   "commons.wikimedia.org",
		Path:   "/w/api.php",
	}
	mwapi := NewMWApi(commons_url)
	fmt.Println(mwapi)
}

func TestLogin(*testing.T) {
	commons_url := url.URL{
		Scheme: "https",
		Host:   "commons.wikimedia.org",
		Path:   "/w/api.php",
	}
	mwapi := NewMWApi(commons_url)
	mwapi.Login(username, password)
}

func TestGet(*testing.T) {
	en_wikipedia_url := url.URL{
		Scheme: "https",
		Host:   "en.wikipedia.org",
		Path:   "/w/api.php",
	}
	params := url.Values{
		"action": {"query"},
		"prop":   {"revisions"},
		"titles": {"Tamil_language"},
		"rvprop": {"content"},
	}
	mwapi := NewMWApi(en_wikipedia_url)
	resp := mwapi.Get(params)
	if resp.StatusCode != 200 {
		rbody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rbody))
		panic("Failed")
	}
}

func TestPost(*testing.T) {
	en_wikipedia_url := url.URL{
		Scheme: "http",
		Host:   "78.46.204.24",
		Path:   "/api.php",
	}
	mwapi := NewMWApi(en_wikipedia_url)
	mwapi.Login(username, password)
	t := mwapi.GetToken("edit")
	params := url.Values{
		"action":  {"edit"},
		"title":   {"Gomwapi-test-page"},
		"section": {"new"},
		"text":    {"This page is created by gomwapi bot. http://github.com/kracekumar/go-mwapi"},
		"token":   {t.Tokens.Edittoken},
	}
	resp := mwapi.PostForm(params)
	if resp.StatusCode != 200 {
		rbody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rbody))
		panic("Failed")
	}
}
