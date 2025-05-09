package bdd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

var responseBody string

func givenAValidPaymentPayload() error {
	bodyData = map[string]any{
		"amount":             50.0,
		"external_reference": "valid-" + generateUUID(),
		"payment_method":     "pix",
	}
	return nil
}

func givenAPaymentPayloadWithExternalReference(ref string) error {
	bodyData = map[string]any{
		"amount":             "50.0",
		"external_reference": ref,
		"payment_method":     "pix",
	}
	return nil
}

func givenAnInvalidJSONPayload() error {
	request, _ = http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader([]byte("invalid-json")))
	request.Header.Set("Content-Type", "application/json")
	return nil
}

func givenAPaymentPayloadWithInvalidFields() error {
	bodyData = map[string]any{
		"amount":             "-10.0",
		"external_reference": "",
		"payment_method":     "",
	}
	return nil
}

func iSendAPostRequestToAuthorizePayment() error {
	if request == nil {
		payload, _ := json.Marshal(bodyData)
		request, _ = http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(payload))
		request.Header.Set("Content-Type", "application/json")
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	responseBody = recorder.Body.String()
	request = nil // reset
	return nil
}

func iSendAnotherPostRequestToAuthorizeWithSameExternalReference(ref string) error {
	body := map[string]interface{}{
		"amount":             99.0,
		"external_reference": ref,
		"payment_method":     "pix",
	}
	payload, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	responseBody = recorder.Body.String()
	return nil
}

func theResponseStatusShouldBe(expected int) error {
	if recorder.Code != expected {
		return fmt.Errorf("expected status %d but got %d", expected, recorder.Code)
	}
	return nil
}

func theResponseShouldContain(text string) error {
	if !strings.Contains(responseBody, text) {
		return fmt.Errorf("expected response to contain %q but got %q", text, responseBody)
	}
	return nil
}

func generateUUID() string {
	return uuid.New().String()
}

func InitializeScenarioPaymentAPI(ctx *godog.ScenarioContext) {
	ctx.Step(`^a valid payment payload$`, givenAValidPaymentPayload)
	ctx.Step(`^a valid payment payload with external_reference "([^"]*)"$`, givenAPaymentPayloadWithExternalReference)
	ctx.Step(`^an invalid JSON payload$`, givenAnInvalidJSONPayload)
	ctx.Step(`^a payment payload with invalid fields$`, givenAPaymentPayloadWithInvalidFields)

	ctx.Step(`^I send a POST request to /api/v1/payments/authorize$`, iSendAPostRequestToAuthorizePayment)
	ctx.Step(`^I send another POST request to /api/v1/payments/authorize with same external_reference "([^"]*)"$`, iSendAnotherPostRequestToAuthorizeWithSameExternalReference)

	ctx.Step(`^the response status should be (\d+)$`, theResponseStatusShouldBe)
	ctx.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
}
