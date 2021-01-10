package googleoauth

import (
	"encoding/json"
	"go_exp/structs"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

type MyGoauth struct {
	OauthConf *oauth2.Config
	AppConf   structs.AppConfig
}

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func getLoginURL(state string, oauthConf *oauth2.Config) string {
	var url = oauthConf.AuthCodeURL(state)
	return url
}

// func to create application's access token
func createAccessToken(email string, secret string) (string, error) {
	// set token expiration time
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	scrt := []byte(secret)
	tokenString, err := token.SignedString(scrt)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// *********** Route handler ******************

func (mygo *MyGoauth) AuthStartHandler(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, getLoginURL("1", mygo.OauthConf))
	// c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL("1") + "'><button>Login with Google!</button> </a> </body></html>"))
}

func (mygo *MyGoauth) AuthCallbackHandler(c echo.Context) error {
	// Handle the exchange code to initiate a transport
	googleToken, err := mygo.OauthConf.Exchange(oauth2.NoContext, c.QueryParam("code"))
	if err != nil {
		log.Println("[Callback] Fail exchange google token")
		return c.JSON(http.StatusBadRequest, err)
	} else {
		// log.Println("token exchange: ", tok)
	}

	client := mygo.OauthConf.Client(oauth2.NoContext, googleToken)
	// send request to google to get user info based on permissions
	googleUserInfoResponse, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println("[Callback] Fail to get userinfo response")
		return c.JSON(http.StatusBadRequest, err)
	}
	defer googleUserInfoResponse.Body.Close()
	data, _ := ioutil.ReadAll(googleUserInfoResponse.Body)

	var userInfo User
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println("[Callback] Error decode userinfo from google: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	// log.Println("userinfo: ", userInfo)
	accessToken, err := createAccessToken(userInfo.Email, mygo.AppConf.Secret)
	if err != nil {
		log.Println("[Callback] Error creating access token: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	// log.Println("accessToken: ", accessToken)
	// send back to client
	return c.Redirect(http.StatusPermanentRedirect, "http://localhost"+mygo.AppConf.ClientPort+"?token="+accessToken)
}
