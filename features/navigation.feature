Feature: Track navigation
  As a HEOS CLI user
  I want to skip between tracks

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Play next track
    Given the speaker accepts navigation commands
    When I play next for player 1
    Then the operation should succeed

  Scenario: Play previous track
    Given the speaker accepts navigation commands
    When I play previous for player 1
    Then the operation should succeed
