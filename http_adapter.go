package megad2561

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/echokepler/megad2561/internal/formserializer"
	"net/http"
	"net/url"
	"strconv"
)

type HTTPAdapter struct {
	Host string
}

/**
* Get - запрашивает данные
**/
func (adapter *HTTPAdapter) Get(params ServiceValues) (ServiceValues, error) {
	values := ServiceValues{}

	queries := adapter.convertToQueries(params)
	uri := adapter.makeURL(queries)

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	form := doc.Find("form")

	fs := formserializer.Serialize(form)

	for _, checkbox := range fs.Checkboxes {
		values.Add(checkbox.Name, strconv.FormatBool(checkbox.Value))
	}

	for _, input := range fs.Texts {
		values.Add(input.Name, input.Value)
	}

	for _, selectInput := range fs.Selects {
		for _, selectedValues := range selectInput.Value {
			values.Set(selectInput.Name, selectedValues)
		}
	}

	return values, nil
}

/**
* Post - отправляет значения в сервис
**/
func (adapter *HTTPAdapter) Post(values ServiceValues) error {
	formattedValues, err := adapter.formatValues(values)
	if err != nil {
		return err
	}

	uri := adapter.makeURL(url.Values(formattedValues))

	fmt.Println(uri)

	res, err := http.Get(uri)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

/**
* convertToQueries преобразует ServiceValues в url.Values
**/
func (adapter *HTTPAdapter) convertToQueries(values ServiceValues) url.Values {
	queries := make(url.Values)

	// Перекладываем из ParamsValues в url.Values из-за риска смены дочернего типа
	for k, params := range values {
		for _, param := range params {
			queries.Add(k, param)
		}
	}

	return queries
}

/**
* makeURL Собирает итоговый URL
**/
func (adapter *HTTPAdapter) makeURL(queries url.Values) string {
	return fmt.Sprintf("http://%v/?%v", adapter.Host, queries.Encode())
}

func (adapter *HTTPAdapter) formatValues(values ServiceValues) (ServiceValues, error) {
	for key := range values {
		if values.IsBool(key) {
			boolValue, err := strconv.ParseBool(values.Get(key))
			if err != nil {
				return nil, err
			}
			if boolValue {
				values.Add(key, "1")
			} else {
				values.Del(key)
			}
		}
	}

	return values, nil
}
