package foo
import . "github.com/gabrielfalcao/gspec"

func Feature(){
	Describe("The sum operation", func(){
		var number int

		Given("the number 10", func(){
			number = 10
		})
		When("I add 5", func(){
			number += 5
		})

		Then("It should equal 15", func(){
			Expect(number).To.Equal(15)
			It("Should really be 15", func(){
				Expect(number).To.Equal(15)
			})
		})
	})
}
