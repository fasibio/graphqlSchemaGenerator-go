package main

import (
	"reflect"
	"testing"
)

func Test_getSchemaObj(t *testing.T) {

	schemas := []string{
		`type test  
{
  id: String
  name: String
}`,
	}
	result := getSchemaObj(schemas[0])
	if result.name != "test" {
		t.Errorf("getSchemaObj() = %v, want %v", result.name, "test")
	}
	type args struct {
		schemaStr string
	}
	tests := []struct {
		name string
		args args
		want schema
	}{
	// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSchemaObj(tt.args.schemaStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSchemaObj() = %v, want %v", got, tt.want)
			}
		})
	}
}
