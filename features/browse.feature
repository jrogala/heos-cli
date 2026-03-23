Feature: Browse and play
  As a HEOS CLI user
  I want to browse music sources and play URLs

  Background:
    Given I am connected to a HEOS speaker

  Scenario: List music sources
    Given the speaker has music sources
    When I list music sources
    Then I should get 3 sources

  Scenario: Play a URL
    Given the speaker accepts play-url
    When I play URL "http://example.com/stream.mp3" on player 1
    Then the operation should succeed
