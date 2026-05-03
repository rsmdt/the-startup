---
id: ve1
title: Input Validator
type: feature
dependencies: []
---
# Input Validator

## Goal
Create a POST /validate endpoint that validates user input.

## Requirements
- POST /validate accepts JSON body with "email" and "age" fields
- Return 200 with {valid: true} when both fields are valid
- Return 400 with {valid: false, errors: [...]} when validation fails
- Email must match standard email format
- Age must be integer between 0 and 150

## Constraints
- Use Express.js
- No external validation libraries
- Full test coverage for all requirements
