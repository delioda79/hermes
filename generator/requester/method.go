package requester

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strings"
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

	resultType := results[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name

	mtdName := mtd.Names[0].Name
	if len(params) == 1 {
		switch pr := params[0].Type.(type) {
		case *ast.StarExpr:
			name := pr.X.(*ast.Ident).Name
			mtdStr = `
// ` + mtdName + ` ...
func (cl *default` + nameSp + `Client) ` + mtdName + `(msg ` + name + `) (*` + resultType + `,error) {

	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	resBts, err := cl.rqstr.Sock().Request("` + nameSp + `.` + mtdName + `", bts)
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
	resBts, err := cl.rqstr.Sock().Request("` + nameSp + `.` + mtdName + `", []byte{})
	if err != nil {
		return nil, err
	}
	rsp := &` + resultType + `{}
	json.Unmarshal(resBts, rsp)
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

		resultType := results[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name

		mtdName := mtd.Names[0].Name
		if len(params) == 1 {
			switch pr := params[0].Type.(type) {
			case *ast.StarExpr:
				name := pr.X.(*ast.Ident).Name
				mtdStr = mtdName + `(msg ` + name + `) (*` + resultType + `,error)`
			default:
				fmt.Printf("Parameter: %+v\n", reflect.TypeOf(pr))

			}
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
