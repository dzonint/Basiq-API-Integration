package basiq

import (
	"basiq/config"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_UserAverageSpendingBySubClass(t *testing.T) {
	var accessToken string
	go getToken(&accessToken)
	time.Sleep(time.Second * 2) // Kinda silly, but need to ensure we get the token before we proceed.

	user := user{
		Email:  "gavin@hooli.com",
		Mobile: "+61410888666",
	}

	err := createUser(&user, &accessToken)
	if err != nil {
		t.Error(err)
	}

	loginInfo, _ := json.Marshal(config.Config.Login)
	job, err := initConnection(&loginInfo, &user.Id, &accessToken)
	if err != nil {
		t.Error(err)
	}

	transactionsUrl, err := getTransactionsUrl(&job, &accessToken)
	if err != nil {
		t.Error(err)
	}

	transactionList, err := getTransactions(&transactionsUrl, &accessToken)
	if err != nil {
		t.Error(err)
	}

	subClassAverageSpendingMap := getAverageSpendingBySubClass(&transactionList.Data)

	for ind, val := range subClassAverageSpendingMap {
		fmt.Println(ind, ":", val)
	}
}
