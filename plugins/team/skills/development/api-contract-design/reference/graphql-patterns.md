# GraphQL Patterns

Schema design, query/mutation patterns, and N+1 prevention strategies.

---

## Schema Design Principles

```graphql
# Use clear, descriptive type names
type User {
  id: ID!
  email: String!
  displayName: String!
  createdAt: DateTime!

  # Relationships with clear naming
  organization: Organization
  orders(first: Int, after: String): OrderConnection!
}

# Use connections for paginated lists
type OrderConnection {
  edges: [OrderEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type OrderEdge {
  node: Order!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

## Query Design

```graphql
type Query {
  # Single resource by ID
  user(id: ID!): User

  # List with filtering and pagination
  users(
    filter: UserFilter
    first: Int
    after: String
    orderBy: UserOrderBy
  ): UserConnection!

  # Viewer pattern for current user
  viewer: User
}

input UserFilter {
  status: UserStatus
  organizationId: ID
  searchQuery: String
}

enum UserOrderBy {
  CREATED_AT_ASC
  CREATED_AT_DESC
  NAME_ASC
  NAME_DESC
}
```

## Mutation Design

```graphql
type Mutation {
  # Use input types for complex mutations
  createUser(input: CreateUserInput!): CreateUserPayload!
  updateUser(input: UpdateUserInput!): UpdateUserPayload!
  deleteUser(id: ID!): DeleteUserPayload!
}

input CreateUserInput {
  email: String!
  displayName: String!
  organizationId: ID
}

# Payload types for consistent responses
type CreateUserPayload {
  user: User
  errors: [UserError!]!
}

type UserError {
  field: String
  code: String!
  message: String!
}
```

## N+1 Query Prevention

| Strategy | Description |
|----------|-------------|
| **DataLoader** | Batch and cache database lookups per request |
| **Query complexity analysis** | Score and limit query cost before execution |
| **Depth limiting** | Cap nested query depth (typically 5-10 levels) |
| **Field-level cost** | Assign cost per field, reject expensive queries |
| **Persisted queries** | Pre-approve queries for production use |
