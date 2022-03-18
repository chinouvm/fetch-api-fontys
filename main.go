// Chinou van Maris - 2022
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/chinouvm/fetch_with_go/util"
	"github.com/scorredoira/email"
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
	fmt.Println("Config ingeladen")

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

	file, _ := os.Create("output.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	for _, person := range data {
		if strings.HasPrefix(person.DisplayName, *surnameArgs) || strings.HasPrefix(person.SurName, *surnameArgs) && len(person.GivenName) > *givennameArgs {
			if person.DisplayName[0:1] == person.SurName[0:1] {
				_, err := file.WriteString("Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n")
				if err != nil {
					log.Fatal(err)
				}
			} else {
				_, err := file.WriteString("Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nDisplay: " + person.DisplayName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n")
				if err != nil {
					log.Fatal(err)
				}
			}
		}	
	}

	m := email.NewMessage("Application Output", "Filters: \n" + "Achternaam begint met: " + *surnameArgs + "\nVoornaam langer dan: " + strconv.Itoa(*givennameArgs) + " letter(s)!\n\n")
	m.From = mail.Address{Name: "Golang Application", Address: config.Email.From}
	m.To = []string{config.Email.To}
	m.AddHeader("Subject", "Output from the Fontys API Fetch!")

	if err := m.Attach("output.txt"); err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", config.Email.From, config.Email.SmtpPassword, config.Email.Mailserver)
	if err := email.Send(config.Email.Mailserver + ":" + config.Email.Mailport, auth, m); err != nil {
		log.Fatal(err)
	}

	if err == nil {
		fmt.Printf("Email succesvol verzonden naar %s", config.Email.To)
	}
}

