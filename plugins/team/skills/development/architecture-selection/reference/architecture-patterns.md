# Architecture Pattern Details

Detailed descriptions, diagrams, and trade-offs for each architecture pattern.

---

## Monolithic Architecture

A single deployable unit containing all functionality.

```
┌─────────────────────────────────────────────────────────────┐
│                    Monolithic Application                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Web UI     │  │  API Layer  │  │  Admin UI   │         │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘         │
│         │                │                │                 │
│         └────────────────┼────────────────┘                 │
│                          │                                  │
│  ┌───────────────────────┴───────────────────────────┐     │
│  │              Business Logic Layer                  │     │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐        │     │
│  │  │ Orders   │  │ Users    │  │ Products │        │     │
│  │  └──────────┘  └──────────┘  └──────────┘        │     │
│  └───────────────────────┬───────────────────────────┘     │
│                          │                                  │
│  ┌───────────────────────┴───────────────────────────┐     │
│  │              Data Access Layer                     │     │
│  └───────────────────────┬───────────────────────────┘     │
│                          │                                  │
└──────────────────────────┼──────────────────────────────────┘
                           │
                    ┌──────┴──────┐
                    │  Database   │
                    └─────────────┘
```

**When to Use:**
- Small team (< 10 developers)
- Simple domain
- Rapid iteration needed
- Limited infrastructure expertise

**Trade-offs:**

| Pros | Cons |
|------|------|
| Simple deployment | Limited scalability |
| Easy debugging | Large codebase to manage |
| Single codebase | Technology lock-in |
| Fast development initially | Team coupling |
| Transactional consistency | Full redeploy for changes |

---

## Microservices Architecture

Independently deployable services organized around business capabilities.

```
┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐
│ Web UI │   │Mobile  │   │ Admin  │   │External│
└───┬────┘   └───┬────┘   └───┬────┘   └───┬────┘
    │            │            │            │
    └────────────┴────────────┴────────────┘
                       │
              ┌────────┴────────┐
              │   API Gateway   │
              └────────┬────────┘
                       │
    ┌──────────────────┼──────────────────┐
    │                  │                  │
┌───┴───┐         ┌────┴───┐         ┌───┴───┐
│ Order │         │ User   │         │Product│
│Service│         │Service │         │Service│
├───────┤         ├────────┤         ├───────┤
│  DB   │         │   DB   │         │  DB   │
└───────┘         └────────┘         └───────┘
    │                  │                  │
    └──────────────────┴──────────────────┘
                       │
              ┌────────┴────────┐
              │  Message Bus    │
              └─────────────────┘
```

**When to Use:**
- Large team (> 20 developers)
- Complex, evolving domain
- Independent scaling needed
- Different tech stacks for different services
- High availability requirements

**Trade-offs:**

| Pros | Cons |
|------|------|
| Independent deployment | Operational complexity |
| Technology flexibility | Network latency |
| Team autonomy | Distributed debugging |
| Targeted scaling | Data consistency challenges |
| Fault isolation | More infrastructure |

---

## Event-Driven Architecture

Services communicate through events rather than direct calls.

```
┌─────────────────────────────────────────────────────────────┐
│                      Event Bus / Broker                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   OrderPlaced    UserCreated    PaymentReceived             │
│                                                             │
└──────┬──────────────┬──────────────┬───────────────────────┘
       │              │              │
       ▼              ▼              ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│   Order     │ │   User      │ │  Payment    │
│   Service   │ │   Service   │ │   Service   │
│             │ │             │ │             │
│ Publishes:  │ │ Publishes:  │ │ Publishes:  │
│ OrderPlaced │ │ UserCreated │ │ PaymentRcvd │
│             │ │             │ │             │
│ Subscribes: │ │ Subscribes: │ │ Subscribes: │
│ PaymentRcvd │ │ OrderPlaced │ │ OrderPlaced │
└─────────────┘ └─────────────┘ └─────────────┘
```

**When to Use:**
- Loose coupling required
- Asynchronous processing acceptable
- Complex workflows spanning multiple services
- Audit trail needed
- Event sourcing scenarios

**Trade-offs:**

| Pros | Cons |
|------|------|
| Temporal decoupling | Eventual consistency |
| Natural audit log | Complex debugging |
| Scalability | Message ordering challenges |
| Extensibility | Infrastructure requirements |
| Resilience | Learning curve |

---

## Serverless Architecture

Functions executed on-demand without managing servers.

```
┌────────────────────────────────────────────────────────────┐
│                         Client                              │
└────────────────────────────┬───────────────────────────────┘
                             │
┌────────────────────────────┴───────────────────────────────┐
│                      API Gateway                            │
└────────────────────────────┬───────────────────────────────┘
                             │
    ┌────────────────────────┼────────────────────────┐
    │                        │                        │
    ▼                        ▼                        ▼
┌──────────┐          ┌──────────┐          ┌──────────┐
│ Function │          │ Function │          │ Function │
│ GetUser  │          │CreateOrder│         │ SendEmail│
└────┬─────┘          └────┬─────┘          └────┬─────┘
     │                     │                     │
     ▼                     ▼                     ▼
┌──────────┐          ┌──────────┐          ┌──────────┐
│ Database │          │  Queue   │          │  Email   │
│          │          │          │          │ Service  │
└──────────┘          └──────────┘          └──────────┘
```

**When to Use:**
- Variable/unpredictable workloads
- Event-triggered processing
- Cost optimization for low traffic
- Rapid development needed
- Short-running operations

**Trade-offs:**

| Pros | Cons |
|------|------|
| No server management | Cold start latency |
| Pay-per-use | Execution time limits |
| Auto-scaling | Vendor lock-in |
| Rapid deployment | Complex local development |
| Reduced ops burden | Stateless constraints |
