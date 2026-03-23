Feature: Group commands
  As a HEOS CLI user
  I want to manage speaker groups

  Background:
    Given I am connected to a HEOS speaker

  Scenario: List groups
    Given the speaker has 2 groups
    When I list groups
    Then I should get 2 groups

  Scenario: No groups
    Given the speaker has 0 groups
    When I list groups
    Then I should get 0 groups
