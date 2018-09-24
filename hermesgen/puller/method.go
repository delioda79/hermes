package puller

import (
	"go/ast"
	"log"
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

	mtdName := mtd.Names[0].Name
	if len(params) == 1 {
		name := utils.GetNameFromTopLevelNode(params[0].Type)
		mtdStr = `
	handler.Add(serviceNmsp + ".` + mtdName + `", func(msg interface{}, rsp ...*[]byte) error {
		inParam := &` + name + `{}
		arg, ok := msg.([]byte)

		if !ok {
			fmt.Printf("Wrong message sent %v", arg)
			return fmt.Errorf("Wrong message sent %v", arg)
		}

		err := json.Unmarshal(arg, inParam)
		if err != nil {
			fmt.Println("Error unmarshaling: ", err)
			return err
		}

		hdl.` + mtdName + `(inParam)
		return nil
	})
`
	} else {
		mtdStr = `
	handler.Add(serviceNmsp + ".` + mtdName + ` ", func(msg interface{}, rsp ...*[]byte) error {
		hdl.` + mtdName + `()
		return nil
	})
`
	}

	return mtdStr
}

func makeMethods(nameSp string, lst *ast.FieldList) string {
	methods := []string{}
	for _, mtd := range lst.List {
		//fmt.Printf("Int: %+v\n\n\n\n\n", mtd)
		mtdStr := makeMethod(nameSp, mtd)
		methods = append(methods, mtdStr)
	}
	methodsStr := strings.Join(methods, "\n")
	return methodsStr
}
