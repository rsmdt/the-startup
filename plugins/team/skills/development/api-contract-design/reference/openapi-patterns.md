# OpenAPI/Swagger Patterns

Specification structure and reusable component patterns for OpenAPI 3.1.

---

## Specification Structure

```yaml
openapi: 3.1.0
info:
  title: Example API
  version: 1.0.0
  description: API description with markdown support
  contact:
    name: API Support
    url: https://example.com/support

servers:
  - url: https://api.example.com/v1
    description: Production
  - url: https://api.staging.example.com/v1
    description: Staging

security:
  - bearerAuth: []

paths:
  /users:
    get:
      operationId: listUsers
      summary: List all users
      tags: [Users]
      # ... operation details

components:
  schemas:
    User:
      type: object
      required: [id, email]
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
          format: email
```

## Reusable Components

```yaml
components:
  schemas:
    # Reusable pagination
    PaginationMeta:
      type: object
      properties:
        total:
          type: integer
        page:
          type: integer
        perPage:
          type: integer

    # Reusable error
    Error:
      type: object
      required: [code, message]
      properties:
        code:
          type: string
        message:
          type: string

  parameters:
    # Reusable query params
    PageParam:
      name: page
      in: query
      schema:
        type: integer
        default: 1
        minimum: 1

  responses:
    # Reusable responses
    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
```
