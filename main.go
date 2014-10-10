package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
)

	var from, to, subject string
	var smtpServer, smtpUser, smtpPass, corsUrl string
	var port int

func main() {


	flag.IntVar(&port, "port", 8001, "Port to listen on")	
	flag.StringVar(&from, "from", "", "User that the email will be from")
	flag.StringVar(&to, "to", "", "List of recipients for the email")
	flag.StringVar(&subject, "subject", "Contact Form Details", "subject")
	flag.StringVar(&smtpServer, "smtpserver", "", "Outgoing SMTP Server")
	flag.StringVar(&smtpUser, "smtpuser", "", "SMTP User")
	flag.StringVar(&smtpPass, "smtppass", "", "SMTP password")
	flag.StringVar(&corsUrl, "corsurl", "", "The URL to put in the CORS Access-Control-Allow-Origin header")
	flag.Parse()

	if smtpServer == "" || smtpUser == "" || smtpPass == "" ||
		from == "" || to == "" {
		flag.PrintDefaults()
		return
	}


	_, err := mail.ParseAddress(from) 
	if (err != nil) {
		log.Fatal("Error parsing From address: ", err)
		return
	}

	_, err = mail.ParseAddressList(to) 
	if (err != nil) {
		log.Fatal("Error parsing To addresses: ", err)
		return
	}



	http.HandleFunc("/excusemeihaveaquestion", contactFormHandler)
	var server = ":" + strconv.Itoa(port)

	http.ListenAndServe(server, nil)	

}


func contactFormHandler(w http.ResponseWriter, r *http.Request) {
	if corsUrl != "" {
		w.Header().Set("Access-Control-Allow-Origin", corsUrl)
	}
	w.Header().Set("Content-Type","application/json; charset=utf-8")

	fromAddress, _ := mail.ParseAddress(from) 
	toAddresses, _ := mail.ParseAddressList(to) 

	response := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		response["status"] = "FAILURE"
		response["message"] = err.Error()
		jsonStr, _ := json.Marshal(response)
		fmt.Fprintf(w, string(jsonStr)) 	
		return	
	}

	body := "";
	for k, v := range r.Form {
		if k != "body" {
			body += fmt.Sprintf("%s: %s\r\n", k, v)
		}
	}
	body += "---------------------\r\n"
	body += r.Form.Get("body")


	err = sendContactFormEmail(smtpUser,
		smtpPass,
		smtpServer,
		toAddresses,
		fromAddress,
		subject,
		body,
		);
	if err == nil {
		response["status"] = "SUCCESS"
		jsonStr, _ := json.Marshal(response)
		fmt.Fprintf(w, string(jsonStr))
	} else {
		response["status"] = "FAILURE"
		response["message"] = err.Error()
		jsonStr, _ := json.Marshal(response)
		fmt.Fprintf(w, string(jsonStr)) 		
	}
}



func sendContactFormEmail(smtpUser string, 
	smtpPass string, 
	smtpServer string, 
	to []*mail.Address, 
	from *mail.Address, 
	subject string, 
	body string,) error {

	recipients := ""
	toAddresses := make([]string, 0, 10)
	for _, addr := range to {
		toAddresses = append(toAddresses, addr.Address)
		recipients += fmt.Sprintf("%s, ", addr.String())
	}

	header := make(map[string]string)
	header["From"] = from.String();
	header["To"] = recipients
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
 
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n"
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPass,
		smtpServer,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":25",
		auth,
		from.Address,
		toAddresses,
		[]byte(message),
	)
	if err != nil {
		log.Printf("Error sending mail: %v", err)
		return err
	} else {
		log.Printf("Sent mail to %s without error: ", recipients);
	}
	log.Println("Finished trying to send email")
	return nil

}