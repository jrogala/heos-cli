Feature: System commands
  As a HEOS CLI user
  I want to check system status and manage my account

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Heartbeat
    Given the speaker responds to heartbeat
    When I send a heartbeat
    Then the operation should succeed

  Scenario: Check account when signed in
    Given the speaker has account "user@example.com" signed in
    When I check the account
    Then I should see username "user@example.com"

  Scenario: Check account when signed out
    Given the speaker has no account signed in
    When I check the account
    Then I should see signed out

  Scenario: Sign in with valid credentials
    Given the speaker accepts sign-in
    When I sign in with "user@example.com" and "password123"
    Then the operation should succeed

  Scenario: Sign in with invalid credentials
    Given the speaker rejects sign-in
    When I sign in with "user@example.com" and "wrong"
    Then the operation should fail
    And the error should contain "Invalid_Credentials"

  Scenario: Sign out
    Given I am connected to a HEOS speaker
    When I sign out
    Then the operation should succeed
