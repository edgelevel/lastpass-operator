package lastpass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/codeskyblue/go-sh"
)

// LastPassSecret represents a LastPass secret
// For more examples see example/lpass-examples.txt
// https://mholt.github.io/json-to-go/
type LastPassSecret struct {
	ID              string `json:"id"`
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

// VerifyCliExistsOrDie verifies that lastpass-cli is properly installed
func VerifyCliExistsOrDie() {
	out, err := sh.Command("which", "lpass").Output()
	if err != nil || "" == string(out) {
		panic(fmt.Sprintf("lpass binary not found: [%s]", err))
	}
	log.Printf("lpass binary found")
}

// Login using lastpass-cli
func Login(username string, password string) error {
        _, err := sh.Command("lpass", "status").Output()
	log.Printf("Checking if already logged in")
	if err != nil {
		log.Printf("Doing login")
		out, err := sh.NewSession().SetEnv("LPASS_DISABLE_PINENTRY", "1").Command("echo", password).Command("lpass", "login", "--trust", username).Output()
		if err != nil || "" == string(out) {
			// sometimes returns error: "Error: HTTP response code said error" even if the credentials are valid
			return fmt.Errorf("verify credentials, unable to login: %s", err)
		}
	}
	log.Printf("Succesfully logged in")
	return nil
}

// Logout using lastpass-cli
func Logout() {
	// logout is not relevant in operator context.
        log.Printf("Should logged out, but doing nothing")
}

// RequestSecrets returns one or more secrets using lastpass-cli
func RequestSecrets(group string, name string) ([]LastPassSecret, error) {

	fullName := buildFullName(group, name)
	secrets := []LastPassSecret{}

	log.Printf("Request secrets: [%s]", fullName)

	// lpass show <GROUP>/<NAME> --json --expand-multi
	out, err := sh.Command("lpass", "show", fullName, "--json", "--expand-multi").Output()
	if err != nil {
		return secrets, fmt.Errorf("invalid secrets: [%s] - %s", fullName, err)
	}

	// uncomment for debug
	//log.Printf("Secret response: %s", out)

	// decode JSON structure into Go structure
	jsonErr := json.Unmarshal([]byte(out), &secrets)
	if jsonErr != nil {
		return secrets, fmt.Errorf("invalid JSON: [%s] - %s", fullName, err)
	}

	log.Printf("Found [%d] secrets", len(secrets))

	return secrets, nil
}

// returns <GROUP>/<NAME> or <NAME>
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
