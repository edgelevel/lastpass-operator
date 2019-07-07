package lastpass

import (
	"github.com/codeskyblue/go-sh"
	"log"
)

func VerifyCliExistsOrDie() {
	out, err := sh.Command("which", "lpass").Output()
	if err != nil || "" == string(out) {
		panic("lpass binary not found")
	}
	log.Printf("lpass binary found")
}

func LoginOrDie(username string, password string) {
	// echo <PASSWORD> | LPASS_DISABLE_PINENTRY=1 lpass login --trust <USERNAME>
	out, err := sh.NewSession().SetEnv("LPASS_DISABLE_PINENTRY", "1").Command("echo", password).Command("lpass", "login", "--trust", username).Output()
	if err != nil || "" == string(out) {
		panic("Unable to login: verify credentials")
	}
	log.Printf("Succesfully logged in")
}
