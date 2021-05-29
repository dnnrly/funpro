Feature: application help

    @Integration
    Scenario: Can get a 200 response using an ALB request
        Given I make a "GET" to "http://funpro:8080/alb-responder/ok"
        Then the response code is 200
        # And the "alb-responder" recording at "ok" matches "simple-ok.json"