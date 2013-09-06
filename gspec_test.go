package gspec

import (
	"os"
	"strings"
	"fmt"
	"testing"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func fail(description string, items ...interface{}){
	fmt.Printf("\033[1;31mAborted because " + strings.Replace(description, ":", ":\033[0m", -1) + "\033[0m\n", items...)
	os.Exit(1)
}

func failMessageMismatch(expected, got string){
	fail("Expected the error message:\n\033[0;33m%s\n\033[0;31mGot instead:\n\033[0;33m%s\n", expected, got)
}
func ensureItPanickedWithMessage(expected string) {
    err := recover()
    got := fmt.Sprintf("%s", err)

    if err != nil && got != expected {
	    failMessageMismatch(expected, got)
    }
}
func ensureNoPanicWhen(situation string) {
	err := recover()
	if err != nil {
		fail("%s should not have failed, but it did: \n%s", situation, err)
	}
}

func (s *S) TestBasicAxiom(c *gocheck.C) {
	// My first test axiom, a string should be comparable to another
	EXPECTED_MESSAGE := "expected \"foo\" to equal \"bar\""
	defer ensureItPanickedWithMessage(EXPECTED_MESSAGE)

	Expect("foo").To.Equal("bar")
}

func (s *S) TestEqualsPassing(c *gocheck.C) {
	// Second axiom: comparing two equal strings should just yield

	defer ensureNoPanicWhen("Comparing the same strings")
	Expect("foo").To.Equal("foo")
}

func (s *S) TestNotEqualsFailing(c *gocheck.C) {
	// it might look silly but, congnitively speaking it's a big
	// lift when debugging a hard test :)
	EXPECTED_MESSAGE := "expected \"foo\" to not equal \"foo\""

	defer func(){
		err := recover()
		msg := fmt.Sprintf("%s", err)
		if err != nil && msg != EXPECTED_MESSAGE {
			fail("Expected the error message:\n\033[0;33m%s\n\033[0;31mGot instead:\n\033[0;33m%s\n", EXPECTED_MESSAGE, msg)
		}
	}()
	Expect("foo").Not.To.Equal("foo")
}
