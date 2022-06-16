package tsl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 测试校验错误物模型
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

// 测试校验正确物模型 校验物模型实体数据
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
		err = test.Model.ValidateEntity(v)
		if err != nil {
			tLog("Error (%s)\n", err.Error())
		}
	}
	return nil
}

// 测试转换为简单模型
func TestToEntity(t *testing.T) {
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
		return executeToEntityTests(t, path)
	})
	if err != nil {
		t.Errorf("Error (%s)\n", err.Error())
	}
}

func executeToEntityTests(t *testing.T, path string) error {
	var test struct {
		Model *Thing
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
	filename := filepath.Base(path)
	if test.Model == nil {
		t.Logf("file: %s, Expected model but got nil\n", filename)
		return nil
	}
	err = test.Model.ValidateSpec()
	if err != nil {
		t.Logf("Error (%s)\n", err.Error())
		return nil
	}

	//fmt.Println(test.Model.ToEntityString())
	// fmt.Println(test.Model.Random("thing.service.property.set", false))
	// fmt.Println(test.Model.Random("thing.service.property.set", true))
	// fmt.Println(test.Model.Random("thing.event.property.post", false))

	for i := 0; i < 1; i++ {
		result, err := test.Model.Random("thing.event.property.post", false)
		if err != nil {
			t.Errorf("Error (%s)\n", err.Error())
			return err
		}
		fmt.Println(string(result))
	}

	for i := 0; i < 1; i++ {
		result, err := test.Model.Random("thing.service.property.get", false)
		if err != nil {
			t.Errorf("Error (%s)\n", err.Error())
			return err
		}
		fmt.Println(string(result))
	}
	// fmt.Println(test.Model.Random("thing.service.reset", true))

	return nil
}
