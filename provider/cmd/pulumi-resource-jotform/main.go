// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "github.com/jotform/jotform-api-go"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Version is initialized by the Go linker to contain the semver of this build.

func main() {
	// To do add metadata, example at
	// https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go
	p.RunProvider("jotform", "0.0.1",
		// We tell the provider what resources it needs to support.

		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource[Form, FormArgs, FormState](),
				infer.Resource[FormQuestions, FormQuestionsArgs, FormQuestionsState](),
			},
			// This doesn't really do what I hoped
			Config: infer.Config[*Config](),
		}))
}

type Config struct {
	ApiKey  string `pulumi:"api_key,optional" provider:"secret"`
	Version string `pulumi:"version,optional"`
}
