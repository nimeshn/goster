package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	fbApppId       string = "859675590797437"
	fbAppSecret    string = "d14e3f34f8d29f304ccb1fc99dc39316"
	gpClientId     string = "204285913760-l1ffdja39fq76g8srsqueqimlhldtss0.apps.googleusercontent.com"
	gpClientSecret string = "MAdrhu-g2P-JxrrBqXL5Jett"
	gpRedirectUri  string = ""
)

type SignInDetails struct {
	accessToken string
	userId      string
	email       string
	sex         string
	FN          string
	expiresAt   uint64
}

type LoginController struct {
	Name    string
	Queries map[string]map[string]string
}

var (
	loginController *LoginController = &LoginController{
		Name: "LoginController",
		Queries: map[string]map[string]string{
			"mysql": map[string]string{
				"getSessionVars": "SELECT value FROM sessionVariables WHERE key = ?",
				"setSessionVars": "CALL SetSessionVar(?,?)",
				"sesionExists":   "SELECT Id FROM Sessions WHERE Id = ?",
				"signIn":         "CALL SignInUser(?,?,?,?,?,?,?)",
				"signOut":        "DELETE FROM sessions where Id = ?",
			},
		},
	}
)

//use cURL to execute get request
func httpGetResult(url string) (rsp string) {
	resp, err := http.Get(url)
	Check(err)
	defer resp.Body.Close()
	rspBytes, err := ioutil.ReadAll(resp.Body)
	Check(err)
	rsp = string(rspBytes)
	return
}

//Gets Session Variables
func (c *LoginController) GetSessionVar(key string) (value string, err error) {
	err = GetDB().QueryRow(c.Queries["mysql"]["getSesionVar"], key).
		Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	fmt.Println("LoginController.GetSessionVar executed", key)
	return
}

//Sets Session Variable
func (c *LoginController) SetSessionVar(key, value string) (ok bool, err error) {
	_, err = GetDB().Exec(c.Queries["mysql"]["setSessionVar"], key, value)
	if err != nil {
		return
	}
	fmt.Println("LoginController.SetSessionVar executed", key, value)
	return true, nil
}

func (c *LoginController) ValidateSessionToken(token string) (ok bool, err error) {
	var value string
	err = GetDB().QueryRow(c.Queries["mysql"]["sessionExists"], token).
		Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	if value != "" {
		ok = true
	}
	fmt.Println("LoginController.ValidateSessionToken executed", token)
	return
}

//SignIn using fbTokenId
func (c *LoginController) SignIn(postedTokens map[string]string) (result map[string]string, err error) {
	var sessionId, memberId, NewSignUp, firstName, token, tokenType, fbToken, gpToken string
	var usrDetail *SignInDetails
	fbToken = postedTokens[`fbToken`]
	gpToken = postedTokens[`gpToken`]

	if fbToken != "" {
		token = fbToken
		tokenType = `FB`
		usrDetail, err = c.VerifyFBAccessToken(fbToken)
		Check(err)
	} else if gpToken != "" {
		token = gpToken
		tokenType = `GP`
		usrDetail, err = c.VerifyGPAccessToken(gpToken)
		Check(err)
	} else {
		err = errors.New(`Incomplete signin details.`)
		return
	}
	//
	if usrDetail == nil {
		err = errors.New(`failed fetching signin details.`)
		return
	}

	err = GetDB().QueryRow(c.Queries["mysql"]["signIn"], token, tokenType,
		usrDetail.userId, usrDetail.email, usrDetail.sex, usrDetail.FN, usrDetail.expiresAt).
		Scan(&sessionId, &memberId, &NewSignUp, &firstName)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Invalid login.")
		}
		return
	}

	//if we got a valid sessionId
	if sessionId != "" {
		result = map[string]string{
			"sessionId": sessionId,
			"memberId":  memberId,
			"NewSignUp": NewSignUp,
			"FN":        firstName,
		}
	}
	return
}

//SignOut using accessToken
func (c *LoginController) SignOut(token string) (ok bool, err error) {
	_, err = GetDB().Exec(c.Queries["mysql"]["signOut"], token)
	if err != nil {
		return
	}
	fmt.Println("LoginController.SignOut executed")
	return true, nil
}

