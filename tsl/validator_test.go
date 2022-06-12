package tsl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInvalidModel(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	wd = filepath.Join(wd, "testdata/invalid_model")

	err = filepath.Walk(wd, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() && path != wd {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(fileInfo.Name(), ".json") {
			return nil
		}
		return executeInvalidTests(t, path)
	})
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
	}
}
func executeInvalidTests(t *testing.T, path string) error {
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
		return err
	}
	defer file.Close()

	fmt.Println(file.Name())

	var testThing Thing
	d := json.NewDecoder(file)
	err = d.Decode(&testThing)
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
		return err
	}
	err = testThing.ValidateSpec()
	if err == nil {
		t.Errorf("file: %s, Expected error but got nil\n", filepath.Base(path))
	}
	t.Logf("file: %s, go Expected(%s)\n", filepath.Base(path), err.Error())
	return nil
}

func TestValidModel(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	wd = filepath.Join(wd, "testdata/model")

	err = filepath.Walk(wd, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() && path != wd {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(fileInfo.Name(), ".json") {
			return nil
		}
		return executeValidTests(t, path)
	})
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
	}
}
func executeValidTests(t *testing.T, path string) error {
	var test struct {
		Model    *Thing
		Entities []*ThingEntity
		IsFail   bool
	}

	file, err := os.Open(path)
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
		return err
	}
	defer file.Close()

	fmt.Println(file.Name())

	d := json.NewDecoder(file)
	err = d.Decode(&test)
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
		return err
	}
	// 控制是否是失败测试
	tLog := t.Errorf
	if test.IsFail {
		tLog = t.Logf
	}
	filename := filepath.Base(path)
	if test.Model == nil {
		tLog("file: %s, Expected model but got nil\n", filename)
		return nil
	}
	err = test.Model.ValidateSpec()
	if err != nil {
		tLog("Error (%s)\n", err.Error())
		return nil
	}
	if test.Entities == nil {
		t.Logf("Entities is empty \n")
		return nil
	}
	for _, v := range test.Entities {

		strs := strings.Split(v.Method, ".")
		if len(strs) != 4 {
			tLog("file: %s, Expected method but got %s\n", filename, v.Method)
			continue
		}
		id := strs[2]
		if strings.Compare(id, "property") == 0 {
			id = strs[3]
		}
		switch strs[1] {
		case "service":
			err = test.Model.ValidateService(id, v.Params, v.Data)
		case "event":
			err = test.Model.ValidateEvent(id, v.Params)
		default:
			tLog("file: %s, err method got %s,\n", filename, v.Method)
			return nil
		}
		if err != nil {
			tLog("Error (%s)\n", err.Error())
		}
	}
	return nil
}
