package gateways

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"prolice-slack-bot/types"
	"prolice-slack-bot/types/contract"

	"github.com/joho/godotenv"
)

func PullRequestById(id string) types.PullRequest {
	godotenv.Load(".env")

	baseUrl := os.Getenv("ADO_BASE_URL")
	adoAuth := os.Getenv("ADO_AUTH")
	adoCookie := os.Getenv("ADO_COOKIE")

	url := fmt.Sprintf("%s/pullrequests/%s?api-version=7.0", baseUrl, id)
	method := "GET"

	log.Printf("url %s", url)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", adoAuth)
	req.Header.Add("Cookie", adoCookie)

	log.Printf("Starting Request")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	log.Printf("Request Finished")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject contract.PullRequestResponse
	json.Unmarshal(body, &responseObject)

	log.Printf(responseObject.Status)

	return contract.ToType(responseObject)
}

func PullRequestThreads(id string) {
	godotenv.Load(".env")

	baseUrl := os.Getenv("ADO_BASE_URL")
	adoAuth := os.Getenv("ADO_AUTH")
	adoCookie := os.Getenv("ADO_COOKIE")

	url := fmt.Sprintf("%s/repositories/%s?api-version=7.0", baseUrl, id)
	method := "GET"

	log.Printf("url %s", url)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", adoAuth)
	req.Header.Add("Cookie", adoCookie)

	log.Printf("Starting Request")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	log.Printf("Request Finished")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject contract.PullRequestResponse
	json.Unmarshal(body, &responseObject)

	log.Printf(responseObject.Status)

	//return contract.ToType(responseObject)
}
