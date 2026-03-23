Feature: Volume control
  As a HEOS CLI user
  I want to control player and group volume

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Get player volume
    Given player 1 has volume 42
    When I get volume for player 1
    Then the volume should be "42"

  Scenario: Set player volume
    Given the speaker accepts volume changes
    When I set volume to 50 for player 1
    Then the operation should succeed

  Scenario: Volume up
    Given the speaker accepts volume changes
    When I increase volume for player 1
    Then the operation should succeed

  Scenario: Volume down
    Given the speaker accepts volume changes
    When I decrease volume for player 1
    Then the operation should succeed

  Scenario: Get group volume
    Given group 1 has volume 60
    When I get volume for group 1
    Then the volume should be "60"

  Scenario: Set group volume
    Given the speaker accepts group volume changes
    When I set volume to 30 for group 1
    Then the operation should succeed
