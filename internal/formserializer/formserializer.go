package formserializer

import (
	"github.com/PuerkitoBio/goquery"
)

type formElementType uint8

const (
	checkbox formElementType = iota
	text
	sel
)

type FormElement struct {
	Name  string
	Type  formElementType
	Value string
}

type Field struct {
	Name string
}

type Checkbox struct {
	Field
	Value bool
}

type InputText struct {
	Field
	Value string
}

type Select struct {
	Field
	Value []string
}

type FormSerializer struct {
	Checkboxes []Checkbox
	Texts      []InputText
	Selects    []Select
}

func Serialize(target *goquery.Selection) FormSerializer {
	fs := FormSerializer{}

	target.Find("input").Each(func(i int, s *goquery.Selection) {
		t, ok := s.Attr("type")

		if ok {
			switch t {
			case "checkbox":
				fs.Checkboxes = append(fs.Checkboxes, parseCheckbox(s))
			case "hidden", "text":
				fs.Texts = append(fs.Texts, parseInputText(s))
			}
		} else {
			fs.Texts = append(fs.Texts, parseInputText(s))
		}
	})

	target.ChildrenFiltered("select").Each(func(i int, s *goquery.Selection) {
		fs.Selects = append(fs.Selects, parseSelect(s))
	})

	return fs
}

/**
* parseCheckbox парсит DOM и возвращает объект Checkbox
**/
func parseCheckbox(node *goquery.Selection) Checkbox {
	checkbox := Checkbox{}

	name, ok := node.Attr("name")
	if ok {
		checkbox.Name = name
	}

	checked, ok := node.Attr("checked")
	if ok && checked == "" {
		checkbox.Value = true
	} else {
		checkbox.Value = false
	}

	return checkbox
}

/**
* parseInputText парсит DOM и возвращает объект InputText
**/
func parseInputText(node *goquery.Selection) InputText {
	input := InputText{}

	name, ok := node.Attr("name")
	if ok {
		input.Name = name
	}

	value, ok := node.Attr("value")
	if ok {
		input.Value = value
	}

	return input
}

/**
* parseSelect парсит DOM и возвращает объект Select
**/
func parseSelect(node *goquery.Selection) Select {
	selectInput := Select{}

	name, ok := node.Attr("name")
	if ok {
		selectInput.Name = name
	}

	node.ChildrenFiltered("option").Each(func(i int, s *goquery.Selection) {
		_, ok := s.Attr("selected")
		if ok {
			value, ok := s.Attr("value")
			if ok {
				selectInput.Value = append(selectInput.Value, value)
			}
		}
	})

	return selectInput
}
