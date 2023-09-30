package main

import (
	"encoding/json"
	"errors"
	"fmt"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"reflect"
	"sort"
	"strconv"
)

type FormQuestions struct {
}

type FormQuestionsQuestion struct {
	Type string `pulumi:"type" json:"type"`
	Text string `pulumi:"text" json:"text"`
	// Undocumented Jotform requirement: Order must be int or stringified int
	Order int64  `pulumi:"order" json:"order"`
	Name  string `pulumi:"name" json:"name"`
	// Undocumented Jotform requirement: qid must be a number or stringified number
	Id int64 `pulumi:"qid,optional" json:"qid,optional"`
	// Additional properties depending on Type. See https://api.jotform.com/docs/properties/index.php
	Properties map[string]string `pulumi:"properties,optional"`
}

type FormQuestionsArgs struct {
	FormId    string                  `pulumi:"form_id" json:"form_id,optional"`
	Questions []FormQuestionsQuestion `pulumi:"questions" json:"questions"`
}

type FormQuestionsState struct {
	FormQuestionsArgs
}

func (fq FormQuestions) Create(ctx p.Context, name string, input FormQuestionsArgs, preview bool) (string, FormQuestionsState, error) {
	config := infer.GetConfig[Config](ctx)
	JotformClient.ApiKey = config.ApiKey
	JotformClient.SetDebugMode(true)
	state := input.initState()
	if preview {
		return name, state, nil
	}
	state, err := fq.create(input)
	return name, state, err
}

func (fq FormQuestions) Delete(ctx p.Context, id string, props FormQuestionsState) error {
	config := infer.GetConfig[Config](ctx)
	JotformClient.ApiKey = config.ApiKey
	err := fq.delete(props)
	return err
}

func (fq FormQuestions) Read(ctx p.Context, id string, inputs FormQuestionsArgs, state FormQuestionsState) (
	canonicalID string, normalizedInputs FormQuestionsArgs, normalizedState FormQuestionsState, err error) {
	config := infer.GetConfig[Config](ctx)
	JotformClient.ApiKey = config.ApiKey
	normalizedInputs, normalizedState, err = fq.read(inputs, state)
	return id, normalizedInputs, normalizedState, err
}

func convertToInt64(v any) (int64, error) {
	var i int64
	var err error
	switch v2 := v.(type) {
	case string:
		i, err = strconv.ParseInt(v2, 10, 64)
		if err != nil {
			return 0, err
		}
	case float64:
		i = int64(v2)
	case nil:
		i = 0
	default:
		return i, errors.New(fmt.Sprintf("Not converting type %s", reflect.TypeOf(v)))
	}
	return i, nil

}
func (fqq *FormQuestionsQuestion) UnmarshalJSON(b []byte) error {
	type FormQuestionsQuestion1 struct {
		Type string `pulumi:"type" json:"type"`
		Text string `pulumi:"text" json:"text"`
		// Undocumented Jotform requirement: Order can be int or string
		Order interface{} `pulumi:"order" json:"order"`
		Name  string      `pulumi:"name" json:"name"`
		// We will convert this to int
		Id interface{} `pulumi:"qid" json:"qid"`
		// Additional properties depending on Type. See https://api.jotform.com/docs/properties/index.php
		Properties map[string]interface{} `pulumi:"properties"`
	}
	var fqq1 FormQuestionsQuestion1
	err := json.Unmarshal(b, &fqq1)
	if err != nil {
		return err
	}
	err2 := json.Unmarshal(b, &fqq1.Properties)
	if err2 != nil {
		return err2
	}

	t := reflect.TypeOf(fqq1)

	for i := 0; i < t.NumField(); i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag != "" {
			delete(fqq1.Properties, jsonTag)
		}
	}

	*fqq = FormQuestionsQuestion{
		Name:       fqq1.Name,
		Type:       fqq1.Type,
		Text:       fqq1.Text,
		Properties: map[string]string{},
	}
	fqq.Order, err = convertToInt64(fqq1.Order)
	if err != nil {
		return err
	}
	fqq.Id, err = convertToInt64(fqq1.Id)
	if err != nil {
		return err
	}

	// In attempt at simplicity convert all other Properties to strings
	for k, v := range fqq1.Properties {
		switch v3 := v.(type) {
		case string:
			fqq.Properties[k] = v3
		case int64:
			fqq.Properties[k] = strconv.FormatInt(v3, 10)
		case int:
			fqq.Properties[k] = strconv.Itoa(v3)
		default:
			return errors.New(fmt.Sprintf("Can't convert value of type %s", reflect.TypeOf(v3)))
		}
	}
	return nil
}

