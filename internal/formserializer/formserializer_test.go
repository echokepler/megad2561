package formserializer_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/echokepler/megad2561/internal/formserializer"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSerializeForms(t *testing.T) {
	file, err := os.Open("./mock/forms.html")
	if err != nil {
		t.Error(err)
	}

	document, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
	}

	fs := formserializer.SerializeForms(document)

	t.Run("Should be serialized two forms", func(t *testing.T) {
		assert.Equal(t, 2, len(fs.Texts))

		assert.Equal(t, "val1", fs.Texts[0].Value)
		assert.Equal(t, "val2", fs.Texts[1].Value)
	})
}

func TestFormSerializer_Serialize(t *testing.T) {
	file, err := os.Open("./mock/form.html")
	if err != nil {
		t.Error(err)
	}

	document, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
	}

	form := document.Find("form")

	fs := formserializer.Serialize(form)

	t.Run("Should be return correct length form elements", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, 2, len(fs.Selects))
		assert.Equal(t, 3, len(fs.Texts))
	})

	t.Run("Should be return correct selects values", func(t *testing.T) {
		t.Parallel()

		expectedSelects := []formserializer.Select{
			{
				Field: formserializer.Field{
					Name: "first_select",
				},
				Value: []string{"first", "third"},
			},
			{
				Field: formserializer.Field{
					Name: "second_select",
				},
				Value: []string{"second"},
			},
		}

		for i, input := range expectedSelects {
			assert.Equal(t, input.Name, fs.Selects[i].Name)
			assert.EqualValues(t, input.Value, fs.Selects[i].Value)
		}
	})

	t.Run("Should be return correct inputs values", func(t *testing.T) {
		expectedSelects := []formserializer.InputText{
			{
				Field: formserializer.Field{
					Name: "hidden_input",
				},
				Value: "hidden_value",
			},
			{
				Field: formserializer.Field{
					Name: "without_type",
				},
				Value: "value",
			},
			{
				Field: formserializer.Field{
					Name: "base_input",
				},
				Value: "value",
			},
		}

		for i, input := range expectedSelects {
			assert.Equal(t, input.Name, fs.Texts[i].Name)
			assert.EqualValues(t, input.Value, fs.Texts[i].Value)
		}
	})

	t.Run("Should correct return checkboxes", func(t *testing.T) {
		expectedValues := []formserializer.Checkbox{
			{
				Field: formserializer.Field{
					Name: "checked",
				},
				Value: true,
			},
			{
				Field: formserializer.Field{
					Name: "not_checked",
				},
				Value: false,
			},
		}

		for i, input := range expectedValues {
			assert.Equal(t, input.Name, fs.Checkboxes[i].Name)
			assert.Equal(t, input.Value, fs.Checkboxes[i].Value)
		}
	})

	if file.Close() != nil {
		t.Error(file.Close())
	}
}
