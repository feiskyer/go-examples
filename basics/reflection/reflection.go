/*
 * Reflect deep structures of api.PodSpec and export them as grpc proto
 *
 * export GOPATH=~/go/src/k8s.io/kubernetes/Godeps/_workspace/:$GOPATH
 * go run reflection.go
 *
 */
package main

import (
	"fmt"
	"reflect"

	"k8s.io/kubernetes/pkg/api"
	"strings"
)

var typeMapping = map[string]string{
	"int":              "int32",
	"string":           "string",
	"int64":            "int64",
	"bool":             "bool",
	"unversioned.Time": "uint64",
	"IntstrKind":       "string",
	"IntOrString":      "string",
}

var protos = make(map[string]*Proto)

type Proto struct {
	Name   string
	Type   string
	fields []*Proto
}

func reflectData(t reflect.Type) {
	switch t.Kind() {
	case reflect.Struct:
		if proto, ok := protos[t.String()]; ok {
			proto.Type = t.String()
			proto.fields = make([]*Proto, 0, 1)
		} else {
			protos[t.String()] = &Proto{
				Type:   t.String(),
				fields: make([]*Proto, 0, 1),
			}
		}

		proto := protos[t.String()]
		for i := 0; i < t.NumField(); i++ {
			stop := false
			typeField := t.Field(i)
			typeString := typeField.Type.String()
			if strings.HasPrefix(typeString, "util.") {
				typeString = "string"
				stop = true
			} else if typeString == "big.Int" {
				typeString = "uint64"
				stop = true
			}

			p := &Proto{
				Name: typeField.Name,
				Type: typeString,
			}

			proto.fields = append(proto.fields, p)
			if !stop {
				reflectData(typeField.Type)
			}
		}
	case reflect.Map:
		reflectData(t.Key())
		reflectData(t.Elem())
	case reflect.Slice, reflect.Ptr:
		reflectData(t.Elem())
	default:
		//do nothing
	}
}

func convertTypes(inType string) string {
	typeString := strings.Replace(inType, "*", "", -1)
	typeString = strings.Replace(typeString, "api.", "", -1)
	typeString = strings.Replace(typeString, "unversioned.", "", -1)
	typeString = strings.Replace(typeString, "[]", "repeated ", -1)

	if t, ok := typeMapping[typeString]; ok {
		typeString = t
	}

	return typeString
}

func main() {
	reflectData(reflect.TypeOf(&api.PodSpec{}))

	for k, v := range protos {
		fmt.Printf("message %s {\n", k)
		for i, f := range v.fields {
			fmt.Printf("\t%s %s = %d;\n", convertTypes(f.Type), f.Name, i+1)
		}
		fmt.Println("}")
	}
}
