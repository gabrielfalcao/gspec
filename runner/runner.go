package runner
import (
	"os"
	"os/exec"
	"fmt"
	"crypto/md5"
	"go/token"
	"github.com/gabrielfalcao/gspec/filesystem"
	. "github.com/gabrielfalcao/gspec/scanner"
)

type SpecSet struct {
	Nodes []filesystem.Node
	TokenFiles token.FileSet
	currentPos int
}

func NewSpecSet (nodes []filesystem.Node) (set SpecSet) {
	fset := token.NewFileSet()
	set = SpecSet{nodes, *fset, 0}
	return
}

type Spec struct {
	Node *filesystem.Node
	Parent *SpecSet
}

func (self *Spec) Parse () (imports []byte, describes[]byte) {
	imports, describes = ParseFile(self.Node, self.Parent.TokenFiles)
	return
}

func (self *SpecSet) Length () int {
	return len(self.Nodes)
}

func (self *SpecSet) Next () (child *Spec) {
	if (self.Length() > self.currentPos) {
		node := &self.Nodes[self.currentPos]
		child = &Spec{node, self}
		self.currentPos++
		return
	}
	return nil
}



type Runner struct {
	RootNode *filesystem.Node
	WorkingDirectory string
}

func ReportError(message string, params ...interface{}){
	fmt.Printf("\033[0;31mERROR: \033[1;37m%s\033[0m\n", fmt.Sprintf(message, params...))
}


func (self *Runner) Run() {
	var err error
	if (err != nil) {
		ReportError("%s", err)
		return
	}

	curnode := filesystem.Node{self.WorkingDirectory}
	files, err := self.RootNode.ListFiles()
	if (err != nil) {
		ReportError("path \"%s\" does not exist", os.Args[1])
		return
	}
	fset := token.NewFileSet()

	for _, file := range files {
		imports, specs := ParseFile(file, *fset)
		if len(specs) == 0 {
			continue
		}
		hash := md5.New()
		hash.Write([]byte(file.Path()))

		specFileName := fmt.Sprintf("gspec_%x.go", hash.Sum(nil))
		specFile, err := self.RootNode.NewFile(specFileName)
		if err != nil {
			ReportError("%s", err)
			return
		}
		specFile.Write([]byte("package main\n\n"))
		specFile.Write(imports)
		specFile.Write([]byte("\nfunc main() {\n"))
		specFile.Write(specs)
		specFile.Write([]byte("\n}\n"))
		specFile.Close()
		specnode := self.RootNode.Join(specFileName)

		fmt.Printf("\033[35mHawkEye\033[1;37m is running \033[0;33m%s\033[0m...\n", curnode.RelPath(file))

		cmd := exec.Command("go", "run", specnode.Path())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		specnode.Delete()
		if err != nil {
			ReportError("exec: %s %s", err, specnode.Path())
			return
		}
	}
}
