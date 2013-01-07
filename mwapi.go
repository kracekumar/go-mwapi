package mwapi

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// cookiejar for persistent connection
type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

//check error
func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func FetchResponseBody(resp *http.Response) []byte {
	rbody, reader_error := ioutil.ReadAll(resp.Body)
	check_error(reader_error)
	resp.Body.Close()
	return rbody
}

//Struct for decoding json during login
type InnerLogin struct {
	Result, Token, Cookieprefix, Sessionid string
}

type OuterLogin struct {
	Login InnerLogin
}

type BaseToken struct {
	Tokens *Token
}

type Token struct {
	Edittoken  string
	Watchtoken string
}

type MWApi struct {
	//struct contains URL for making requests
	Url url.URL
	/* All the request will be sent using client, using client makes
	   easier to maintain session
	*/
	jar    *Jar
	client *http.Client
	format string
	Tokens *Token
}

func NewMWApi(m_url url.URL) *MWApi {
	jar := new(Jar)
	client := &http.Client{nil, nil, jar}
	tokens := new(Token)
	return &MWApi{m_url, jar, client, "json", tokens}
}

func (m MWApi) Get(params url.Values) *http.Response {
	params.Add("format", m.format)
	m.Url.RawQuery = params.Encode()
	resp, err := m.client.Get(m.Url.String())
	check_error(err)
	return resp
}

func (m MWApi) GetToken(token_type string) *BaseToken {
	params := url.Values{
		"action": {"tokens"},
		"prop":   {"info"},
		"format": {"format"},
		"type":   {token_type},
	}
	params.Encode()
	resp := m.Get(params)
	var t BaseToken
	rbody := FetchResponseBody(resp)
	json_error := json.Unmarshal(rbody, &t)
	check_error(json_error)
	return &t
}

func (m MWApi) SetTokens(token_type string) {
	t := m.GetToken(token_type)
	m.Tokens = t.Tokens
}

func (m MWApi) PostForm(data url.Values) *http.Response {
	data.Add("format", m.format)
	resp, err := m.client.PostForm(m.Url.String(), data)
	check_error(err)
	resp.Body.Close()
	return resp
}

func (m MWApi) Login(username, password string) {
	login_query := url.Values{
		"action":     {"login"},
		"lgname":     {username},
		"lgpassword": {password},
		"format":     {m.format},
	}
	resp, err := m.client.PostForm(m.Url.String(), login_query)
	check_error(err)

	rbody := FetchResponseBody(resp)

	var l OuterLogin
	json_error := json.Unmarshal(rbody, &l)
	check_error(json_error)
	login_query.Add("lgtoken", l.Login.Token)

	second_resp, err := m.client.PostForm(m.Url.String(), login_query)
	check_error(err)
	rbody = FetchResponseBody(second_resp)

	var r OuterLogin
	json_error_2 := json.Unmarshal(rbody, &r)
	check_error(json_error_2)

	if r.Login.Result != "Success" {
		panic("Authentication failed")
	}
}
