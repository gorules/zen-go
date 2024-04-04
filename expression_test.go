package zen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluateExpression(t *testing.T) {
	type TestCase[T any] struct {
		expression string
		output     T
		context    any
	}

	// Example usage with int
	intTestCases := []TestCase[int]{
		{expression: "1 + 1", output: 2},
		{expression: "2 + 2", output: 4},
		{expression: "10 + a", output: 14, context: map[string]int{"a": 4}},
	}

	// Example usage with string
	stringTestCases := []TestCase[string]{
		{expression: `"hello" + " " + "world"`, output: "hello world"},
		{expression: `"foo" + "bar"`, output: "foobar"},
	}

	for _, intTestCase := range intTestCases {
		res, err := EvaluateExpression[int](intTestCase.expression, intTestCase.context)
		assert.NoError(t, err)
		assert.Equal(t, intTestCase.output, res)
	}

	for _, stringTestCase := range stringTestCases {
		res, err := EvaluateExpression[string](stringTestCase.expression, stringTestCase.context)
		assert.NoError(t, err)
		assert.Equal(t, stringTestCase.output, res)
	}
}

func TestEvaluateUnaryExpression(t *testing.T) {
	type TestCase struct {
		expression string
		output     bool
		context    any
	}

	testCases := []TestCase{
		{
			expression: "> 10",
			output:     false,
			context:    map[string]any{"$": 5},
		},
		{
			expression: "> 10",
			output:     true,
			context:    map[string]any{"$": 15},
		},
		{
			expression: "'US', 'GB'",
			output:     true,
			context:    map[string]any{"$": "US"},
		},
		{
			expression: "'US', 'GB'",
			output:     false,
			context:    map[string]any{"$": "AA"},
		},
	}

	for _, testCase := range testCases {
		isTrue, err := EvaluateUnaryExpression(testCase.expression, testCase.context)
		assert.NoError(t, err)
		assert.Equal(t, testCase.output, isTrue)
	}
}

func TestRenderTemplate(t *testing.T) {
	type TestCase[T any] struct {
		template string
		output   any
		context  any
	}

	intTestCases := []TestCase[int]{
		{
			template: "{{ a + b }}",
			output:   15,
			context:  map[string]any{"a": 5, "b": 10},
		},
	}

	stringTestCases := []TestCase[string]{
		{
			template: "Hello: {{ a + b }}",
			output:   "Hello: 15",
			context:  map[string]any{"a": 5, "b": 10},
		},
	}

	for _, testCase := range intTestCases {
		isTrue, err := RenderTemplate[int](testCase.template, testCase.context)
		assert.NoError(t, err)
		assert.Equal(t, testCase.output, isTrue)
	}

	for _, testCase := range stringTestCases {
		isTrue, err := RenderTemplate[string](testCase.template, testCase.context)
		assert.NoError(t, err)
		assert.Equal(t, testCase.output, isTrue)
	}
}
