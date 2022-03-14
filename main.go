// Chinou van Maris - 2022
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
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
	fmt.Println("API Fetch wordt gestart!")
	surnameArgs := flag.String("achternaamfilter", "M", "Welke letter er gebruikt wordt om achternaam te filteren")
	givennameArgs := flag.Int("naamlengte", 3, "Hoelang de voornaam maximaal mag zijn.")
	flag.Parse()

	config, err := util.LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	} 
	fmt.Println("Config loaded")

	client := &http.Client{}
	req, err := http.NewRequest("GET", config.Api.Address, nil)
	if err != nil {
    log.Fatal(err)
	}

	req.Header = http.Header{
    "Content-Type": []string{"application/json"},
    "Authorization": []string{config.Api.AuthToken},
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

	mailbody := ""

	// Mail Setup from config
	from := config.Email.From
	password := config.Email.SmtpPassword
	toEmail := config.Email.To
	to := []string{toEmail}
	host := config.Email.Mailserver
	port := config.Email.Mailport
	address := host + ":" + port
	subject := "Subject: Automatische verzonden mail vanuit Go Applicatie\n"
	mailbody += "\n\nFilters: \n" + "Achternaam begint met: " + *surnameArgs + "\nVoornaam langer dan: " + strconv.Itoa(*givennameArgs) + " letter(s)!\n\n" + "\n----------------\n"

	// ----------------------------------

	for _, person := range data {
		if strings.HasPrefix(person.DisplayName, *surnameArgs) && len(person.GivenName) > *givennameArgs {
			fmt.Println("Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n")
			mailbody += "Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n"				
		}
	}

	message := []byte(subject + mailbody)
	auth := smtp.PlainAuth("", from, password, host)
	mailerr := smtp.SendMail(address, auth, from, to, message)
	if mailerr != nil {
		fmt.Println("Er is iets fout gegaan met het versturen van de mail! Kijk goed of de config juist is")
	} else {
		fmt.Println("Email verzonden!")
	}
}

