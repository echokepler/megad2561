# Megad2561
Реализует программный интерфейс для работы с контроллером megad2561

## Примеры

### Смена megadID
```go
package main

import (
	"github.com/echokepler/megad2561"
)

func main() {
	options := megad2561.OptionsController{
		Host: "192.168.88.14",
		Password: "sec",
	}
	controller, err := megad2561.NewController(options); if err != nil {
		panic(err)
	}

	controller.MegadID = "new"

	err = controller.ApplyConfigsChanges(); if err != nil {
		panic(err)
	}
}
```

### Изменение комманды порта типа IN
```go
package main

import (
	"github.com/echokepler/megad2561"
)

func main() {
	options := megad2561.OptionsController{
		Host: "192.168.88.14",
		Password: "sec",
	}
	controller, err := megad2561.NewController(options); if err != nil {
		panic(err)
	}

	port, err := controller.Ports.GetByID(4, megad2561.InputType); if err != nil {
		panic(err)
	}

	portInput := port.(*megad2561.PortInput)
	portInput.Commands = "22:2|g0:0;g1:0;22:0"

	err = controller.Ports.Set(port); if err != nil {
		panic(err)
	}
}
```


## Задачи

**Порты:**
- [x] Input
- [ ] Output
- [ ] АЦП
- [ ] Dsen

**Конфиги:**
- [x] Основной конфиг
- [x] MegadID
- [ ] Cron
- [ ] Program
- [ ] СКУД


**Остальное:**
- [ ] Оповещать о неуспешной авторизации
