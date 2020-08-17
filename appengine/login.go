package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// Cookie names.
	CookieNameState    = "_scrobble_auth_state" // value: state param (string)
	CookieNameUserInfo = "_scrobble_user_info"  // value: JSON marshaled CookieData (string)

	CookieDomain = AppDomain

	// Maximum age for cookies.
	CookieAgeState    = 10 * time.Minute
	CookieAgeUserInfo = 15 * 24 * time.Hour
)

const (
	redirectParam = "redirect_to"
)

func loginURLWithRedirect(v string) string {
	vals := url.Values{redirectParam: []string{v}}
	return "/login?" + vals.Encode()
}

func logoutURLWithRedirect(v string) string {
	vals := url.Values{redirectParam: []string{v}}
	return "/logout?" + vals.Encode()
}

func determineRedirect(r *http.Request) string {
	if v := r.FormValue(redirectParam); v != "" {
		return v
	}
	return "/"
}

// CookieData is the data in the cookie.
type CookieData struct {
	Email string
}

// GoogleUserInfo is the data provided from Google's userinfo endpoint.
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func isSecureCookieExpired(err error) bool {
	// Hacky way to check for cookie expired, since the securecookie package
	// does not appear to provide a better way.
	return err != nil && strings.Contains(err.Error(), "expired timestamp")
}

// secureCookieCodec costructs a securecookie encoder/decoder for the given keys.
func secureCookieCodec(hashKey, blockKey string) (*securecookie.SecureCookie, error) {
	hash, err := base64.StdEncoding.DecodeString(hashKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash key: %s", err)
	}
	block, err := base64.StdEncoding.DecodeString(blockKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block key: %s", err)
	}
	return securecookie.New(hash, block), nil
}

func stateCookieCodec(hashKey, blockKey string) (*securecookie.SecureCookie, error) {
	s, err := secureCookieCodec(hashKey, blockKey)
	if err != nil {
		return nil, err
	}
	return s.MaxAge(int(CookieAgeState / time.Second)), nil
}

func userInfoCookieCodec(hashKey, blockKey string) (*securecookie.SecureCookie, error) {
	s, err := secureCookieCodec(hashKey, blockKey)
	if err != nil {
		return nil, err
	}
	return s.MaxAge(int(CookieAgeUserInfo / time.Second)), nil
}

// ErrNoUser is returned by currentUser when there is no user logged in.
var ErrNoUser = errors.New("no current user")

type UserInfo CookieData

// currentUser returns the currently logged in user's info. If there is no
// user logged in, the error is ErrNoUser.
func (s *server) currentUser(r *http.Request) (UserInfo, error) {
	userInfoCodec, err := userInfoCookieCodec(s.secret.CookieHashKey, s.secret.CookieBlockKey)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to build user info cookie codec: %s", err)
	}

	c, err := r.Cookie(CookieNameUserInfo)
	if err == http.ErrNoCookie {
		return UserInfo{}, ErrNoUser
	}
	if err != nil && err != http.ErrNoCookie {
		return UserInfo{}, fmt.Errorf("failed to get cookie: %s", err)
	}

	var jsonCookieData string
	err = userInfoCodec.Decode(CookieNameUserInfo, c.Value, &jsonCookieData)
	if isSecureCookieExpired(err) {
		return UserInfo{}, ErrNoUser
	}
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to decode cookie: %s", err)
	}

	var g CookieData
	if err := json.Unmarshal([]byte(jsonCookieData), &g); err != nil {
		return UserInfo{}, fmt.Errorf("failed to json-unmarshal cookie data: %s", err)
	}

	return UserInfo(g), nil
}

// Produces the "state" param value for use during OAuth.
func oauthStateParam() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to read rand: %s", err))
	}
	return base64.StdEncoding.EncodeToString(b)
}

func oauthConfig(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"profile", "email", "openid"},
		RedirectURL:  "https://" + AppDomain + "/googleAuth",
		Endpoint:     google.Endpoint,
	}
}

