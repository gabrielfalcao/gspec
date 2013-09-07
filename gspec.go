package main
import (
	"os"
	"fmt"
	"github.com/gabrielfalcao/gspec/filesystem"
	"github.com/gabrielfalcao/gspec/runner"
)

func main() {
	var locationOfSpecs filesystem.Node
	var application runner.Runner
	var err error

	if len(os.Args) == 1 {
		locationOfSpecs, err = filesystem.GetNode(".")
	} else {
		locationOfSpecs, err = filesystem.GetNode(os.Args[1])
	}
	if (err != nil) {
		fmt.Println("ERROR:", err)
		return
	}
	workingDirectory, _ := os.Getwd()

	application = runner.Runner{&locationOfSpecs, workingDirectory}
	application.Run()
}
