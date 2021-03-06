package core

import "net/url"

type TypeServiceError int

const (
	UnAuthorized TypeServiceError = iota
)

// ServiceAdapter общий интерфейс сервис адаптера,
// далее по нему могут быть реализованы нативный и кастомные адаптеры
type ServiceAdapter interface {
	Get(params ServiceValues) (ServiceValues, error)
	Post(values ServiceValues) error
}

type ErrorService struct {
	Type TypeServiceError
	Err  error
}

func (ce *ErrorService) Error() string {
	return ce.Err.Error()
}

// По сути ServiceValues нужен в качестве подстраховки на случай смены структуры этих типов.
// Использование url.Values является более подходящим,
// т.к все общение с контроллером далее происходит по query параметрам
type (
	// ServiceValues нужен для возвращения данных из сервиса
	ServiceValues url.Values
)

func (sv ServiceValues) Set(key string, value string) {
	sv[key] = append(sv[key], value)
}

func (sv ServiceValues) Add(key string, value string) {
	sv[key] = []string{value}
}

func (sv ServiceValues) Get(key string) string {
	if sv[key] == nil {
		return ""
	} else if len(sv[key]) > 0 {
		return sv[key][0]
	}

	return sv[key][0]
}

func (sv ServiceValues) Has(key string) bool {
	return sv[key] != nil
}

func (sv ServiceValues) Del(key string) {
	delete(sv, key)
}

func (sv ServiceValues) Encode() string {
	return url.Values(sv).Encode()
}

func (sv ServiceValues) Parse(str string) (ServiceValues, error) {
	values, err := url.ParseQuery(str)
	if err != nil {
		return nil, err
	}

	return ServiceValues(values), nil
}

func (sv ServiceValues) IsEmpty() bool {
	return len(sv) == 0
}

func (sv ServiceValues) IsBool(key string) bool {
	switch sv.Get(key) {
	case "true", "false":
		return true
	default:
		return false
	}
}
