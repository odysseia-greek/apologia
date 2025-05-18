Feature: Sokrates GraphQL API for Media Quiz Functionality
  As a user of the Odysseia Greek application
  I want to interact with the Sokrates GraphQL API
  So that I can access and complete media quizzes

  Background:
    Given the graphql backend is running

  Scenario: Check health status of the Media service
    When I query the health status
    Then the service "media" should be healthy
    And the version information should be available for "media"
    And basic database health info should be available for "media"

  Scenario: Answer a media quiz
    When I query for media quiz options
    And I use the media options to create a question
    And I submit each option once
    Then I should have 3 incorrect and 1 correct answer
    And The Progress should be 3 incorrect and 1 correct answer
