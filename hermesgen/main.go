package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"

	"bitbucket.org/ConsentSystems/mango-micro/hermesgen/puller"
	"bitbucket.org/ConsentSystems/mango-micro/hermesgen/pusher"
	"bitbucket.org/ConsentSystems/mango-micro/hermesgen/replier"
	"bitbucket.org/ConsentSystems/mango-micro/hermesgen/requester"
)

// CodeGenerator generates the code for the server and client
type CodeGenerator interface {
	MakeHeader(name string) string
	MakeBody(t *ast.TypeSpec) string
}

func main() {
	var fileToParse string
	var outputStr string
	var genMode string
	flag.StringVar(&fileToParse, "f", "", "File to parse")
	flag.StringVar(&outputStr, "o", "./", "Output location folder")
	flag.StringVar(&genMode, "m", "pp", "generate mode, available: pp(defaul - push pull), rpc (req/rep), ps (pub sub)")

	flag.Parse()

	if fileToParse == "" {
		log.Fatal("No file provided")
	}

	fset := token.NewFileSet()
	fmt.Println("Parsing: ", fileToParse)
	f, err := parser.ParseFile(fset, fileToParse, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Parsing error: ", err)
		return
	}

	var srvGen, clGen CodeGenerator

	switch genMode {
	case "pp":
		srvGen = &puller.Generator{}
		clGen = &pusher.Generator{}
		break
	case "rpc":
		srvGen = &replier.Generator{}
		clGen = &requester.Generator{}
	case "ps":
		log.Fatal("Not implemented yet")
	default:
		log.Fatal("Not implemented yet")
	}

	target := srvGen.MakeHeader(f.Name.Name)
	targetClient := clGen.MakeHeader(f.Name.Name)

	// Adding the registry import
	target = strings.Join([]string{target, `import  "bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"`, "\n"}, "")
	targetClient = strings.Join([]string{targetClient, `import "bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"`, "\n"}, "")

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		// find variable declarations
		case *ast.TypeSpec:
			// which are public
			if t.Name.IsExported() {
				switch t.Type.(type) {
				// and are interfaces
				case *ast.InterfaceType:
					server := srvGen.MakeBody(t)
					target = target + server + "\n"

					client := clGen.MakeBody(t)
					targetClient = targetClient + client
					break
				default:
					//fmt.Println("Non INt: ", t)
				}
			}
		case *ast.ImportSpec:
			target = strings.Join([]string{target, `import `, t.Path.Value, "\n"}, "")
			targetClient = strings.Join([]string{targetClient, `import `, t.Path.Value, "\n"}, "")
		default:
			if reflect.TypeOf(t) != nil {
				//fmt.Println(reflect.TypeOf(t))
			}
		}
		return true
	})

	out, err := os.Create(outputStr + "server.go")
	if err != nil {
		log.Fatal("Impossible to save file: ", err.Error())
	}
	defer out.Close()

	out.WriteString(target)

	outCl, err := os.Create(outputStr + "client.go")
	if err != nil {
		log.Fatal("Impossible to save file: ", err.Error())
	}
	defer outCl.Close()
	outCl.WriteString(targetClient)

	return
}
