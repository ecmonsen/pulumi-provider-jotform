package main

import (
	"encoding/json"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"strconv"
)

// Jotform API docs: https://api.jotform.com/docs
// Jotform property reference: https://api.jotform.com/docs/properties/index.php

// Each resource has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// The methods are defined in https://github.com/pulumi/pulumi-go-provider/blob/main/infer/resource.go.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resource are compared.
// - Update: Mutate a resource in place.
// - Read: Get the state of a resource from the backing provider.
// - Delete: Custom logic when the resource is deleted.
// - Annotate: Describe fields and set defaults for a resource.
// - WireDependencies: Control how outputs and secrets flows through values.
type Form struct {
}

type EmailArgs struct {
	Type    string `pulumi:"type" json:"type"`
	From    string `pulumi:"from" json:"from"`
	To      string `pulumi:"to" json:"to"`
	Subject string `pulumi:"subject" json:"subject"`
	Html    bool   `pulumi:"html" json:"html"`
	Body    string `pulumi:"body" json:"body"`
}

type FormArgs struct {
	Title  string      `pulumi:"title" json:"title"`
	Emails []EmailArgs `pulumi:"emails" json:"emails"`
}

type FormState struct {
	FormArgs
	FormId    string `pulumi:"form_id" json:"id"`
	Username  string `pulumi:"username" json:"username"`
	Title     string `pulumi:"title" json:"title"`
	Height    string `pulumi:"height" json:"height"`
	Status    string `pulumi:"status" json:"status"`
	CreatedAt string `pulumi:"created_at" json:"created_at"`
	UpdatedAt string `pulumi:"updated_at" json:"updated_at"`
	// Uninitialized "any" causes a panic later on
	//New       any    `pulumi:"new" json:"new"`
	Url string `pulumi:"url" json:"url"`
}

func (f Form) Create(ctx p.Context, name string, input FormArgs, preview bool) (string, FormState, error) {
	// from https://github.com/pulumi/pulumi-go-provider/blob/d8f8412f9990c708b0a3dcba792bacf72bf2788b/tests/grpc/config/provider/provider.go#L26
	config := infer.GetConfig[Config](ctx)
	JotformClient.ApiKey = config.ApiKey
	state := input.initState()
	if preview {
		return name, state, nil
	}
	state, err := f.create(input)
	return name, state, err
}

func (f Form) Delete(ctx p.Context, id string, props FormState) error {
	config := infer.GetConfig[Config](ctx)
	JotformClient.ApiKey = config.ApiKey
	err := f.delete(props)
	return err
}

func (f Form) Read(ctx p.Context, id string, inputs FormArgs, state FormState) (
	canonicalID string, normalizedInputs FormArgs, normalizedState FormState, err error) {
	config := infer.GetConfig[Config](ctx)
	JotformClient.ApiKey = config.ApiKey
	normalizedInputs, normalizedState, err = f.read(inputs, state)
	return id, normalizedInputs, normalizedState, err
}

func (fa FormArgs) initState() FormState {
	emails := fa.Emails
	if emails == nil {
		//
		emails = make([]EmailArgs, 0, 1)
	}

	state := FormState{FormArgs: FormArgs{Title: fa.Title, Emails: emails}}
	return state
}

func (fa FormArgs) AsMap() map[string]any {
	formData := map[string]any{
		"properties": map[string]string{
			"title": fa.Title,
		},
	}
	//if fa.Questions != nil {
	//      formData["questions"] = questionArgsAsMap(fa.Questions)
	//}
	if fa.Emails != nil {
		emails := make(map[string]any, len(fa.Emails))
		for index, email := range fa.Emails {
			var htmlString string = "0"
			if email.Html {
				htmlString = "1"
			}
			emails[strconv.Itoa(index)] = map[string]string{
				"type":    email.Type,
				"from":    email.From,
				"to":      email.To,
				"subject": email.Subject,
				"html":    htmlString,
				"body":    email.Body,
			}
		}
		formData["emails"] = emails
	}
	return formData
}

func (f Form) create(input FormArgs) (FormState, error) {
	state := input.initState()
	formData := input.AsMap()
	formResultBytes, createFormErr := JotformClient.CreateForm(formData)
	if createFormErr != nil {
		return state, createFormErr
	}
	err := json.Unmarshal(formResultBytes, &state)
	if err != nil {
		return state, err
	}
	return state, nil
}

func (f Form) delete(state FormState) error {
	_, err := JotformClient.DeleteForm(state.FormId)
	return err
}

func (f Form) read(args FormArgs, state FormState) (FormArgs, FormState, error) {
	formResultBytes, getFormErr := JotformClient.GetForm(state.FormId)
	if getFormErr != nil {
		return args, state, getFormErr
	}
	err := json.Unmarshal(formResultBytes, &state)
	if err != nil {
		return args, state, err
	}
	// TODO emails
	return args, state, nil
}
