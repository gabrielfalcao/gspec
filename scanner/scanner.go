package gspec
import (
	"fmt"
	"bytes"
	"go/ast"
	"go/token"
	"go/parser"
	"go/format"
	"falcao.it/gspec/filesystem"
)

func HandleError(pos token.Position, msg string) {
}

func ParseFile(node filesystem.Node, fset token.FileSet) (imports []byte, specs []byte) {
	if !node.IsGoFile() {
		return []byte{}, []byte{}
	}
	src := node.Read()
	var importBuf bytes.Buffer
	var specBuf bytes.Buffer

	file, err := parser.ParseFile(&fset, node.Path(), src, parser.ParseComments)
	if err != nil {
	 	fmt.Println("ERROR", err)
	}
	phase := 0
	for _, importSpec := range file.Imports {
		var name string
		if importSpec.Name == nil {
			name = ""
		} else {
			name = importSpec.Name.Name
		}
		importBuf.Write([]byte(fmt.Sprintf("import %s %s\n",name, importSpec.Path.Value)))
	}
	PHASES := 8
	ast.Inspect(file, func(n ast.Node) bool {
		if phase > 0 && phase < PHASES {
			phase ++
			return true
		}
		if phase == PHASES {
			format.Node(&specBuf, &fset, n)
			return false
		}

		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "Feature" {
				phase = 1
			}
		}
		return true
	})

	return importBuf.Bytes(), specBuf.Bytes()
}
