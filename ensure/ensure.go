package ensure

import (
	"os"
	"fmt"
	"strings"
)

func fail(description string, items ...interface{}){
	fmt.Printf("\033[1;31mAborted because " + strings.Replace(description, ":", ":\033[0m", -1) + "\033[0m\n", items...)
	os.Exit(1)
}

func failMessageMismatch(expected, got string){
	fail("Expected the error message:\n\033[0;33m%s\n\033[0;31mGot instead:\n\033[0;33m%s\n", expected, got)
}
func ItPanickedWithMessage(expected string) {
    err := recover()
    got := fmt.Sprintf("%s", err)

    if err != nil && got != expected {
	    failMessageMismatch(expected, got)
    }
}
func NoPanicWhen(situation string) {
	err := recover()
	if err != nil {
		fail("%s should not have failed, but it did: \n%s", situation, err)
	}
}
