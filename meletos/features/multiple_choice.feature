Feature: Sokrates GraphQL API for Multiple Choice Quiz Functionality
  As a user of the Odysseia Greek application
  I want to interact with the Sokrates GraphQL API
  So that I can access and complete multiple choice quizzes

  Background:
    Given the graphql backend is running


  Scenario: Check health status of the Media service
    When I query the health status
    Then the service "multiple-choice" should be healthy
    And the version information should be available for "multiple-choice"
    And basic database health info should be available for "multiple-choice"

  Scenario: Answer a multiple quiz
    When I query for multiple choice quiz options
    And I use the multiple choice options to create a question
    And I submit each multiple choice option once
    Then I should have 3 incorrect and 1 correct answer
    And The Progress should be 3 incorrect and 1 correct answer