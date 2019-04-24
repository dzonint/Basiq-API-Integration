package basiq

import (
	"basiq/config"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type accessTokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type user struct {
	UserType string `json:"type"`
	Id       string `json:"id"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Links    links  `json:"links"`
}

type job struct {
	ResourceType string  `json:"type"`
	Id           string  `json:"id"`
	Created      string  `json:"created"`
	Updated      string  `json:"updated"`
	Steps        []steps `json:"steps"`
	Links        links   `json:"links"`
}

type steps struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	Result result `json:"result"`
}

type result struct {
	ResultType string `json:"type"`
	Url        string `json:"url"`
}

type links struct {
	Self        string `json:"self"`
	Source      string `json:"source"`
	Next        string `json:"next"`
	Account     string `json:"account"`
	Institution string `json:"institution"`
	Connection  string `json:"connection"`
}

type transactionList struct {
	ResponseType string        `json:"type"`
	Count        int           `json:"count"`
	Size         int           `json:"size"`
	Data         []transaction `json:"data"`
	Links        links         `json:"links"`
}

type transaction struct {
	TransactionType string   `json:"type"`
	Id              string   `json:"id"`
	Status          string   `json:"status"`
	Description     string   `json:"description"`
	PostDate        string   `json:"postDate"`
	TransactionDate string   `json:"transactionDate"`
	Amount          string   `json:"amount"`
	Balance         string   `json:"balance"`
	BankCategory    string   `json:"bankCategory"`
	Account         string   `json:"account"`
	Institution     string   `json:"institution"`
	Connection      string   `json:"connection"`
	Direction       string   `json:"direction"`
	Class           string   `json:"class"`
	SubClass        subClass `json:"subClass"`
	Links           links    `json:"links"`
}

type subClass struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

func getAuth() (accessTokenData, error) {
	postData := []byte(`scope = ` + config.Config.Auth.Scope)
	req, err := http.NewRequest("POST", config.Config.Auth.EndpointUrl, bytes.NewBuffer(postData))
	if err != nil {
		return accessTokenData{}, err
	}

	req.Header.Add("basiq-version", config.Config.Auth.BasiqVersion)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+config.Config.Auth.APIKEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return accessTokenData{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return accessTokenData{}, err
	}

	accessTokenData := accessTokenData{}
	json.Unmarshal(body, &accessTokenData)
	if accessTokenData.AccessToken == "" {
		return accessTokenData, errors.New("access token not found in body response - the error probably occured due to malformed request")
	}

	return accessTokenData, nil
}

func getToken(accessToken *string) {
	for {
		accessTokenData, err := getAuth()
		if err != nil {
			panic(err)
		}

		*accessToken = accessTokenData.AccessToken
		time.Sleep(time.Second * time.Duration(accessTokenData.ExpiresIn))
	}
}

func createUser(usr *user, accessToken *string) error {
	if usr.Email == "" && usr.Mobile == "" {
		return errors.New("email or mobile are required for user creation")
	}

	postData, err := json.Marshal(*usr)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.Config.UserCreation.EndpointUrl, bytes.NewBuffer(postData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+*accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(body), &usr)
	if usr.Id == "" {
		return errors.New("user id not found in body response - the error probably occured due to malformed request")
	}

	return nil
}

func initConnection(loginInfo *[]byte, userId *string, accessToken *string) (job, error) {
	endpointURL := strings.Replace(config.Config.Connect.EndpointUrl, "[USER_ID]", *userId, -1)
	req, err := http.NewRequest("POST", endpointURL, bytes.NewBuffer(*loginInfo))
	if err != nil {
		return job{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+*accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return job{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return job{}, err
	}

	job := job{}
	json.Unmarshal(body, &job)

	return job, nil
}

func getJobStatus(job *job, accessToken *string) error {
	endpointURL := config.Config.Job.EndpointUrl + job.Id
	req, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+*accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, &job)

	return nil
}

func getTransactionsUrl(job *job, accessToken *string) (transactionsUrl string, err error) {
	// Do polling for 2 minutes.
	for i := 0; i < 24; i++ {
		err = getJobStatus(job, accessToken)
		if err != nil {
			return "", err
		}
		if job.Steps[2].Status == "success" {
			return job.Steps[2].Result.Url, nil
		}
		time.Sleep(time.Second * 5)
	}

	return "", errors.New("getJobStatus request timeout")
}

func getTransactions(transactionsUrl *string, accessToken *string) (transactionList, error) {
	endpointURL := config.Config.BaseUrl + *transactionsUrl
	req, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return transactionList{}, err
	}

	req.Header.Add("Authorization", "Bearer "+*accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return transactionList{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return transactionList{}, err
	}

	transactionList := transactionList{}
	json.Unmarshal(body, &transactionList)

	return transactionList, nil
}

func getAverageSpendingBySubClass(transactionData *[]transaction) map[string]float64 {
	subClassValuesMap := make(map[string][]float64)
	for _, val := range *transactionData {
		amount, err := strconv.ParseFloat(val.Amount, 64)
		if err != nil || amount > 0 {
			continue
		}
		subClassValuesMap[val.SubClass.Code+" - "+val.SubClass.Title] = append(subClassValuesMap[val.SubClass.Code+"-"+val.SubClass.Title], amount)
	}

	subClassAverageSpendingMap := make(map[string]float64)
	for ind, val := range subClassValuesMap {
		sum := 0.0
		for _, value := range val {
			sum += math.Abs(value)
		}
		average := sum / float64(len(val))

		subClassAverageSpendingMap[ind] = average
	}

	return subClassAverageSpendingMap
}
