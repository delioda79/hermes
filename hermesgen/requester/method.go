package requester

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"reflect"
	"strings"

	"bitbucket.org/ConsentSystems/mango-micro/hermesgen/utils"
)

func makeMethod(nameSp string, mtd *ast.Field) string {
	mtdStr := ""
	mtdTp := mtd.Type.(*ast.FuncType)
	params := mtdTp.Params.List
	if len(params) > 1 {
		log.Fatal("Methods can accept maximum 1 parameter")
	}

	results := mtdTp.Results.List
	if len(results) != 2 {
		log.Fatal("Methods need to return exactly 2 params")
	}

	resultType := ""
	switch damn := results[0].Type.(*ast.StarExpr).X.(type) {
	case *ast.Ident:
		resultType = results[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
	case *ast.SelectorExpr:
		slr := results[0].Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
		fmt.Printf("Selector: %+v\n", slr)
		resultType = slr.X.(*ast.Ident).Name + "." + slr.Sel.Name
	default:
		fmt.Println("Type is: ", reflect.TypeOf(damn))
		os.Exit(1)
	}

	mtdName := mtd.Names[0].Name
	if len(params) == 1 {
		switch pr := params[0].Type.(type) {

		case *ast.StarExpr:
			name := ""
			switch damn := pr.X.(type) {
			case *ast.Ident:
				name = pr.X.(*ast.Ident).Name
			case *ast.SelectorExpr:
				slr := pr.X.(*ast.SelectorExpr)
				fmt.Printf("Selector: %+v\n", slr)
				name = slr.X.(*ast.Ident).Name + "." + slr.Sel.Name
			default:
				fmt.Println("Type is: ", reflect.TypeOf(damn))
				os.Exit(1)
			}
			mtdStr = `
// ` + mtdName + ` ...
func (cl *default` + nameSp + `Client) ` + mtdName + `(msg ` + name + `) (*` + resultType + `,error) {

	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	sck := cl.rqstr.Sock()
	sck.SetDeadline(cl.deadline)
	resBts, err := sck.Request("` + nameSp + `.` + mtdName + `", bts)
	if err != nil {
		return nil, err
	}
	resArr := &[]*[]byte{}
	json.Unmarshal(resBts, resArr)
	rsp := &` + resultType + `{}
	json.Unmarshal(*(*resArr)[0], rsp)
	if len(*(*resArr)[1]) > 0 {
		return nil, errors.New(string(*(*resArr)[1]))
	}
	return rsp, nil
}
`

		default:
			fmt.Printf("Parameter: %+v\n", reflect.TypeOf(pr))

		}
	} else {
		mtdStr = `
// ` + mtdName + ` ...
func (cl *default` + nameSp + `Client) ` + mtdName + `() (*` + resultType + `,error) {
	sck := cl.rqstr.Sock()
	sck.SetDeadline(cl.deadline)
	resBts, err := sck.Request("` + nameSp + `.` + mtdName + `", []byte{})
	if err != nil {
		return nil, err
	}
	resArr := &[]*[]byte{}
	json.Unmarshal(resBts, resArr)
	rsp := &` + resultType + `{}
	json.Unmarshal(*(*resArr)[0], rsp)
	if len(*(*resArr)[1]) > 0 {
		return nil, errors.New(string(*(*resArr)[1]))
	}
	return rsp, nil
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
	rqstr requester.Server
	deadline time.Duration
}

// SetDeadline Sets the deadline for the requests
func (cl *default` + nameSp + `Client) SetDeadline(dr time.Duration) {
	cl.deadline = dr
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

		results := mtdTp.Results.List
		if len(results) != 2 {
			log.Fatal("Methods need to return exactly 2 params")
		}

		resultType := utils.GetNameFromTopLevelNode(results[0].Type)

		mtdName := mtd.Names[0].Name
		if len(params) == 1 {
			name := utils.GetNameFromTopLevelNode(params[0].Type)
			mtdStr = mtdName + `(msg ` + name + `) (*` + resultType + `,error)`
		} else {
			mtdStr = mtdName + `() (*` + resultType + `,error)`
		}

		methods = append(methods, mtdStr)
	}
	methodsStr := `
// ` + nameSp + ` ...
type ` + nameSp + ` interface {
	` + strings.Join(methods, "\n	") + `
}
`
	return methodsStr
}
