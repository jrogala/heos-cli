Feature: Play mode
  As a HEOS CLI user
  I want to control repeat and shuffle settings

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Get play mode
    Given player 1 has play mode repeat "on_all" shuffle "off"
    When I get play mode for player 1
    Then the repeat mode should be "on_all"
    And the shuffle mode should be "off"

  Scenario: Get shuffle mode
    Given player 1 has play mode repeat "off" shuffle "on"
    When I get play mode for player 1
    Then the shuffle mode should be "on"

  Scenario: Set play mode
    Given the speaker accepts play mode changes
    When I set play mode repeat "on_one" shuffle "on" for player 1
    Then the operation should succeed
