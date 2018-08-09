package utils

import (
	"fmt"
	"os"
	"reflect"

	"go/ast"
)

// GetNameFromTopLevelNode retrieves the name from a top level node
func GetNameFromTopLevelNode(xpr ast.Expr) string {
	mtdStr := ""
	switch pr := xpr.(type) {
	case *ast.StarExpr:
		mtdStr = GetNameFromStarExpr(pr)
	default:
		fmt.Printf("Parameter: %+v\n", reflect.TypeOf(pr))

	}

	return mtdStr
}

// GetNameFromNode retrieves the name from a subnode
func GetNameFromNode(xpr ast.Expr) string {
	name := ""
	switch damn := xpr.(type) {
	case *ast.Ident:
		name = xpr.(*ast.Ident).Name
	case *ast.SelectorExpr:
		slr := xpr.(*ast.SelectorExpr)
		fmt.Printf("Selector: %+v\n", slr)
		name = slr.X.(*ast.Ident).Name + "." + slr.Sel.Name
	default:
		fmt.Println("Type is: ", reflect.TypeOf(damn))
		os.Exit(1)
	}

	return name
}

// GetNameFromStarExpr retrieves the name from a star expression
func GetNameFromStarExpr(xpr *ast.StarExpr) string {
	resultType := ""
	switch damn := xpr.X.(type) {
	case *ast.Ident:
		resultType = xpr.X.(*ast.Ident).Name
	case *ast.SelectorExpr:
		slr := xpr.X.(*ast.SelectorExpr)
		fmt.Printf("Selector: %+v\n", slr)
		resultType = slr.X.(*ast.Ident).Name + "." + slr.Sel.Name
	default:
		fmt.Println("Type is: ", reflect.TypeOf(damn))
		os.Exit(1)
	}

	return resultType
}
