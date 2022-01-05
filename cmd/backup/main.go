package main

import (
	"flag"
	"github.com/echokepler/megad2561/internal/backup"
	"log"
	"os"
)

func main() {
	var (
		addr           = flag.String("host", "http://192.168.88.14", "Сетевой адрес веб сервера megad")
		pwd            = flag.String("password", "sec", "Пароль доступка к megad")
		backupFilePath = flag.String("backup_file", "", "Путь до бэкап файла")
		outputDirPath  = flag.String("output", "", "Путь до папки куда сохранится бэкап")
	)

	flag.Parse()

	var (
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

		needGeneration = len(*outputDirPath) > 0
		needUploading  = len(*backupFilePath) > 0
		mode           backup.WorkMode
	)

	if needGeneration {
		mode = backup.GENERATE
	} else if needUploading {
		mode = backup.UPLOAD
	}

	infoLog.Println("Creating backup megad")
	infoLog.Printf("Address: %v Password: %v Mode: %v", *addr, *pwd, mode)

	if mode == backup.GENERATE {
		backup.GenerateBackup(*outputDirPath, *addr, *pwd)
	} else if mode == backup.UPLOAD {
		backup.UploadBackup(*backupFilePath, *addr, *pwd)
	}
}

// Загрузка бэкапа -> $ backupmegad --addr=192.168.88.14 --password=sec --backup_file=~/WorkSpace/apphub/megad_backup_15032021.yaml
// Создание бэкапа -> $ backupmegad --addr=192.168.88.14 --password=sec --output=~/WorkSpace/apphub
