# Template: GraphQL Schema Design

## Purpose

Use this template to design GraphQL schemas with consistent patterns for types, queries, mutations, and subscriptions. Follows Relay-style connections for pagination and input/payload patterns for mutations.

## Template

```graphql
# =============================================================================
# [API Name] GraphQL Schema
# =============================================================================
#
# Description: [Brief description of the API's purpose]
# Version: 1.0.0
# Last Updated: [Date]
#
# Conventions:
# - Types use PascalCase
# - Fields use camelCase
# - Enums use SCREAMING_SNAKE_CASE
# - Connections follow Relay specification
# - Mutations use Input/Payload pattern
# =============================================================================

# -----------------------------------------------------------------------------
# SCALARS
# -----------------------------------------------------------------------------

"""
A datetime string in ISO 8601 format: YYYY-MM-DDTHH:mm:ss.sssZ
"""
scalar DateTime

"""
A UUID string following RFC 4122
"""
scalar UUID

"""
Arbitrary JSON object for flexible data
"""
scalar JSON

"""
A URL string
"""
scalar URL

"""
An email address string
"""
scalar Email

# -----------------------------------------------------------------------------
# ENUMS
# -----------------------------------------------------------------------------

"""
Sort direction for ordered results
"""
enum SortDirection {
  ASC
  DESC
}

"""
Status values for [Resource]
"""
enum [Resource]Status {
  DRAFT
  ACTIVE
  ARCHIVED
  DELETED
}

# -----------------------------------------------------------------------------
# INTERFACES
# -----------------------------------------------------------------------------

"""
Node interface for Relay-compliant object identification
"""
interface Node {
  """
  Global unique identifier
  """
  id: ID!
}

"""
Timestamped interface for entities with audit fields
"""
interface Timestamped {
  """
  When the entity was created
  """
  createdAt: DateTime!

  """
  When the entity was last updated
  """
  updatedAt: DateTime!
}

# -----------------------------------------------------------------------------
# TYPES: Core Entities
# -----------------------------------------------------------------------------

"""
[Description of what this type represents]
"""
type [Resource] implements Node & Timestamped {
  """
  Unique identifier for the [resource]
  """
  id: ID!

  """
  [Description of field]
  """
  [field]: [Type]!

  """
  [Description of optional field]
  """
  [optionalField]: [Type]

  """
  Current status of the [resource]
  """
  status: [Resource]Status!

  """
  When the [resource] was created
  """
  createdAt: DateTime!

  """
  When the [resource] was last updated
  """
  updatedAt: DateTime!

  # ---------------------------------------------------------------------------
  # Relationships
  # ---------------------------------------------------------------------------

  """
  The [related entity] this [resource] belongs to
  """
  [parent]: [ParentType]!

  """
  [Children] associated with this [resource]
  """
  [children](
    first: Int
    after: String
    last: Int
    before: String
    filter: [Child]Filter
    orderBy: [Child]OrderBy
  ): [Child]Connection!
}

# -----------------------------------------------------------------------------
# TYPES: Connections (Relay-style Pagination)
# -----------------------------------------------------------------------------

"""
A connection to a list of [Resource] items
"""
type [Resource]Connection {
  """
  A list of edges
  """
  edges: [[Resource]Edge!]!

  """
  Information to aid in pagination
  """
  pageInfo: PageInfo!

  """
  Total number of items in the connection
  """
  totalCount: Int!
}

"""
An edge in a [Resource] connection
"""
type [Resource]Edge {
  """
  The item at the end of the edge
  """
  node: [Resource]!

  """
  A cursor for use in pagination
  """
  cursor: String!
}

"""
Information about pagination in a connection
"""
type PageInfo {
  """
  When paginating forwards, are there more items?
  """
  hasNextPage: Boolean!

  """
  When paginating backwards, are there more items?
  """
  hasPreviousPage: Boolean!

  """
  When paginating backwards, the cursor to continue
  """
  startCursor: String

  """
  When paginating forwards, the cursor to continue
  """
  endCursor: String
}

# -----------------------------------------------------------------------------
# INPUTS: Filtering and Ordering
# -----------------------------------------------------------------------------

"""
Filter criteria for [Resource] queries
"""
input [Resource]Filter {
  """
  Filter by status
  """
  status: [Resource]Status

  """
  Filter by [parent] ID
  """
  [parent]Id: ID

  """
  Filter by creation date range
  """
  createdAfter: DateTime

  """
  Filter by creation date range
  """
  createdBefore: DateTime

  """
  Full-text search query
  """
  searchQuery: String

  """
  Combine multiple filters with AND
  """
  AND: [[Resource]Filter!]

  """
  Combine multiple filters with OR
  """
  OR: [[Resource]Filter!]
}

"""
Ordering options for [Resource] queries
"""
input [Resource]OrderBy {
  """
  Field to order by
  """
  field: [Resource]OrderField!

  """
  Sort direction
  """
  direction: SortDirection!
}

"""
Fields available for ordering [Resource] queries
"""
enum [Resource]OrderField {
  CREATED_AT
  UPDATED_AT
  [FIELD_NAME]
}

# -----------------------------------------------------------------------------
# QUERIES
# -----------------------------------------------------------------------------

type Query {
  # ---------------------------------------------------------------------------
  # Node Query (Relay requirement)
  # ---------------------------------------------------------------------------

  """
  Fetch any node by its global ID
  """
  node(
    """
    The global ID of the node
    """
    id: ID!
  ): Node

  """
  Fetch multiple nodes by their global IDs
  """
  nodes(
    """
    The global IDs of the nodes
    """
    ids: [ID!]!
  ): [Node]!

  # ---------------------------------------------------------------------------
  # Viewer Query (Current User Context)
  # ---------------------------------------------------------------------------

  """
  The currently authenticated user
  """
  viewer: User

  # ---------------------------------------------------------------------------
  # [Resource] Queries
  # ---------------------------------------------------------------------------

  """
  Fetch a single [resource] by ID
  """
  [resource](
    """
    The unique identifier of the [resource]
    """
    id: ID!
  ): [Resource]

  """
  Fetch a paginated list of [resources]
  """
  [resources](
    """
    Number of items to fetch (forward pagination)
    """
    first: Int

    """
    Cursor to start after (forward pagination)
    """
    after: String

    """
    Number of items to fetch (backward pagination)
    """
    last: Int

    """
    Cursor to start before (backward pagination)
    """
    before: String

    """
    Filter criteria
    """
    filter: [Resource]Filter

    """
    Ordering specification
    """
    orderBy: [Resource]OrderBy
  ): [Resource]Connection!
}

# -----------------------------------------------------------------------------
# MUTATIONS
# -----------------------------------------------------------------------------

type Mutation {
  # ---------------------------------------------------------------------------
  # [Resource] Mutations
  # ---------------------------------------------------------------------------

  """
  Create a new [resource]
  """
  create[Resource](
    """
    Input data for the new [resource]
    """
    input: Create[Resource]Input!
  ): Create[Resource]Payload!

  """
  Update an existing [resource]
  """
  update[Resource](
    """
    Input data for the update
    """
    input: Update[Resource]Input!
  ): Update[Resource]Payload!

  """
  Delete a [resource]
  """
  delete[Resource](
    """
    ID of the [resource] to delete
    """
    id: ID!
  ): Delete[Resource]Payload!

  """
  [Custom action] on a [resource]
  """
  [action][Resource](
    """
    Input data for the action
    """
    input: [Action][Resource]Input!
  ): [Action][Resource]Payload!
}

# -----------------------------------------------------------------------------
# INPUTS: Mutations
# -----------------------------------------------------------------------------

"""
Input for creating a [resource]
"""
input Create[Resource]Input {
  """
  [Description of required field]
  """
  [requiredField]: [Type]!

  """
  [Description of optional field]
  """
  [optionalField]: [Type]

  """
  Idempotency key for safe retries
  """
  clientMutationId: String
}

"""
Input for updating a [resource]
"""
input Update[Resource]Input {
  """
  ID of the [resource] to update
  """
  id: ID!

  """
  [Description of updatable field]
  """
  [field]: [Type]

  """
  Idempotency key for safe retries
  """
  clientMutationId: String
}

"""
Input for [custom action] on a [resource]
"""
input [Action][Resource]Input {
  """
  ID of the [resource]
  """
  id: ID!

  """
  [Action-specific parameters]
  """
  [param]: [Type]!

  """
  Idempotency key for safe retries
  """
  clientMutationId: String
}

# -----------------------------------------------------------------------------
# PAYLOADS: Mutations
# -----------------------------------------------------------------------------

"""
Payload returned from create[Resource] mutation
"""
type Create[Resource]Payload {
  """
  The created [resource], null if errors occurred
  """
  [resource]: [Resource]

  """
  Errors that occurred during the mutation
  """
  errors: [UserError!]!

  """
  Client mutation ID for request correlation
  """
  clientMutationId: String
}

"""
Payload returned from update[Resource] mutation
"""
type Update[Resource]Payload {
  """
  The updated [resource], null if errors occurred
  """
  [resource]: [Resource]

  """
  Errors that occurred during the mutation
  """
  errors: [UserError!]!

  """
  Client mutation ID for request correlation
  """
  clientMutationId: String
}

"""
Payload returned from delete[Resource] mutation
"""
type Delete[Resource]Payload {
  """
  ID of the deleted [resource]
  """
  deletedId: ID

  """
  Errors that occurred during the mutation
  """
  errors: [UserError!]!

  """
  Client mutation ID for request correlation
  """
  clientMutationId: String
}

"""
Payload returned from [action][Resource] mutation
"""
type [Action][Resource]Payload {
  """
  The [resource] after the action
  """
  [resource]: [Resource]

  """
  [Action-specific result fields]
  """
  [resultField]: [Type]

  """
  Errors that occurred during the mutation
  """
  errors: [UserError!]!

  """
  Client mutation ID for request correlation
  """
  clientMutationId: String
}

# -----------------------------------------------------------------------------
# ERROR TYPES
# -----------------------------------------------------------------------------

"""
A user-facing error from a mutation
"""
type UserError {
  """
  The field that caused the error, if applicable
  Path format: ["input", "fieldName"]
  """
  path: [String!]

  """
  Machine-readable error code
  """
  code: UserErrorCode!

  """
  Human-readable error message
  """
  message: String!
}

"""
Error codes for user-facing errors
"""
enum UserErrorCode {
  """
  A required field was not provided
  """
  REQUIRED_FIELD

  """
  The provided value is invalid
  """
  INVALID_VALUE

  """
  The resource was not found
  """
  NOT_FOUND

  """
  The operation conflicts with current state
  """
  CONFLICT

  """
  The user is not authorized for this operation
  """
  UNAUTHORIZED

  """
  A business rule was violated
  """
  BUSINESS_RULE_VIOLATION

  """
  Rate limit exceeded
  """
  RATE_LIMITED
}

# -----------------------------------------------------------------------------
# SUBSCRIPTIONS (if applicable)
# -----------------------------------------------------------------------------

type Subscription {
  """
  Subscribe to [resource] changes
  """
  [resource]Changed(
    """
    Optional filter to specific [resource] ID
    """
    id: ID
  ): [Resource]ChangePayload!
}

"""
Payload for [resource] change subscription
"""
type [Resource]ChangePayload {
  """
  The type of change that occurred
  """
  changeType: ChangeType!

  """
  The [resource] after the change (null for DELETE)
  """
  [resource]: [Resource]

  """
  ID of the affected [resource]
  """
  [resource]Id: ID!
}

"""
Types of changes for subscriptions
"""
enum ChangeType {
  CREATED
  UPDATED
  DELETED
}
```

