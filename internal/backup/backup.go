package backup

import (
	"github.com/echokepler/megad2561/client"
	"github.com/echokepler/megad2561/configs"
	"os"
)

type WorkMode int

//go:generate stringer -type=WorkMode -output=work_mode_string.go
const (
	_ WorkMode = iota
	GENERATE
	UPLOAD
)

func GenerateBackup(path string, host string, pwd string) {
	options := client.OptionsController{
		Host:     host,
		Password: pwd,
	}
	controller, err := client.NewController(options)
	if err != nil {
		panic(err)
	}

	pathFile, err := GeneratePath(path)
	if err != nil {
		panic(err)
	}

	fileService, err := NewFileService(pathFile)
	if err != nil {
		panic(err)
	}

	controller.ChangeService(fileService)

	err = controller.MegadIDConfig.ChangeSettings(func(config configs.MegaIDSettings) configs.MegaIDSettings {
		return config
	})

	if err != nil {
		handleError(err, fileService)
	}

	for _, port := range controller.Ports.Records {
		_ = controller.Ports.Set(port)
	}

	err = controller.MainConfig.Update(func(settings configs.MainSettings) configs.MainSettings {
		return settings
	})

	if err != nil {
		handleError(err, fileService)
	}
}

// UploadBackup WIP
func UploadBackup(path string, host string, pwd string) {}

func handleError(err error, service *FileService) {
	rErr := os.Remove(service.file.Name())
	if rErr != nil {
		panic(rErr)
	}

	panic(err)
}
