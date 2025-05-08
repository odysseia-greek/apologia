Feature: Sokrates GraphQL API for Dialogue Quiz Functionality
  As a user of the Odysseia Greek application
  I want to interact with the Sokrates GraphQL API
  So that I can access and complete dialogue quizzes

  Background:
    Given the graphql backend is running


  Scenario: Check health status of the Dialogue service
    When I query the health status
    Then the service "dialogue" should be healthy
    And the version information should be available for "dialogue"
    And basic database health info should be available for "dialogue"

  Scenario: Answer a Dialogue quiz with a mistake
    When I query for dialogue quiz options
    And I use the dialogue options to create a question
    And I submit with at least one section wronly placed
    Then the percentage should be lower than 100
    And wronglyPlaced should hold a reference to the correct place

  @wip
  Scenario: Answer a Dialogue quiz without a mistake
    When I query for dialogue quiz options
    And I use the dialogue options to create a question
    And I submit with a perfect input
    Then the percentage should be 100
    And wronglyPlaced should be empty