package replier

import (
	"fmt"
	"go/ast"
	"log"
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

	for _, v := range mtdTp.Results.List {
		fmt.Printf("Method: %+v\n", v)
	}

	if len(mtdTp.Results.List) != 2 {
		log.Fatal("An RPC method ust return exactly two parameters, first being a pointer to struct, second being an error")
	}

	if _, ok := mtdTp.Results.List[0].Type.(*ast.StarExpr); !ok {
		log.Fatal("First returned parameter should be pointer to struct")
	}

	if mtdTp.Results.List[1].Type.(*ast.Ident).Name != "error" {
		log.Fatal("Second returned parameter should be an error")
	}
	mtdName := mtd.Names[0].Name
	if len(params) == 1 {
		switch pr := params[0].Type.(type) {
		case *ast.StarExpr:
			// name := ""
			// switch damn := pr.X.(type) {
			// case *ast.Ident:
			// 	name = pr.X.(*ast.Ident).Name
			// case *ast.SelectorExpr:
			// 	slr := pr.X.(*ast.SelectorExpr)
			// 	fmt.Printf("Selector: %+v\n", slr)
			// 	name = slr.X.(*ast.Ident).Name + "." + slr.Sel.Name
			// default:
			// 	fmt.Println("Type is: ", reflect.TypeOf(damn))
			// 	os.Exit(1)
			// }
			name := utils.GetNameFromTopLevelNode(params[0].Type)
			mtdStr = `
	handler.Add("` + nameSp + `.` + mtdName + `", func(in interface{}, out ...*[]byte) error {
		*out[0] = []byte{}
		*out[1] = []byte{}
		req := &` + name + `{}
		err := json.Unmarshal(in.([]byte), req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		rsp, err := hdl.` + mtdName + `(req)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		bts, err := json.Marshal(rsp)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		*out[0] = bts
		return nil
	})
`

		default:
			fmt.Printf("Parameter: %+v\n", reflect.TypeOf(pr))

		}
	} else {
		mtdStr = `

	handler.Add("` + nameSp + `.` + mtdName + `", func(in interface{}, out ...*[]byte) error {
		*out[0] = []byte{}
		*out[1] = []byte{}
		fmt.Println("RECEIVED HOOK")

		rsp, err := hdl.` + mtdName + `()
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		bts, err := json.Marshal(rsp)
		if err != nil {
			*out[1] = []byte(err.Error())
			return err
		}

		*out[0] = bts
		return nil
	})
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
	return methodsStr
}
