Feature: Mute control
  As a HEOS CLI user
  I want to control mute state on players and groups

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Player is muted
    Given player 1 is muted
    When I get mute state for player 1
    Then the mute state should be "on"

  Scenario: Player is not muted
    Given player 1 is not muted
    When I get mute state for player 1
    Then the mute state should be "off"

  Scenario: Mute a player
    Given the speaker accepts mute changes
    When I mute player 1
    Then the operation should succeed

  Scenario: Unmute a player
    Given the speaker accepts mute changes
    When I unmute player 1
    Then the operation should succeed

  Scenario: Toggle mute
    Given the speaker accepts mute changes
    When I toggle mute for player 1
    Then the operation should succeed

  Scenario: Group is muted
    Given group 1 is muted
    When I get mute state for group 1
    Then the mute state should be "on"

  Scenario: Group is not muted
    Given group 1 is not muted
    When I get mute state for group 1
    Then the mute state should be "off"

  Scenario: Toggle group mute
    Given the speaker accepts group mute changes
    When I toggle mute for group 1
    Then the operation should succeed
