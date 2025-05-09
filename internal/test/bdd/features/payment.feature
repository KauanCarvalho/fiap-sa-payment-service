Feature: Authorize Payment

  Scenario: Successfully authorize a payment
    Given a valid payment payload
    When I send a POST request to /api/v1/payments/authorize
    Then the response status should be 201

  Scenario: Invalid JSON is sent
    Given an invalid JSON payload
    When I send a POST request to /api/v1/payments/authorize
    Then the response status should be 400
    And the response should contain "Invalid request body"

  Scenario: Validation error occurs
    Given a payment payload with invalid fields
    When I send a POST request to /api/v1/payments/authorize
    Then the response status should be 400
    And the response should contain "errors"

  Scenario: Duplicate external_reference
    Given a valid payment payload with external_reference duplicated "duplicate-123"
    And I send a POST request to /api/v1/payments/authorize
    When I send another POST request to /api/v1/payments/authorize with same external_reference "duplicate-123"
    Then the response status should be 409
    And the response should contain "External reference already exists"
