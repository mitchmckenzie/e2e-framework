/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/conf"
)

// EnvFunc represents a user-defined operation that
// can be used to customized the behavior of the
// environment.
type EnvFunc func(context.Context, *conf.Config) error

// Environment represents an environment where
// features can be tested.
type Environment interface {
	Config() *conf.Config
	Context() context.Context

	// Setup registers environment operations that are executed once
	// prior to the environment being ready and prior to any test.
	Setup(...EnvFunc) Environment

	// BeforeTest registers funcs that are executed before each Env.Test(...)
	BeforeTest(...EnvFunc) Environment

	// Test executes a test feature
	Test(context.Context, *testing.T, Feature)

	// AfterTest registers funcs that are executed after each Env.Test(...)
	AfterTest(...EnvFunc) Environment

	// Finish registers funcs that are executed at the end.
	Finish(...EnvFunc) Environment

	// Run Launches the test suite from within a TestMain
	Run(context.Context, *testing.M) int
}

type Labels map[string]string
type Feature interface {
	// Name is a descriptive text for the feature
	Name() string
	// Labels returns a map of feature labels
	Labels() Labels
	// Steps testing tasks to test the feature
	Steps() []Step
}

type Level uint8

const (
	LevelSetup Level = iota
	LevelAssess
	LevelTeardown
)

type StepFunc func (context.Context, *testing.T, *conf.Config)

type Step interface {
	// Name is the step name
	Name() string
	// Level action level {setup|requirement|assertion|teardown}
	Level() Level
	// Func is the operation for the step
	Func() StepFunc
}