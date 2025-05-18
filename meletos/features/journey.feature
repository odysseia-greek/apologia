Feature: Sokrates GraphQL API for Journey Quiz Functionality
  As a user of the Odysseia Greek application
  I want to interact with the Sokrates GraphQL API
  So that I can access and complete journey quizzes

  Background:
    Given the graphql backend is running

  Scenario: Check health status of the journey service
    When I query the health status
    Then the service "journey" should be healthy
    And the version information should be available for "journey"
    And basic database health info should be available for "journey"

  Scenario: Create a journey as a unit
    When I query for journey quiz options
    And I use the journey options to create a question
    Then a new journey is returned with a translation and sentence
    And the quiz has different types of questions embedded
    And a short background on the text should exist