## Usage Instructions

1. Copy the template into your schema file(s)
2. Replace all `[bracketed]` placeholders with actual values
3. Remove sections that don't apply (e.g., Subscriptions if not needed)
4. Add your custom scalars if using libraries like graphql-scalars
5. Ensure all types implement Node if using Relay
6. Add descriptions to ALL types and fields

## Design Patterns Applied

### Relay Specification Compliance

- **Node interface**: All fetchable types implement `Node` with global ID
- **Connection pattern**: Paginated lists use `edges`, `node`, `cursor`, `pageInfo`
- **Global IDs**: Use base64-encoded `type:id` format for Relay compatibility

### Input/Payload Pattern

- **Inputs**: All mutation arguments wrapped in a single `input` argument
- **Payloads**: All mutations return a payload with the result AND errors array
- **clientMutationId**: Included for idempotency and request correlation

### Error Handling

- **UserError type**: Field-level errors with path, code, and message
- **Non-null errors array**: Always returns `errors: [UserError!]!` (empty if success)
- **Nullable result**: Main result field is nullable (null when errors exist)

## Customization Options

### For Simpler APIs

Remove Relay patterns if not needed:

```graphql
# Instead of Connection pattern
type Query {
  users(limit: Int, offset: Int): [User!]!
}
```