func (fqa FormQuestionsArgs) initState() FormQuestionsState {
	questions := fqa.Questions
	if questions == nil {
		//
		questions = make([]FormQuestionsQuestion, 0, 1)
	}

	state := FormQuestionsState{FormQuestionsArgs: FormQuestionsArgs{FormId: fqa.FormId, Questions: questions}}
	return state
}

type altFormQuestionsArgs struct {
	FormId    string                           `pulumi:"form_id" json:"form_id"`
	Questions map[string]FormQuestionsQuestion `pulumi:"questions" json:"questions"`
}

func (a altFormQuestionsArgs) Convert() (FormQuestionsArgs, error) {
	formQuestionsArgs := FormQuestionsArgs{FormId: a.FormId}
	questions, err := questionsAsSlice(a.Questions)
	formQuestionsArgs.Questions = questions
	return formQuestionsArgs, err
}

// Sometimes Jotform API uses array of questions, sometimes an array of qid (which must be an integer) and question.
func (fqa FormQuestionsArgs) questionsAsMap() map[string]any {
	questionsMap := make(map[string]any)
	for i, v := range fqa.Questions {
		questionsMap[strconv.Itoa(i)] = v.AsMap()
	}
	return questionsMap
}

func questionsAsSlice(questionsMap map[string]FormQuestionsQuestion) ([]FormQuestionsQuestion, error) {
	questionsSlice := make([]FormQuestionsQuestion, 0, len(questionsMap))
	for _, v := range questionsMap {
		questionsSlice = append(questionsSlice, v)
	}
	sort.Slice(questionsSlice, func(i, j int) bool {
		return questionsSlice[i].Order < questionsSlice[j].Order
	})
	return questionsSlice, nil
}

func (question FormQuestionsQuestion) AsMap() map[string]any {
	var m = map[string]any{
		"type":  question.Type,
		"text":  question.Text,
		"order": question.Order,
		"name":  question.Name,
	}
	for k, v := range question.Properties {
		_, exists := m[k]
		if !exists {
			m[k] = v
		} else {
			// Ignore keys such as 'type' that are already present in the map
			fmt.Println(fmt.Sprintf("Warning: property key '%s' ignored", k))
		}
	}
	return m
}

func (fq FormQuestions) create(args FormQuestionsArgs) (FormQuestionsState, error) {
	state := args.initState()
	questionsJson, jErr := json.Marshal(map[string]any{"questions": args.questionsAsMap()})
	if jErr != nil {
		return state, jErr
	}
	questionBytes, createErr := JotformClient.CreateFormQuestions(args.FormId, questionsJson)
	if createErr != nil {
		return state, createErr
	}
	questions, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestion](questionBytes, questionsAsSlice)
	state.FormId = args.FormId
	state.Questions = questions
	return state, err
}

func jsonUnmarshalSliceOrMap[T any](jsonBytes []byte, convertMapToSlice func(map[string]T) ([]T, error)) ([]T, error) {
	var ts = make([]T, 0, 10)
	err := json.Unmarshal(jsonBytes, &ts)
	if err != nil {
		var tm = make(map[string]T)
		err = json.Unmarshal(jsonBytes, &tm)
		if err != nil {
			return ts, err
		}

		ts, err = convertMapToSlice(tm)
	}
	return ts, err
}

func (fq FormQuestions) read(args FormQuestionsArgs, state FormQuestionsState) (FormQuestionsArgs, FormQuestionsState, error) {
	questionsBytes, getQuestionsErr := JotformClient.GetFormQuestions(args.FormId)
	if getQuestionsErr != nil {
		return args, state, getQuestionsErr
	}
	questions, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestion](questionsBytes, questionsAsSlice)
	if err != nil {
		return args, state, err
	}
	args.Questions = questions
	state.FormQuestionsArgs = args

	return args, state, nil
}

func (fq FormQuestions) delete(state FormQuestionsState) error {
	// No API to delete all questions. Questions must be deleted one by one.
	for _, q := range state.Questions {
		_, err := JotformClient.DeleteFormQuestion(state.FormId, strconv.FormatInt(q.Id, 10))
		if err != nil {
			return err
		}
	}
	return nil
}