/********************************************************************
	Start Facebook signin related routines
********************************************************************/
//Get the Fb App Access Token
func (c *LoginController) GetFBAppAccessToken() (fbAppAccessToken string, err error) {
	fbAppAccessToken, err = c.GetSessionVar("fbAppAccessToken")
	Check(err)
	if fbAppAccessToken != "" {
		return
	} else {
		url := fmt.Sprintf("https://graph.facebook.com/oauth/access_token?client_id=%s&client_secret=%s&grant_type=client_credentials",
			fbApppId, fbAppSecret)
		response := httpGetResult(url)
		respArr := strings.Split(response, "=")
		//if successfull then it would be like access_token=<access token value>
		if len(respArr) == 2 {
			c.SetSessionVar("fbAppAccessToken", respArr[1])
			fbAppAccessToken = respArr[1]
		} else { //if failed, then it would be JSON error object
			var rspJson map[string]interface{}
			Check(json.Unmarshal([]byte(response), &rspJson))
			if rspJson != nil && rspJson["error"] != nil {
				err = errors.New(response)
			}
		}
	}
	return
}

func (c *LoginController) GetFBUserDetails(fbToken string) (rspJson map[string]interface{}) {
	response := httpGetResult(fmt.Sprintf("https://graph.facebook.com/me?fields=email,gender,name&access_token=%s", fbToken))
	Check(json.Unmarshal([]byte(response), &rspJson))
	return
}

//Verify FB access token
func (c *LoginController) VerifyFBAccessToken(fbToken string) (signDetail *SignInDetails, err error) {
	var rspJson map[string]interface{}
	appToken, err := c.GetFBAppAccessToken()
	Check(err)
	response := httpGetResult(fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s", fbToken, appToken))
	Check(json.Unmarshal([]byte(response), &rspJson))
	if rspJson != nil {
		data := rspJson["data"].(map[string]interface{})
		if !data["is_valid"].(bool) { //checking if access token is invalid if yes, then throw exception
			rspErr := data["error"].(map[string]interface{})
			err = errors.New(rspErr["message"].(string))
			return
		} else { //checking if access token is valid, then get the userId
			//if the token doesn"t belong to our fbApp
			if data["app_id"].(string) != fbApppId {
				err = errors.New("Invalid login.")
				return
			}
		}
	} else {
		err = errors.New("Invalid fbToken.")
		return
	}
	fbuser := c.GetFBUserDetails(fbToken)
	signDetail = &SignInDetails{
		accessToken: fbToken,
		userId:      fbuser["id"].(string),
		email:       fbuser["email"].(string),
		sex:         "M",
		FN:          fbuser["name"].(string),
		expiresAt:   rspJson["data"].(map[string]interface{})["expires_at"].(uint64),
	}
	if fbuser["gender"].(string) != "male" {
		signDetail.sex = "F"
	}
	return
}

/********************************************************************
	End Facebook signin related routines
********************************************************************/

/********************************************************************
Start Google+ signin related routines
********************************************************************/
func (c *LoginController) VerifyGPAccessToken(gpToken string) (userDetail *SignInDetails, err error) {
	config := &oauth2.Config{
		ClientID:     gpClientId,
		ClientSecret: gpClientSecret,
		// Scope determines which API calls you are authorized to make
		Scopes:   []string{"email"},
		Endpoint: google.Endpoint,
		// Use "postmessage" for the code-flow for server side apps
		RedirectURL: gpRedirectUri,
	}

	tok, err := config.Exchange(oauth2.NoContext, gpToken)
	if err != nil {
		err = fmt.Errorf("Error while exchanging code: %v", err)
		return
	}
	// TODO: return ID token in second parameter from updated oauth2 interface
	// tok.AccessToken, tok.Extra("id_token").(string), nil

	if tok.AccessToken != "" {
		userDetail = &SignInDetails{
			accessToken: gpToken,
			userId:      tok.Extra("sub").(string),
			email:       tok.Extra("email").(string),
			sex:         "M", //Ticket does not contain gender info
			FN:          tok.Extra("name").(string),
			expiresAt:   tok.Extra("exp").(uint64),
		}
	}
	return
}

/********************************************************************
	End Google+ signin related routines
********************************************************************/
