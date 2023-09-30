package main

import (
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

	questions := []FormQuestionsQuestionArgs{
		{Type: "control_textbox", Text: "enter your name", Order: "1", Name: "yourname", Id: "345"},
		{Type: "control_textbox", Text: "enter your height", Order: "2", Name: "yourheight", Id: "567"},
	}
	fqInput := FormQuestionsArgs{FormId: state.FormId, Questions: questions}
	fqState, createQErr := fq.create(fqInput)
	if createQErr != nil {
		t.Errorf("Error creating form questions: %s", createQErr)
	}

	readQuestionsArgs, readQuestionsState, readQuestionsErr := fq.read(fqInput, fqState)
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
