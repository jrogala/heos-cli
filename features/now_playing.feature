Feature: Now playing
  As a HEOS CLI user
  I want to see what is currently playing

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Get now playing media
    Given player 1 is playing "Overnight" by "Parcels" from "Day/Night"
    When I get now playing for player 1
    Then the song should be "Overnight"
    And the artist should be "Parcels"
