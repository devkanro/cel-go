package pb

import (
	"testing"

	proto3pb "github.com/google/cel-go/test/proto3pb"
)

func TestFileDescription_GetTypes(t *testing.T) {
	pbdb := NewDb()
	fd, err := pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Error(err)
	}
	expected := []string{
		"google.expr.proto3.test.TestAllTypes",
		"google.expr.proto3.test.TestAllTypes.NestedMessage",
		"google.expr.proto3.test.TestAllTypes.MapStringStringEntry",
		"google.expr.proto3.test.TestAllTypes.MapInt64NestedTypeEntry",
		"google.expr.proto3.test.NestedTestAllTypes"}
	if len(fd.GetTypeNames()) != len(expected) {
		t.Errorf("got '%v', wanted '%v'", fd.GetTypeNames(), expected)
	}
	for _, tn := range fd.GetTypeNames() {
		var found = false
		for _, elem := range expected {
			if elem == tn {
				found = true
				break
			}
		}
		if !found {
			t.Error("Unexpected type name", tn)
		}
	}
	for _, typeName := range fd.GetTypeNames() {
		td, err := fd.GetTypeDescription(typeName)
		if err != nil {
			t.Error(err)
		}
		if td.Name() != typeName {
			t.Error("Indexed type name not equal to descriptor type name", td, typeName)
		}
	}
}

func TestFileDescription_GetEnumNames(t *testing.T) {
	pbdb := NewDb()
	fd, err := pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]int32{
		"google.expr.proto3.test.TestAllTypes.NestedEnum.FOO": 0,
		"google.expr.proto3.test.TestAllTypes.NestedEnum.BAR": 1,
		"google.expr.proto3.test.TestAllTypes.NestedEnum.BAZ": 2,
		"google.expr.proto3.test.GlobalEnum.GOO":              0,
		"google.expr.proto3.test.GlobalEnum.GAR":              1,
		"google.expr.proto3.test.GlobalEnum.GAZ":              2}
	if len(expected) != len(fd.GetEnumNames()) {
		t.Error("Count of enum names does not match expected'",
			fd.GetEnumNames(), expected)
	}
	for _, enumName := range fd.GetEnumNames() {
		if enumVal, found := expected[enumName]; found {
			ed, err := fd.GetEnumDescription(enumName)
			if err != nil {
				t.Error(err)
			} else if ed.Value() != enumVal {
				t.Errorf("Enum did not have expected value. %s got '%v', wanted '%v'",
					enumName, ed.Value(), enumVal)
			}
		} else {
			t.Errorf("Enum value not found for: %s", enumName)
		}
	}
}

func TestFileDescription_GetImportedEnumNames(t *testing.T) {
	pbdb := NewDb()
	fds, err := CollectFileDescriptorSet(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatal(err)
	}
	for _, fd := range fds.GetFile() {
		_, err = pbdb.RegisterDescriptor(fd)
		if err != nil {
			t.Fatal(err)
		}
	}
	imported := map[string]int32{
		"google.expr.proto3.test.ImportedGlobalEnum.IMPORT_FOO": 0,
		"google.expr.proto3.test.ImportedGlobalEnum.IMPORT_BAR": 1,
		"google.expr.proto3.test.ImportedGlobalEnum.IMPORT_BAZ": 2,
	}
	for enumName, value := range imported {
		ed, err := pbdb.DescribeEnum(enumName)
		if err != nil {
			t.Error(err)
		}
		if ed.Value() != value {
			t.Errorf("Got %v, wanted %v for enum %s", ed, value, enumName)
		}
	}
}
