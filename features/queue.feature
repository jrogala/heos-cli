Feature: Queue management
  As a HEOS CLI user
  I want to manage the player queue

  Background:
    Given I am connected to a HEOS speaker

  Scenario: Get queue with items
    Given player 1 has 5 items in queue
    When I get queue for player 1
    Then I should get 5 queue items

  Scenario: Get empty queue
    Given player 1 has 0 items in queue
    When I get queue for player 1
    Then I should get 0 queue items

  Scenario: Play queue item
    Given the speaker accepts queue operations
    When I play queue item 3 for player 1
    Then the operation should succeed

  Scenario: Clear queue
    Given the speaker accepts queue operations
    When I clear queue for player 1
    Then the operation should succeed
