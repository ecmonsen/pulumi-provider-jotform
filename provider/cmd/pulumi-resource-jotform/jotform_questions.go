package main

import (
	"encoding/json"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"sort"
	"strconv"
)

type FormQuestions struct {
}

type FormQuestionsQuestionArgs struct {
	Type  string `pulumi:"type" json:"type"`
	Text  string `pulumi:"text" json:"text"`
	Order string `pulumi:"order" json:"order"`
	Name  string `pulumi:"name" json:"name"`
	Id    string `pulumi:"id" json:"qid"`
}

type FormQuestionsArgs struct {
	FormId    string                      `pulumi:"form_id" json:"form_id,optional"`
	Questions []FormQuestionsQuestionArgs `pulumi:"questions" json:"questions"`
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

func (fqa FormQuestionsArgs) initState() FormQuestionsState {
	questions := fqa.Questions
	if questions == nil {
		//
		questions = make([]FormQuestionsQuestionArgs, 0, 1)
	}

	state := FormQuestionsState{FormQuestionsArgs: FormQuestionsArgs{FormId: fqa.FormId, Questions: questions}}
	return state
}

type altFormQuestionsArgs struct {
	FormId    string                               `pulumi:"form_id" json:"form_id"`
	Questions map[string]FormQuestionsQuestionArgs `pulumi:"questions" json:"questions"`
}

func (a altFormQuestionsArgs) Convert() (FormQuestionsArgs, error) {
	formQuestionsArgs := FormQuestionsArgs{FormId: a.FormId}
	questions, err := questionsAsSlice(a.Questions)
	formQuestionsArgs.Questions = questions
	return formQuestionsArgs, err
}

// Sometimes Jotform API uses array of questions, sometimes an array of qid (which must be an integer) and question.
func (fqa FormQuestionsArgs) questionsAsMap() map[string]FormQuestionsQuestionArgs {
	questionsMap := make(map[string]FormQuestionsQuestionArgs)
	for i, v := range fqa.Questions {
		questionsMap[strconv.Itoa(i)] = v
	}
	return questionsMap
}

func questionsAsSlice(questionsMap map[string]FormQuestionsQuestionArgs) ([]FormQuestionsQuestionArgs, error) {
	questionsSlice := make([]FormQuestionsQuestionArgs, 0, len(questionsMap))
	orderArray := make([]int64, 0, len(questionsMap))
	for _, v := range questionsMap {
		order, orderErr := strconv.ParseInt(v.Order, 10, 64)
		if orderErr != nil {
			return questionsSlice, orderErr
		}
		orderArray = append(orderArray, order)
		questionsSlice = append(questionsSlice, v)
	}
	sort.Slice(questionsSlice, func(i, j int) bool {
		iOrder, _ := strconv.ParseInt(questionsSlice[i].Order, 10, 64)
		jOrder, _ := strconv.ParseInt(questionsSlice[j].Order, 10, 64)
		return iOrder < jOrder
	})
	return questionsSlice, nil
}

func (question FormQuestionsQuestionArgs) AsMap() map[string]string {
	return map[string]string{
		"type":  question.Type,
		"text":  question.Text,
		"order": question.Order,
		"name":  question.Name,
	}
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
	questions, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestionArgs](questionBytes, questionsAsSlice)
	state.FormId = args.FormId
	state.Questions = questions
	return state, err
}

func jsonUnmarshalSliceOrMap[T any](jsonBytes []byte, convert func(map[string]T) ([]T, error)) ([]T, error) {
	var ts = make([]T, 0, 10)
	err := json.Unmarshal(jsonBytes, &ts)
	if err != nil {
		var tm = make(map[string]T)
		err = json.Unmarshal(jsonBytes, &tm)
		if err != nil {
			return ts, err
		}

		ts, err = convert(tm)
	}
	return ts, err
}

func (fq FormQuestions) read(args FormQuestionsArgs, state FormQuestionsState) (FormQuestionsArgs, FormQuestionsState, error) {
	questionsBytes, getQuestionsErr := JotformClient.GetFormQuestions(args.FormId)
	if getQuestionsErr != nil {
		return args, state, getQuestionsErr
	}
	questions, err := jsonUnmarshalSliceOrMap[FormQuestionsQuestionArgs](questionsBytes, questionsAsSlice)
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
		_, err := JotformClient.DeleteFormQuestion(state.FormId, q.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
