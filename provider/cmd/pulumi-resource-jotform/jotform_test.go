package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

var testTime time.Time
var apiKey string

func TestMain(m *testing.M) {
	// Set up for tests
	testTime = time.Now()
	apiKey = os.Getenv("JOTFORM_API_KEY")
	JotformClient.ApiKey = apiKey
	JotformClient.SetOutputType("json")
	//JotformClient.SetDebugMode(true)
	exitVal := m.Run()
	// Clean up here if needed
	os.Exit(exitVal)
}

func Test_questionAsMap(t *testing.T) {
	fqq := FormQuestionsQuestion{Name: "qname", Type: "test_type", Order: 2, Text: "hello",
		Properties: map[string]string{"prop1": "val1", "prop2": "val2", "type": "this_should_be_ignored"}}
	actual := fqq.AsMap()
	expected := map[string]any{
		"name": "qname", "type": "test_type", "order": int64(2), "text": "hello", "prop1": "val1", "prop2": "val2"}
	if !reflect.DeepEqual(actual, expected) {
		t.Logf("Maps are not equal:\nexpected %#v\nreceived %#v", expected, actual)
		t.FailNow()
	}

}

func Test_jsonUnmarshalQuestion(t *testing.T) {
	b := []byte(`{"name":"qname","type":"test_type","order":2,"text":"test_text", "prop1": "val1", "prop2": "val2"}`)
	var actual FormQuestionsQuestion
	err := json.Unmarshal(b, &actual)
	if err != nil {
		t.Errorf("Error during parsing %#v", err)
		t.FailNow()
	}
	expected := FormQuestionsQuestion{
		Id: 0, Name: "qname", Type: "test_type", Order: 2, Text: "test_text", Properties: map[string]string{"prop1": "val1", "prop2": "val2"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Logf("Maps are not equal:\nexpected %#v\nreceived %#v", expected, actual)
		t.FailNow()
	}
}

func Test_jsonUnmarshalSlice(t *testing.T) {
	b := []byte(`[{"name":"qname","type":"test_type","order":2,"text":"test_text", "prop1": "val1", "prop2": "val2"},{"name":"qname2","type":"test_type2","order":3,"text":"test_text2"}]`)

	var expected []FormQuestionsQuestion
	expected = append(expected, FormQuestionsQuestion{
		Name: "qname", Type: "test_type", Order: 2, Text: "test_text", Properties: map[string]string{"prop1": "val1", "prop2": "val2"},
	})
	expected = append(expected, FormQuestionsQuestion{Name: "qname2", Type: "test_type2", Order: 3, Text: "test_text2", Properties: map[string]string{}})

	actual, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestion](b, questionsAsSlice)
	if err != nil {
		t.Errorf("Error during parsing %#v", err)
		t.FailNow()
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Logf("Maps are not equal:\nexpected %#v\nreceived %#v", expected, actual)
		t.FailNow()
	}

}

func Test_jsonUnmarshalMap(t *testing.T) {
	b := []byte(`{"1": {"name":"qname","type":"test_type","order":2,"text":"test_text", "prop1": "val1", "prop2": "val2"},"2":{"name":"qname2","type":"test_type2","order":3,"text":"test_text2"}}`)

	var expected []FormQuestionsQuestion
	expected = append(expected, FormQuestionsQuestion{
		Name: "qname", Type: "test_type", Order: 2, Text: "test_text", Properties: map[string]string{"prop1": "val1", "prop2": "val2"},
	})
	expected = append(expected, FormQuestionsQuestion{Name: "qname2", Type: "test_type2", Order: 3, Text: "test_text2", Properties: map[string]string{}})

	actual, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestion](b, questionsAsSlice)
	if err != nil {
		t.Errorf("Error during parsing %#v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Logf("Maps are not equal:\nexpected %#v\nreceived %#v", expected, actual)
		t.FailNow()
	}

}

func Test_jsonUnmarshalSlice2(t *testing.T) {
	b := []byte(`[{"name":"yourname","order":1,"qid":10,"text":"enter your name","type":"control_textbox"},{"name":"yourheight","order":2,"qid":20,"text":"enter your height","type":"control_textbox"}]`)

	var expected []FormQuestionsQuestion
	expected = append(expected, FormQuestionsQuestion{
		Name: "yourname", Type: "control_textbox", Order: 1, Text: "enter your name", Id: 10, Properties: map[string]string{},
	})
	expected = append(expected, FormQuestionsQuestion{
		Name: "yourheight", Type: "control_textbox", Order: 2, Text: "enter your height", Id: 20, Properties: map[string]string{}})

	actual, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestion](b, questionsAsSlice)
	if err != nil {
		t.Errorf("Error during parsing %#v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Logf("Maps are not equal:\nexpected %#v\nreceived %#v", expected, actual)
		t.FailNow()
	}

}

func Test_createReadAndDelete(t *testing.T) {
	f := Form{}
	formTitle := fmt.Sprintf("unit test form %s", testTime)
	input := FormArgs{formTitle, nil}
	state, err := f.create(input)
	if err != nil {
		t.Errorf("Error creating form: %s", err)
		return
	}
	t.Logf("Created form #{state.FormId}")
	defer func() {
		t.Logf("Deleting form %s", state.FormId)
		deleteErr := f.delete(state)
		if deleteErr != nil {
			t.Errorf("Error deleting form: %s", deleteErr)
		}
	}()
	//readArgs, readState,
	readFormArgs, readFormState, readFormErr := f.read(input, state)
	if readFormErr != nil {
		t.Errorf("Error reading form: %s", readFormErr)
	}
	if readFormArgs.Title != formTitle {
		t.Logf("Expected form title '%s', got '%s'", formTitle, readFormArgs.Title)
		t.FailNow()
	}
	if !reflect.DeepEqual(readFormArgs, input) {
		t.Logf("Unexpected form args")
		t.Logf("%#v", readFormArgs)
		t.FailNow()
	}
	if readFormState.FormId == "" {
		t.Logf("Unexpectedly empty form ID")
		t.Logf("%#v", readFormState)
		t.FailNow()
	}

	// Questions
	fq := FormQuestions{}

	questions := []FormQuestionsQuestion{
		{Type: "control_textbox", Text: "enter your name", Order: 1, Name: "yourname", Properties: map[string]string{}},
		{Type: "control_textbox", Text: "enter your height", Order: 2, Name: "yourheight", Properties: map[string]string{}},
	}
	fqInput := FormQuestionsArgs{FormId: state.FormId, Questions: questions}
	fqState, createQErr := fq.create(fqInput)
	if createQErr != nil {
		t.Errorf("Error creating form questions: %s", createQErr)
	}

	readQuestionsArgs, readQuestionsState, readQuestionsErr := fq.read(fqInput, fqState)

	// We cannot predict the IDs assigned by Jotform
	for i, _ := range readQuestionsArgs.Questions {
		readQuestionsArgs.Questions[i].Id = 0
	}
	if readQuestionsErr != nil {
		t.Errorf("Error reading form questions: %s", readQuestionsErr)
	}
	if !reflect.DeepEqual(readQuestionsArgs, fqInput) {
		t.Logf("Unexpected read questions args")
		t.Logf("expected: %#v", fqInput)
		t.Logf("actual: %#v", readQuestionsArgs)
		t.Fail()
	}
	expectedState := FormQuestionsState{FormQuestionsArgs: fqInput}
	if !reflect.DeepEqual(readQuestionsState, expectedState) {
		t.Logf("Unexpected questions state")
		t.Logf("expected: %#v", expectedState)
		t.Logf("actual: %#v", readQuestionsState)
		t.Fail()
	}
}