### For Real-time Features

Expand subscriptions with more granular events:

```graphql
type Subscription {
  resourceCreated: Resource!
  resourceUpdated(id: ID!): Resource!
  resourceDeleted: ID!
}
```

### For Federation

Add directives for Apollo Federation:

```graphql
type User @key(fields: "id") {
  id: ID!
  # ... fields
}

extend type Query {
  _entities(representations: [_Any!]!): [_Entity]!
  _service: _Service!
}
```

## Anti-patterns to Avoid

```graphql
# BAD: Inconsistent nullability
type User {
  name: String      # Should be String! if always present
  email: String!    # Good
}

# BAD: No connection pattern for lists
type Query {
  allUsers: [User]  # Use UserConnection for pagination
}

# BAD: Generic mutation payloads
type MutationResult {
  success: Boolean!
  message: String
}

# BAD: No error handling in mutations
type Mutation {
  createUser(name: String!): User  # Should return CreateUserPayload
}

# BAD: Deeply nested inputs
input CreateOrderInput {
  items: [CreateOrderItemInput!]!  # OK, one level
}
input CreateOrderItemInput {
  product: CreateProductInput!     # BAD, creates products inline
}
```

## Schema Checklist

Before finalizing:
- [ ] All types have descriptions
- [ ] All fields have descriptions
- [ ] Nullability is intentional and consistent
- [ ] Connections used for paginated lists
- [ ] Mutations return payloads with errors array
- [ ] Enums use SCREAMING_SNAKE_CASE
- [ ] Input types are suffixed with `Input`
- [ ] Payload types are suffixed with `Payload`
- [ ] Filter types are suffixed with `Filter`
- [ ] N+1 query potential identified and addressed
- [ ] Deprecated fields marked with `@deprecated(reason: "...")`
- [ ] No circular type dependencies without interfaces
