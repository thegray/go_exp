package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"go_exp/googleoauth"
	"go_exp/marshalbehaviour"
	"go_exp/structs"
)

var appConf structs.AppConfig
var cred structs.Credentials
var oauthConf *oauth2.Config

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("[INIT] File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	cfgFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Println("[INIT] Config file not found: ", err)
		os.Exit(1)
	}
	json.Unmarshal(cfgFile, &appConf)

	oauthConf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://127.0.0.1" + appConf.Port + "/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func hi(c echo.Context) error {
	return c.JSON(http.StatusOK, "hi")
}

func main() {
	e := echo.New()
	e.Validator = &structs.CustomValidator{Validator: validator.New()}

	mygoauth := &googleoauth.MyGoauth{OauthConf: oauthConf, AppConf: appConf}

	e.POST("/marshal", marshalbehaviour.MarshalTest)
	e.GET("/healthcheck", hi)
	e.GET("/auth/google", mygoauth.AuthStartHandler)             //route to start auth from client
	e.GET("/auth/google/callback", mygoauth.AuthCallbackHandler) //route to handle return from google

	s := &http.Server{
		Addr:         appConf.Port,
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(5) * time.Second,
	}
	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
