Feature: Player commands
  As a HEOS CLI user
  I want to discover and control my HEOS players

  Background:
    Given I am connected to a HEOS speaker

  Scenario: List players
    Given the speaker has 2 players
    When I list players
    Then I should get 2 players

  Scenario: List players when none exist
    Given the speaker has 0 players
    When I list players
    Then I should get 0 players

  Scenario: Get player info
    Given the speaker has player "Bureau" with pid 1
    When I get info for player 1
    Then the player name should be "Bureau"

  Scenario: Get play state
    Given player 1 is "play"
    When I get play state for player 1
    Then the state should be "play"

  Scenario: Get paused state
    Given player 1 is "pause"
    When I get play state for player 1
    Then the state should be "pause"

  Scenario: Set play state
    Given the speaker accepts play state changes
    When I set play state "play" for player 1
    Then the operation should succeed
