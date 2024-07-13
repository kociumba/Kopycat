package internal

import (
	"reflect"
	"runtime"
	"strings"
)

// define possible args
type Args int

const (
	// The full name with all the packages
	FullName Args = iota
	// Only the function name
	OnlyFunc
	// Only the package name
	OnlyPkg
	// Just the function name and the package it has been declared in
	Both
)

// get the local name of the function from a pointer
//
//	args: FullName, OnlyFunc, OnlyPakg, Both
//
// If more than 2 args are passed the function will behave as if no argument were passed
func GetFuncName(f interface{}, arg ...Args) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	switch len(arg) {
	case 0:
		return name
	case 1:
		switch arg[0] {
		case Both:
			return name[strings.LastIndexByte(name, '/')+1:]
		case OnlyPkg:
			slashIndex := strings.LastIndexByte(name, '/')
			dotIndex := strings.LastIndexByte(name, '.')
			if slashIndex < dotIndex {
				return name[slashIndex+1 : dotIndex]
			}
		case OnlyFunc:
			dotIndex := strings.LastIndexByte(name, '.')
			return name[dotIndex+1:]
		}
	}

	return name
}
