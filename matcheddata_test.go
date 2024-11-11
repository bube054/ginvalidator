package ginvalidator

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMatchedData(t *testing.T) {
	type data struct {
		location string
		field    string
		value    string
	}

	tests := []struct {
		name                string
		ctx                 *gin.Context
		insertedData        []data
		expectedMatchedData MatchedData
		expectedErr         error
	}{
		{
			name:                "Nil ctx provided",
			ctx:                 nil,
			insertedData:        []data{{location: "body", field: "testField", value: "testValue"}},
			expectedMatchedData: MatchedData{"body": MatchedDataFieldValues{"testField": "testValue"}},
			expectedErr:         ErrNilCtxMatchedData,
		},
		{
			name:                "Extracted matched body data",
			ctx:                 createTestGinCtx(ginCtxReqOpts{}),
			insertedData:        []data{{location: "body", field: "testField", value: "testValue"}},
			expectedMatchedData: MatchedData{"body": MatchedDataFieldValues{"testField": "testValue"}},
			expectedErr:         nil,
		},
		{
			name:                "Extracted matched cookies data",
			ctx:                 createTestGinCtx(ginCtxReqOpts{}),
			insertedData:        []data{{location: "cookies", field: "testField", value: "testValue"}},
			expectedMatchedData: MatchedData{"cookies": MatchedDataFieldValues{"testField": "testValue"}},
			expectedErr:         nil,
		},
		{
			name:                "Extracted matched headers data",
			ctx:                 createTestGinCtx(ginCtxReqOpts{}),
			insertedData:        []data{{location: "headers", field: "testField", value: "testValue"}},
			expectedMatchedData: MatchedData{"headers": MatchedDataFieldValues{"testField": "testValue"}},
			expectedErr:         nil,
		},
		{
			name:                "Extracted matched params data",
			ctx:                 createTestGinCtx(ginCtxReqOpts{}),
			insertedData:        []data{{location: "params", field: "testField", value: "testValue"}},
			expectedMatchedData: MatchedData{"params": MatchedDataFieldValues{"testField": "testValue"}},
			expectedErr:         nil,
		},
		{
			name:                "Extracted matched query data",
			ctx:                 createTestGinCtx(ginCtxReqOpts{}),
			insertedData:        []data{{location: "query", field: "testField", value: "testValue"}},
			expectedMatchedData: MatchedData{"query": MatchedDataFieldValues{"testField": "testValue"}},
			expectedErr:         nil,
		},
		{
			name: "Multiple fields in body data",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "body", field: "field1", value: "value1"},
				{location: "body", field: "field2", value: "value2"},
				{location: "body", field: "field3", value: "value3"},
			},
			expectedMatchedData: MatchedData{
				"body": MatchedDataFieldValues{
					"field1": "value1",
					"field2": "value2",
					"field3": "value3",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Multiple data from different locations",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "body", field: "bodyField", value: "bodyValue"},
				{location: "cookies", field: "cookieField", value: "cookieValue"},
				{location: "headers", field: "headerField", value: "headerValue"},
				{location: "params", field: "paramField", value: "paramValue"},
				{location: "query", field: "queryField", value: "queryValue"},
			},
			expectedMatchedData: MatchedData{
				"body":    MatchedDataFieldValues{"bodyField": "bodyValue"},
				"cookies": MatchedDataFieldValues{"cookieField": "cookieValue"},
				"headers": MatchedDataFieldValues{"headerField": "headerValue"},
				"params":  MatchedDataFieldValues{"paramField": "paramValue"},
				"query":   MatchedDataFieldValues{"queryField": "queryValue"},
			},
			expectedErr: nil,
		},
		{
			name: "Multiple fields in cookies and headers",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "cookies", field: "cookie1", value: "cookieValue1"},
				{location: "cookies", field: "cookie2", value: "cookieValue2"},
				{location: "headers", field: "header1", value: "headerValue1"},
				{location: "headers", field: "header2", value: "headerValue2"},
			},
			expectedMatchedData: MatchedData{
				"cookies": MatchedDataFieldValues{
					"cookie1": "cookieValue1",
					"cookie2": "cookieValue2",
				},
				"headers": MatchedDataFieldValues{
					"header1": "headerValue1",
					"header2": "headerValue2",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Multiple params and query fields",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "params", field: "param1", value: "paramValue1"},
				{location: "params", field: "param2", value: "paramValue2"},
				{location: "query", field: "query1", value: "queryValue1"},
				{location: "query", field: "query2", value: "queryValue2"},
			},
			expectedMatchedData: MatchedData{
				"params": MatchedDataFieldValues{
					"param1": "paramValue1",
					"param2": "paramValue2",
				},
				"query": MatchedDataFieldValues{
					"query1": "queryValue1",
					"query2": "queryValue2",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Mixed body, query, and headers data",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "body", field: "bodyField1", value: "bodyValue1"},
				{location: "query", field: "queryField1", value: "queryValue1"},
				{location: "headers", field: "headerField1", value: "headerValue1"},
				{location: "body", field: "bodyField2", value: "bodyValue2"},
				{location: "query", field: "queryField2", value: "queryValue2"},
			},
			expectedMatchedData: MatchedData{
				"body": MatchedDataFieldValues{
					"bodyField1": "bodyValue1",
					"bodyField2": "bodyValue2",
				},
				"query": MatchedDataFieldValues{
					"queryField1": "queryValue1",
					"queryField2": "queryValue2",
				},
				"headers": MatchedDataFieldValues{
					"headerField1": "headerValue1",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Overridden fields in body data",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "body", field: "field1", value: "initialValue"},
				{location: "body", field: "field1", value: "overriddenValue"},
			},
			expectedMatchedData: MatchedData{
				"body": MatchedDataFieldValues{
					"field1": "overriddenValue",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Overridden cookies data",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "cookies", field: "cookieField", value: "cookieValue1"},
				{location: "cookies", field: "cookieField", value: "cookieValue2"},
			},
			expectedMatchedData: MatchedData{
				"cookies": MatchedDataFieldValues{
					"cookieField": "cookieValue2",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Overridden headers and body data",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "headers", field: "headerField", value: "headerValue1"},
				{location: "body", field: "headerField", value: "bodyValue1"},
				{location: "body", field: "headerField", value: "bodyValue2"},
			},
			expectedMatchedData: MatchedData{
				"headers": MatchedDataFieldValues{
					"headerField": "headerValue1",
				},
				"body": MatchedDataFieldValues{
					"headerField": "bodyValue2",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Overridden params data",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "params", field: "param1", value: "paramValue1"},
				{location: "params", field: "param1", value: "overriddenParamValue1"},
			},
			expectedMatchedData: MatchedData{
				"params": MatchedDataFieldValues{
					"param1": "overriddenParamValue1",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Mixed overridden data across multiple locations",
			ctx:  createTestGinCtx(ginCtxReqOpts{}),
			insertedData: []data{
				{location: "body", field: "field1", value: "bodyValue1"},
				{location: "query", field: "field1", value: "queryValue1"},
				{location: "body", field: "field1", value: "overriddenBodyValue"},
				{location: "query", field: "field1", value: "overriddenQueryValue"},
			},
			expectedMatchedData: MatchedData{
				"body": MatchedDataFieldValues{
					"field1": "overriddenBodyValue",
				},
				"query": MatchedDataFieldValues{
					"field1": "overriddenQueryValue",
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, d := range test.insertedData {
				saveMatchedDataToCtx(test.ctx, d.location, d.field, d.value)
			}

			actualMatchedData, actualErr := GetMatchedData(test.ctx)

			if actualErr != nil {
				if !errors.Is(actualErr, test.expectedErr) {
					t.Errorf("got %+v, want %+v", actualErr, test.expectedErr)
				}
			} else {
				if !reflect.DeepEqual(actualMatchedData, test.expectedMatchedData) {
					t.Errorf("got %+v, want %+v", actualMatchedData, test.expectedMatchedData)
				}
			}
		})
	}
}
