/*
@Time : 2022/11/30 11:01
@Author : lianyz
@Description :
*/

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	authentication "k8s.io/api/authentication/v1beta1"
)

func main() {
	http.HandleFunc("/authenticate", authenticate)

	log.Println(http.ListenAndServe(":3000", nil))
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var tokenReview authentication.TokenReview
	err := decoder.Decode(&tokenReview)
	if err != nil {
		log.Println("[Error]", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status": authentication.TokenReviewStatus{
				Authenticated: false,
			},
		})
		return
	}

	log.Print("receving request")
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenReview.Spec.Token})
	tokenClient := oauth2.NewClient(context.Background(), tokenSource)
	githubClient := github.NewClient(tokenClient)

	user, _, err := githubClient.Users.Get(context.Background(), "")
	if err != nil {
		log.Println("[Error]", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status": authentication.TokenReviewStatus{
				Authenticated: false,
			},
		})
		return
	}

	log.Printf("[Success] login as %s", *user.Login)
	w.WriteHeader(http.StatusOK)
	tokenReviewStatus := authentication.TokenReviewStatus{
		Authenticated: true,
		User: authentication.UserInfo{
			Username: *user.Login,
			UID:      *user.Login,
		},
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind":       "TokenReview",
		"status":     tokenReviewStatus,
	})
}
