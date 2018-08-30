package configor

import (
	"os"
	"reflect"
	"testing"
)

type testConfig struct {
	Test1 int
	Test2 []struct {
		Test2Ele1 int
		Test2Ele2 int
	}
}

func TestLoadSliceFromEnv(t *testing.T) {
	var tc = testConfig{
		Test1: 1,
		Test2: []struct {
			Test2Ele1 int
			Test2Ele2 int
		}{
			{
				Test2Ele1: 1,
				Test2Ele2: 2,
			},
			{
				Test2Ele1: 3,
				Test2Ele2: 4,
			},
		},
	}

	var result testConfig
	os.Setenv("CONFIGOR_TEST1", "1")
	os.Setenv("CONFIGOR_TEST2_0_TEST2ELE1", "1")
	os.Setenv("CONFIGOR_TEST2_0_TEST2ELE2", "2")

	os.Setenv("CONFIGOR_TEST2_1_TEST2ELE1", "3")
	os.Setenv("CONFIGOR_TEST2_1_TEST2ELE2", "4")
	err := Load(&result)
	if err != nil {
		t.Fatalf("load from env err:%v", err)
	}

	if !reflect.DeepEqual(result, tc) {
		t.Fatalf("unexpected result:%+v", result)
	}
}
