// Chinou van Maris - 2022
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/chinouvm/fetch_with_go/util"
)

type fhictData struct {
	GivenName string `json:"givenName"`
	SurName string `json:"surName"`
	DisplayName string `json:"displayName"`
	TelephoneNumber string `json:"telephoneNumber"`
}


func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.ApiAddress, nil)
	if err != nil {
    log.Fatal(err)
	}

	req.Header = http.Header{
    "Content-Type": []string{"application/json"},
    "Authorization": []string{config.ApiAuth},
	}	

	
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var data []fhictData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	surnameArgs := "M"
	givennameArgs := 3
	mailbody := ""
	
	for i := 1; i < len(os.Args); i++ {

			num, err := strconv.Atoi(os.Args[i])
			if err != nil {
				surnameArgs = os.Args[2]
			} else {
				givennameArgs = num
			}
	}

	// Mail Setup
	from := config.FromEmail
	password := config.SMTPPassword
	toEmail := config.ToEmail
	to := []string{toEmail}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := "Subject: Automatische verzonden mail vanuit Go Applicatie\n"
	mailbody += "\n\nFilters: \n" + "Achternaam begint met: " + surnameArgs + "\nVoornaam langer dan: " + strconv.Itoa(givennameArgs) + " letter(s)!\n\n" + "\n----------------\n"

	// ----------------------------------

	for _, person := range data {
		if strings.HasPrefix(person.DisplayName, surnameArgs) && len(person.GivenName) > givennameArgs {
			fmt.Println("Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n")
			mailbody += "Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n"				
		}
	}

	message := []byte(subject + mailbody)
	auth := smtp.PlainAuth("", from, password, host)
	mailerr := smtp.SendMail(address, auth, from, to, message)
	if mailerr != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("Email verzonden!")
	}
}

