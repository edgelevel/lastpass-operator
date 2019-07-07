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
