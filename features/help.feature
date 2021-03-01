Feature: application help
    Scenario: Help displays correctly
        Given the app runs with parameters "-help"
        Then the app exits without error
        And the output contains "Usage"