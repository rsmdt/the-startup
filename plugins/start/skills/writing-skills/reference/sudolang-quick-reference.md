# SudoLang Quick Reference

Compact syntax reference for writing skills.

---

## Interfaces

User {
  id: String
  displayName             // type inferred when omitted
  permissions[]           // [] for arrays
  metadata{}              // {} for open objects
  nickname?: String       // ? for optional fields
}

---

## Constraints

Constraints {
  require { ... }         // must do
  never { ... }           // must not do
}

Constraints are declarative — describe desired state, not procedures. Can be nested inside functions for per-step rules:

fn documentSection(section) {
  Update SDD with research findings.

  Constraints {
    require {
      PRD requirements for this section are read and understood.
      User has confirmed all architecture decisions before proceeding.
    }
  }
}

---

## Functions

fn foo()                    // inferred — AI fills in behavior (forward declaration)
fn foo() { ... }            // defined — body provided, not executed (definition)
foo() { ... }               // executed — runs immediately (entry point)

fn bar():modifier=value     // modifier shapes behavior at call time

Convention: `fn` = "define this for later". No `fn` = "run this now".

---

## Match

match (target) {
  /^\d+$/    => handleNumber       // regex pattern
  "staged"   => handleStaged       // exact string
  default    => handleDefault      // fallback
}

---

## Pipe Operator

data |> normalize |> filter |> sort

Chains function calls left-to-right. Use for data pipelines and workflow sequencing.

---

## Template Strings

"Hello $name, you have $count items"

Use double quotes. Avoid single quotes — they break syntax highlighting in SudoLang and conflict with natural language apostrophes.

---

## Operators

| Operator | Usage | Example |
|----------|-------|---------|
| `match` | Decision routing | `match (x) { pattern => result }` |
| `\|>` | Function composition / pipelines | `data \|> filter \|> sort` |
| `?` | Optional fields | `nickname?: String` |
| `=>` | Match arm results | `"staged" => git diff --cached` |
| `1..10` | Inclusive ranges | `depth: 1..10` |
| `&&` / `\|\|` | Logical operators | `isAdult && isMember` |
| `!` | Negation | `!isBlocked` |

---

## Disallowed Keywords

Do not use: `class`, `extends`, `constructor`, `super`. SudoLang uses interfaces and composition, not class-based inheritance.
