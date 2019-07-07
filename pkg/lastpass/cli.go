package lastpass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
)

// https://mholt.github.io/json-to-go/
type SecretResponse []struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Fullname        string `json:"fullname"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	LastModifiedGmt string `json:"last_modified_gmt"`
	LastTouch       string `json:"last_touch"`
	Group           string `json:"group"`
	URL             string `json:"url"`
	Note            string `json:"note"`
}

func VerifyCliExistsOrDie() {
	out, err := sh.Command("which", "lpass").Output()
	if err != nil || "" == string(out) {
		panic("lpass binary not found")
	}
	log.Printf("lpass binary found")
}

// FIXME re-attempt, don't die - Error: HTTP response code said error
func LoginOrDie(username string, password string) {
	// echo <PASSWORD> | LPASS_DISABLE_PINENTRY=1 lpass login --trust <USERNAME>
	out, err := sh.NewSession().SetEnv("LPASS_DISABLE_PINENTRY", "1").Command("echo", password).Command("lpass", "login", "--trust", username).Output()
	if err != nil || "" == string(out) {
		panic("Unable to login: verify credentials")
	}
	log.Printf("Succesfully logged in")
}

func RequestSecret(group string, name string) (SecretResponse, error) {

	fullName := buildFullName(group, name)
	response := SecretResponse{}

	log.Printf("Request secret: [group=%s][name=%s][fullName=%s]", group, name, fullName)

	// lpass show <GROUP>/<NAME> --json --expand-multi
	out, err := sh.Command("lpass", "show", fullName, "--json", "--expand-multi").Output()
	if err != nil {
		return response, fmt.Errorf("invalid secret: [%s] - %s", fullName, err)
	}

	log.Printf("Secret response: %s", out)

	// decode JSON structure into Go structure
	jsonErr := json.Unmarshal([]byte(out), &response)
	if jsonErr != nil {
		return response, fmt.Errorf("invalid response: [%s] - %s", fullName, err)
	}

	log.Printf("JSON secret response size: %d", len(response))

	for secret := range response {
		log.Printf("Secret id: %s", response[secret].Id)
	}

	return response, nil
}

func buildFullName(group string, name string) string {
	var b bytes.Buffer
	if group != "" {
		b.WriteString(group)
		b.WriteString("/")
		b.WriteString(name)
	} else {
		b.WriteString(name)
	}
	return b.String()
}
