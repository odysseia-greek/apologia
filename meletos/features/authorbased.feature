Feature: Sokrates GraphQL API for AuthorBased Quiz Functionality
  As a user of the Odysseia Greek application
  I want to interact with the Sokrates GraphQL API
  So that I can access and complete author based quizzes

  Background:
    Given the graphql backend is running

  Scenario: Check health status of the authorbased service
    When I query the health status
    Then the service "author-based" should be healthy
    And the version information should be available for "author-based"
    And basic database health info should be available for "author-based"

  Scenario: Answer a authorbased quiz
    When I query for authorbased quiz options
    And I use the authorbased options to create a question
    And I submit each authorbased option once
    Then I should have 3 incorrect and 1 correct answer
    And The Progress should be 3 incorrect and 1 correct answer

  Scenario: Some quizzes have grammar quizzes embedded
    When I query for authorbased quiz options
    And I create a quiz that has the name "John - The New Testament"
    Then grammar options should be embedded into the quiz for some words

  Scenario: Creating a quiz will also return the metadata for a authorbased sentence in Greek, English and a reference
    When I query for authorbased quiz options
    And I use the authorbased options to create a question
    Then that question has a Greek, English sentence
    And that question has a reference to the text module

  Scenario: Some quizzes have grammar quizzes embedded
    When I query for authorbased quiz options
    And I query the word forms for a segment
    Then the words should be returned as they appear in the text