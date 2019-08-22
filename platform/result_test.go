package platform_test

import (
	"errors"
	"testing"

	"github.com/HoaHuynhSoft/go-core/platform"
)

func TestValidationErrorResult(t *testing.T) {
	tt := []struct {
		Result          *platform.Result
		expectCode      string
		expectedMessage string
	}{
		{platform.ValidationErrorResult(errors.New("Jed_Test")), "UNKNOWN_ERROR", "Jed_Test"},
		{platform.ValidationErrorResult(errors.New("Jed_Test"), "CODE_1"), "CODE_1", "Jed_Test"},
		{platform.ValidationErrorResult(errors.New("Jed_Test"), "CODE_1", "MESSAGE_1"), "CODE_1", "MESSAGE_1"},

		{platform.NotFoundErrorResult(errors.New("Jed_Test")), "UNKNOWN_ERROR", "Jed_Test"},
		{platform.NotFoundErrorResult(errors.New("Jed_Test"), "CODE_1"), "CODE_1", "Jed_Test"},
		{platform.NotFoundErrorResult(errors.New("Jed_Test"), "CODE_1", "MESSAGE_1"), "CODE_1", "MESSAGE_1"},

		{platform.InternalErrorResult(errors.New("Jed_Test")), "UNKNOWN_ERROR", "Jed_Test"},
		{platform.InternalErrorResult(errors.New("Jed_Test"), "CODE_1"), "CODE_1", "Jed_Test"},
		{platform.InternalErrorResult(errors.New("Jed_Test"), "CODE_1", "MESSAGE_1"), "CODE_1", "MESSAGE_1"},
	}

	for _, tc := range tt {
		if tc.Result.IsSuccess {
			t.Errorf("Expected fail result but got %v", tc.Result.IsSuccess)
		}
		errCode := tc.Result.Error.GetCode()
		if errCode != tc.expectCode {
			t.Errorf("Defalt Error Code must be %v but got %v", tc.expectCode, errCode)
		}
		errMessage := tc.Result.Error.GetErrorMessage()
		if errMessage != tc.expectedMessage {
			t.Errorf("Expected Error Message be %v but got %v", tc.expectedMessage, errMessage)
		}
	}
}

func TestSuccessResult(t *testing.T) {
	tt := []interface{}{
		struct {
			testField string
		}{
			testField: "TEST",
		},
		struct {
			testField1 string
		}{
			testField1: "TEST1",
		},
		struct {
			testField2 string
		}{
			testField2: "TEST2",
		},
	}

	for _, tc := range tt {
		result := platform.OKResult(tc)
		if result.Value != tc {
			t.Errorf("value of the result must be the object have been set %v but got %v", tc, result.Value)
		}
		if !result.IsSuccess {
			t.Errorf("expected success result but got %v", result.IsSuccess)
		}
	}
}
