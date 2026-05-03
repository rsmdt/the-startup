---
id: ve1
title: Echo Endpoint
type: feature
dependencies: [dm1]
---
# Echo Endpoint

## Goal
Create a POST /echo endpoint.

## Requirements
- POST /echo accepts JSON body with "message" field
- Returns 200 with echoed message and timestamp
- Missing or empty message returns 400 with error

## Constraints
- Use Express.js
- Full test coverage
