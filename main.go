package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// replace your configs here
	secret      = "-----BEGIN PRIVATE KEY-----\nYOUR-PRIVATE-KEY\n-----END PRIVATE KEY-----\n"
	keyId       = "ABC123DEFG"
	teamId      = "DEF123GHIJ"
	clientId    = "com.mytest.app"
	redirectUrl = "www.example.com"
)

// create client_secret
func GetAppleSecret() string {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": keyId,
		},
		Claims: jwt.MapClaims{
			"iss": teamId,
			"iat": time.Now().Unix(),
			// constraint: exp - iat <= 180 days
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"aud": "https://appleid.apple.com",
			"sub": clientId,
		},
		Method: jwt.SigningMethodES256,
	}

	ecdsaKey, _ := AuthKeyFromBytes([]byte(secret))
	ss, _ := token.SignedString(ecdsaKey)
	return ss
}

// create private key for jwt sign
func AuthKeyFromBytes(key []byte) (*ecdsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("token: AuthKey must be a valid .p8 PEM file")
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, err
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	}

	return pkey, nil
}

// do http request
func HttpRequest(method, addr string, params map[string]string) ([]byte, int, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	var request *http.Request
	var err error
	if request, err = http.NewRequest(method, addr, strings.NewReader(form.Encode())); err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response *http.Response
	if response, err = http.DefaultClient.Do(request); nil != err {
		return nil, 0, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, response.StatusCode, nil
}

func main() {
	// replace your code here
	code := "your.code"

	data, status, err := HttpRequest("POST", "https://appleid.apple.com/auth/token", map[string]string{
		"client_id":     clientId,
		"client_secret": GetAppleSecret(),
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectUrl,
	})

	fmt.Printf("%d\n%v\n%s", status, err, data)
}
