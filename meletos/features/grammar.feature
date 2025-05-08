Feature: Sokrates GraphQL API for Grammar Quiz Functionality
  As a user of the Odysseia Greek application
  I want to interact with the Sokrates GraphQL API
  So that I can access and complete grammar quizzes

  Background:
    Given the graphql backend is running

  Scenario: Check health status of the Grammar service
    When I query the health status
    Then the service "grammar" should be healthy
    And the version information should be available for "grammar"
    And basic database health info should be available for "grammar"

  Scenario: Answer a grammar quiz
    When I query for grammar quiz options
    And I use the grammar options to create a question
    And I submit each grammar option once
    Then I should have 3 incorrect and 1 correct answer
    And The Progress should be 3 incorrect and 1 correct answer
