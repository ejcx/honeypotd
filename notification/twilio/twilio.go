package twilio

import (
  "os"
  "net/http"
  "net/url"
  "fmt"
  "strings"
  "errors"
)

func Notify(s string) error {
  client := &http.Client{}

  apiEndpoint :="https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"
  sid := os.Getenv("TWILIO_ACCOUNT_SID")
  authToken := os.Getenv("TWILIO_AUTH_TOKEN")
  phone := os.Getenv("TO_PHONE")
  fromPhone := os.Getenv("FROM_PHONE")
  if sid == "" || authToken == "" || phone == "" || fromPhone == ""{
    return errors.New("TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TO_PHONE, and FROM_PHONE envvars must be set")
  }
  postEndpoint := fmt.Sprintf(apiEndpoint, sid)

  data := url.Values{}
  data.Set("To", phone)
  data.Set("Body", s)
  data.Set("From", fromPhone)

  req, err := http.NewRequest("POST", postEndpoint, strings.NewReader(data.Encode()))
  if err != nil {
    return err
  }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.SetBasicAuth(sid, authToken)
  _, err = client.Do(req)
  return err
}
