# Strategic Patterns

## Bounded Context

A bounded context defines the boundary within which a domain model applies. The same term can mean different things in different contexts.

```
Example: "Customer" in different contexts

+------------------+  +------------------+  +------------------+
|    Sales         |  |    Support       |  |    Billing       |
|    Context       |  |    Context       |  |    Context       |
+------------------+  +------------------+  +------------------+
| Customer:        |  | Customer:        |  | Customer:        |
| - Leads          |  | - Tickets        |  | - Invoices       |
| - Opportunities  |  | - SLA            |  | - Payment        |
| - Proposals      |  | - Satisfaction   |  | - Credit Limit   |
+------------------+  +------------------+  +------------------+
```

### Context Identification

Ask these questions to find context boundaries:
- Where does the ubiquitous language change?
- Which teams own which concepts?
- Where do integration points naturally occur?
- What could be deployed independently?

## Context Mapping

Define how bounded contexts integrate:

| Pattern | Description | Use When |
|---------|-------------|----------|
| **Shared Kernel** | Shared code between contexts | Close collaboration, same team |
| **Customer-Supplier** | Upstream/downstream relationship | Clear dependency direction |
| **Conformist** | Downstream adopts upstream model | No negotiation power |
| **Anti-Corruption Layer** | Translation layer between models | Protecting domain from external models |
| **Open Host Service** | Published API for integration | Multiple consumers |
| **Published Language** | Shared interchange format | Industry standards exist |

## Ubiquitous Language

The shared vocabulary between developers and domain experts.

### Building Ubiquitous Language

1. EXTRACT terms from domain expert conversations
2. DOCUMENT in a glossary with precise definitions
3. ENFORCE in code — class names, method names, variables
4. EVOLVE as understanding deepens

### Glossary Entry Format

```
Term: Order
Definition: A confirmed request from a customer to purchase
            one or more products at agreed prices.
NOT: A shopping cart (which is an Intent, not an Order)
Context: Sales
```