func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// If the user is already present, just redirect to original destination.
	// (All other cases, continue below.)
	if _, err := s.currentUser(r); err == nil {
		http.Redirect(w, r, determineRedirect(r), http.StatusFound)
		return
	}

	conf := oauthConfig(s.secret.GoogleClientID, s.secret.GoogleClientSecret)
	state := oauthStateParam()

	stateCodec, err := stateCookieCodec(s.secret.CookieHashKey, s.secret.CookieBlockKey)
	if err != nil {
		log.Printf("failed to make state cookie codec: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Set up writing of state cookie.
	encoded, err := stateCodec.Encode(CookieNameState, state)
	if err != nil {
		log.Printf("failed to encode cookie: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieNameState,
		Value:    encoded,
		Domain:   CookieDomain,
		Expires:  time.Now().Add(CookieAgeState),
		Secure:   true,
		HttpOnly: true,
	})

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *server) googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	conf := oauthConfig(s.secret.GoogleClientID, s.secret.GoogleClientSecret)

	stateCodec, err := stateCookieCodec(s.secret.CookieHashKey, s.secret.CookieBlockKey)
	if err != nil {
		log.Printf("failed to make state cookie codec: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	userInfoCodec, err := userInfoCookieCodec(s.secret.CookieHashKey, s.secret.CookieBlockKey)
	if err != nil {
		log.Printf("failed to make user info cookie codec: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	code := r.URL.Query().Get("code")
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("failed to exchange code: %s", err)
		http.Error(w, "bad 'code' value", http.StatusInternalServerError)
		return
	}

	// Verify that incoming state param in URL matches cookie's state value.
	incomingState := r.URL.Query().Get("state")
	c, err := r.Cookie(CookieNameState)
	if err != nil {
		log.Printf("failed to get cookie: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var expectState string
	err = stateCodec.Decode(CookieNameState, c.Value, &expectState)
	if isSecureCookieExpired(err) {
		log.Printf("state cookie expired: %s", err)
		http.Error(w, "error: try logging in again", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("failed to decode state cookie: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if expectState != incomingState {
		log.Printf("state value mismatch: %s != %s", expectState, incomingState)
		http.Error(w, "error: try logging in again", http.StatusBadRequest)
		return
	}

	// Set up deletion of state cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     CookieNameState,
		Value:    "invalidated", // value does not matter
		Domain:   CookieDomain,
		MaxAge:   -1, // delete cookie
		Secure:   true,
		HttpOnly: true,
	})

	// Fetch user info from Google.
	client := conf.Client(ctx, tok)
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		log.Printf("failed to build request: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	req = req.WithContext(ctx)

	rsp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to do request: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer drainAndClose(rsp.Body)

	if rsp.StatusCode/100 != 2 {
		log.Printf("bad status code from google: %d", rsp.StatusCode)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var g GoogleUserInfo
	if err := json.NewDecoder(rsp.Body).Decode(&g); err != nil {
		log.Printf("failed to json-decode google response: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Marshal for writing into cookie.
	b, err := json.Marshal(CookieData{
		Email: g.Email,
	})
	if err != nil {
		log.Printf("failed to json-encode cookie data: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	jsonCookieData := string(b)

	// Set up writing cookie.
	encoded, err := userInfoCodec.Encode(CookieNameUserInfo, jsonCookieData)
	if err != nil {
		log.Printf("failed to encode cookie: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieNameUserInfo,
		Value:    encoded,
		Domain:   CookieDomain,
		Expires:  time.Now().Add(CookieAgeUserInfo),
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound) // TODO: support redirecting to original destination
}

func (s *server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieNameUserInfo,
		Value:    "invalidated", // value does not matter
		Domain:   CookieDomain,
		MaxAge:   -1, // delete cookie
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, determineRedirect(r), http.StatusFound)
}
