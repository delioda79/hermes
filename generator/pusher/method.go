package pusher

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"bitbucket.org/ConsentSystems/mango-micro/generator/utils"
)

func makeMethod(nameSp string, mtd *ast.Field) string {
	mtdStr := ""
	mtdTp := mtd.Type.(*ast.FuncType)
	params := mtdTp.Params.List
	if len(params) > 1 {
		log.Fatal("Methods can accept maximum 1 parameter")
	}

	mtdName := mtd.Names[0].Name
	if len(params) == 1 {
		name := utils.GetNameFromTopLevelNode(params[0].Type)
		mtdStr = `
func (cl *default` + nameSp + `Client) ` + mtdName + `(msg ` + name + `) error {
	bts, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return cl.psh.Push("` + nameSp + `.` + mtdName + `", bts)
}
`
	} else {
		mtdStr = `
func (cl *default` + nameSp + `Client) ` + mtdName + `() error {
	bts := []byte{}
	return cl.psh.Push("` + nameSp + `.` + mtdName + `", bts)
}
`
	}

	return mtdStr
}

func makeMethods(nameSp string, lst *ast.FieldList) string {
	methods := []string{}
	for _, mtd := range lst.List {
		fmt.Printf("Int: %+v\n\n\n\n\n", mtd)
		mtdStr := makeMethod(nameSp, mtd)
		methods = append(methods, mtdStr)
	}
	methodsStr := strings.Join(methods, "\n")

	defaultCl := `
type default` + nameSp + `Client struct{
	psh pusher.Pusher
}

`
	return defaultCl + methodsStr
}

// MakeInterface returns the client interface
func makeInterface(nameSp string, lst *ast.FieldList) string {
	methods := []string{}
	for _, mtd := range lst.List {
		mtdStr := ""
		mtdTp := mtd.Type.(*ast.FuncType)
		params := mtdTp.Params.List
		if len(params) > 1 {
			log.Fatal("Methods can accept maximum 1 parameter")
		}

		mtdName := mtd.Names[0].Name
		if len(params) == 1 {
			name := utils.GetNameFromTopLevelNode(params[0].Type)
			mtdStr = mtdName + `(msg ` + name + `) error`
		} else {
			mtdStr = mtdName + `() error`
		}

		methods = append(methods, mtdStr)
	}
	methodsStr := `
type ` + nameSp + ` interface {
	` + strings.Join(methods, "\n	") + `
}
`
	return methodsStr
}
