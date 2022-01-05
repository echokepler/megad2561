package backup

import (
	"errors"
	"fmt"
	"github.com/echokepler/megad2561/core"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const FileDateFormat = "02_01_2006"

type FileService struct {
	file  *os.File
	store map[string]core.ServiceValues
}

// GeneratePath генерирует имя и путь до файла бекапа
func GeneratePath(dir string) (string, error) {
	dateNow := time.Now()
	fileName := fmt.Sprintf("%v_megad2561.txt", dateNow.Format(FileDateFormat))

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	log.Println(dir, absPath)

	return fmt.Sprintf("%v/%v", absPath, fileName), nil
}

// NewFileService конструктор FileService
//
// Проверяет на наличие файла с которым будет работать и создаст по необходимости
func NewFileService(pathFile string) (*FileService, error) {
	var file *os.File

	log.Println(pathFile)

	_, err := os.Stat(pathFile)

	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(pathFile)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		file, err = os.Create(pathFile)
		if err != nil {
			return nil, err
		}
	}

	fileService := FileService{
		file:  file,
		store: map[string]core.ServiceValues{},
	}

	err = fileService.Parse()
	if err != nil {
		return nil, err
	}

	return &fileService, err
}

func (f *FileService) Parse() error {
	content, err := ioutil.ReadAll(f.file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	if len(lines) > 0 {
		lines = lines[:len(lines)-1] // Remove last empty line
	}

	for index, line := range lines {
		value := line

		sv := core.ServiceValues{}

		f.store[strconv.Itoa(index)], err = sv.Parse(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FileService) Get(params core.ServiceValues) (core.ServiceValues, error) {
	if params.Has("cf") {
		return f.store[params.Get("cf")], nil
	} else if params.Has("pt") {
		return f.store[params.Get("pt")], nil
	}

	return nil, errors.New("not found entity")
}

func (f *FileService) Post(values core.ServiceValues) error {
	targetIds := []string{"cf", "pn"}

	for _, id := range targetIds {
		if !values.Has(id) {
			continue
		}

		f.store[values.Get(id)] = values

		_, err := f.file.WriteString(fmt.Sprintf("%v \n", values.Encode()))
		if err != nil {
			return err
		}
	}

	return nil
}
