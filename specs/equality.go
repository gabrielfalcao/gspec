package specs

import (
	. "github.com/gabrielfalcao/gspec"
	Ensure "github.com/gabrielfalcao/gspec/ensure"
)


func Feature(){
	Describe("The most basic axiom", func(){
		var source, destination string

		Given("the string 'foo'", func(){
			source = "foo"
		})
		When("it gets comparated with the string 'bar'", func(){
			destination = "bar"
		})
		Then("the error message should be of easy cognition", func(){
			defer Ensure.ItPanickedWithMessage("expected \"foo\" to equal \"bar\"")
			Expect("foo").To.Equal("bar")
		})
	})
}
