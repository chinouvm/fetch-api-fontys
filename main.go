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

func handleMail(config util.Config, data []fhictData, surnameArgs string, givennameArgs int) error {
	messageSlice := make([]byte, 1)

	for _, person := range data {
		if strings.HasPrefix(person.DisplayName, surnameArgs) || strings.HasPrefix(person.SurName, surnameArgs) && len(person.GivenName) > givennameArgs {
			if person.DisplayName[0:1] == person.SurName[0:1] {
					messageInBytes := []byte("Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n")
					messageSlice = append(messageSlice, messageInBytes...)

			} else {
				messageInBytes := []byte("Voornaam: " + person.GivenName + "\nAchternaam: " + person.SurName + "\nDisplayname: " + person.DisplayName + "\nTelefoon: " + person.TelephoneNumber + "\n----------------\n")
					messageSlice = append(messageSlice, messageInBytes...)
			}
		}	
	}

	m := email.NewMessage("Application Output", "Filters: \n" + "Achternaam begint met: " + surnameArgs + "\nVoornaam langer dan: " + strconv.Itoa(givennameArgs) + " letter(s)!\n\n")
	m.From = mail.Address{Name: "Golang Application", Address: config.Email.From}
	m.To = []string{config.Email.To}
	m.AddHeader("Subject", "Output from the Fontys API Fetch!")

	if err := m.AttachBuffer("byteslice.txt", messageSlice, false); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", config.Email.From, config.Email.SmtpPassword, config.Email.Mailserver)
	if err := email.Send(config.Email.Mailserver + ":" + config.Email.Mailport, auth, m); err != nil {
		return err
	} else if err == nil {
		fmt.Printf("Email succesvol verzonden naar %s", config.Email.To)
	}

	return nil
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

	handleMail(config, data, *surnameArgs, *givennameArgs)
}

