package erb

import (
	"bytes"
	"encoding/json"
	"erbBot/bot"
	"erbBot/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const UrlApi = "https://erb.minjust.gov.ua/"
const ContentTypeJson = "application/json;charset=UTF-8"

var nameMethods = []string{
	"listDebtorsEndpoint",
}

type resultSearch struct {
	IsSuccess bool `json:"isSuccess"`
	Results   []struct {
		ID              int64     `json:"ID"`
		RootID          int64     `json:"rootID"`
		LastName        string    `json:"lastName"`
		FirstName       string    `json:"firstName"`
		MiddleName      string    `json:"middleName"`
		BirthDate       time.Time `json:"birthDate"`
		Publisher       string    `json:"publisher"`
		DepartmentCode  string    `json:"departmentCode"`
		DepartmentName  string    `json:"departmentName"`
		DepartmentPhone string    `json:"departmentPhone"`
		Executor        string    `json:"executor"`
		ExecutorPhone   string    `json:"executorPhone"`
		ExecutorEmail   string    `json:"executorEmail"`
		DeductionType   string    `json:"deductionType"`
		VpNum           string    `json:"vpNum"`
	} `json:"results"`
	Rows        int       `json:"rows"`
	RequestDate time.Time `json:"requestDate"`
	IsOverflow  bool      `json:"isOverflow"`
}

type querySearch struct {
	SearchType string `json:"searchType"`
	Paging     string `json:"paging"`
	Filter     struct {
		LastName     string      `json:"LastName"`
		FirstName    string      `json:"FirstName"`
		MiddleName   string      `json:"MiddleName"`
		BirthDate    interface{} `json:"BirthDate"`
		IdentCode    string      `json:"IdentCode"`
		CategoryCode string      `json:"categoryCode"`
	} `json:"filter"`
}

func SearchPerson(firstName, lastName, middleName string) resultSearch {
	var query querySearch

	//prepare query
	query.SearchType = "1"
	query.Paging = "1"
	query.Filter.LastName = lastName
	query.Filter.FirstName = firstName
	query.Filter.MiddleName = middleName
	query.Filter.CategoryCode = ""
	query.Filter.IdentCode = ""
	query.Filter.BirthDate = nil

	queryJson, err := json.Marshal(query)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewReader(queryJson)

	content := requestPost("listDebtorsEndpoint", ContentTypeJson, body)
	var result resultSearch

	err = json.Unmarshal(content, &result)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func SearchAllUser(s models.Storage) {
	users := s.GetAllUsers()
	for _, user := range users {
		data := strings.Split(user.DataSearch, " ")

		if data[0] != "" && data[1] != "" || data[2] != "" {
			result := SearchPerson(data[1], data[0], data[2])
			if result.Rows > 0 {
				user.IsSearching = false
				bot.SendMessage(user.IDChat, "Ви знайденні у БД боржників")
				s.SetIsSearching(&user)
			}
		}
	}
}

func requestPost(method string, contentType string, body io.Reader) []byte {
	url := UrlApi + method
	response, err := http.Post(url, contentType, body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return data
}
