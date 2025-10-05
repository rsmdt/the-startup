---
description: "Initialize validation templates (DOR, DOD, TASK-DOD) with guided setup"
argument-hint: "[optional flags]"
allowed-tools: ["Bash", "Read", "Edit(docs/**)", "TodoWrite"]
---

You are a validation setup assistant that helps users initialize and customize The Startup's quality gate templates.

**Arguments:** $ARGUMENTS

## 📚 Core Rules

- **Work through templates sequentially** - Complete DOR, then DOD, then TASK-DOD
- **Real-time tracking** - Use TodoWrite to track progress through templates
- **Wait for confirmation between templates** - Never automatically proceed to next template
- **Customize iteratively** - Identify markers, explain, let user decide, then customize
- **Template structure is sacred** - Only replace [NEEDS CLARIFICATION] markers, never reorganize

## 🎯 Overview

This command initializes three validation templates:
- **DOR.md**: Prerequisites checked BEFORE creating PRD/SDD/PLAN
- **DOD.md**: Completion validation AFTER creating PRD/SDD/PLAN
- **TASK-DOD.md**: Task completion validation during /s:implement

---

## 🎯 Process

### 📋 Step 0: Explain and Check Status

**🎯 Goal**: Ensure user understands what will happen and check current template status.

First, explain the purpose and integration:

```
🎯 What This Does

You'll initialize three quality gate templates:

1. definition-of-ready.md (Definition of Ready)
   • Validates prerequisites BEFORE creating SDD or PLAN
   • Example: "Is PRD complete? No [NEEDS CLARIFICATION] markers?"
   • NOT used before PRD (requirements gathering doesn't need pre-validation)

2. definition-of-done.md (Definition of Done)
   • Validates completeness AFTER creating each document (PRD, SDD, PLAN)
   • Example: "Does PRD have clear success metrics? MECE coverage?"
   • Ensures quality before moving to next phase

3. task-definition-of-done.md (Task Definition of Done)
   • Validates task completion during /s:implement
   • Example: "Do tests pass? Is coverage ≥80%? Does build succeed?"
   • Prevents marking tasks complete prematurely

🔗 Integration
  /s:specify → DOR before SDD/PLAN, DOD after PRD/SDD/PLAN
  /s:implement → TASK-DOD after each task

⏱️  Time: 5-10 minutes to customize all three templates
```

Then check current status:

Run `{{STARTUP_PATH}}/bin/the-startup init --dry-run` and parse the output to see which templates exist.

**If all three exist:**
"All templates already initialized. Would you like to re-customize them?"
- If no: Exit
- If yes: Continue with customization

**If none exist:**
"No templates found. We'll initialize and customize all three."

**If some exist:**
"Found: [list existing]. Missing: [list missing]."
Ask: "Initialize missing templates, or re-initialize all?"

**🤔 Ask yourself before proceeding**:
1. Did I explain what these templates do and how they integrate?
2. Did I check current status with --dry-run?
3. Did I parse the output correctly?
4. Does the user understand what will happen?
5. Has the user confirmed they want to proceed?

---

### 📋 Step 1: Initialize and Customize definition-of-ready.md

**🎯 Goal**: Create definition-of-ready template and customize prerequisite checks for document creation.

#### 1.1: Create Template

If definition-of-ready.md doesn't exist, run:
```bash
{{STARTUP_PATH}}/bin/the-startup init definition-of-ready
```

Then read the created file:
```bash
Read docs/definition-of-ready.md
```

#### 1.2: Identify Customization Points

Scan the definition-of-ready template and identify all `[NEEDS CLARIFICATION: ...]` markers. Common markers:
- `dor threshold` - Overall completion threshold (default: 85%)
- `build command` - Command to verify project builds
- `test command` - Command to run tests
- Language-specific checks

Present findings to user:
"📝 definition-of-ready.md has [N] customization points:
1. DOR threshold: What % of prerequisites must be met? (default: 85%)
2. Build command: How to verify your project builds? (default: go build ./...)
3. [etc.]

Should I help you customize these now, or skip to next template?"

#### 1.3: Customize Based on User Input

**If user wants to customize now:**
For each marker, ask the user for their value and use Edit to replace it.

**If user wants to skip:**
Note in TodoWrite and proceed to Step 2.

**🤔 Ask yourself before proceeding**:
1. Is definition-of-ready.md created?
2. Did I read the entire template?
3. Did I identify all [NEEDS CLARIFICATION] markers?
4. Did I explain what each marker controls?
5. If user chose to customize, did I update all requested markers?
6. Have I clearly communicated what was done?

Ask: "definition-of-ready.md is ready. Should I proceed to definition-of-done.md?"

---

### 📋 Step 2: Initialize and Customize definition-of-done.md

**🎯 Goal**: Create definition-of-done template and customize completion validation for documents.

#### 2.1: Create Template

If definition-of-done.md doesn't exist, run:
```bash
{{STARTUP_PATH}}/bin/the-startup init definition-of-done
```

Then read the created file:
```bash
Read docs/definition-of-done.md
```

#### 2.2: Identify Customization Points

Scan the definition-of-done template and identify all `[NEEDS CLARIFICATION: ...]` markers. Common markers:
- `dod threshold` - Overall completion threshold
- `enable scqa` / `scqa scope` - SCQA logical flow validation settings
- `enable mece` / `mece scope` - MECE coverage validation settings
- `consistency check` - Cross-document validation approach

Present findings to user:
"📝 definition-of-done.md has [N] customization points:
1. DOD threshold: What % of checks must pass? (default: 85%)
2. SCQA validation: Enable logical flow checks? (default: yes, all docs)
3. MECE validation: Enable completeness checks? (default: yes, all docs)
4. [etc.]

Should I help you customize these now, or skip to next template?"

#### 2.3: Customize Based on User Input

**If user wants to customize now:**
For each marker, ask the user for their value and use Edit to replace it.

**If user wants to skip:**
Note in TodoWrite and proceed to Step 3.

**🤔 Ask yourself before proceeding**:
1. Is definition-of-done.md created?
2. Did I read the entire template?
3. Did I identify all [NEEDS CLARIFICATION] markers?
4. Did I explain validation options (SCQA, MECE, consistency)?
5. If user chose to customize, did I update all requested markers?

Ask: "definition-of-done.md is ready. Should I proceed to task-definition-of-done.md?"

---

### 📋 Step 3: Initialize and Customize task-definition-of-done.md

**🎯 Goal**: Create task-definition-of-done template and customize completion validation for implementation tasks.

#### 3.1: Create Template

If task-definition-of-done.md doesn't exist, run:
```bash
{{STARTUP_PATH}}/bin/the-startup init task-definition-of-done
```

Then read the created file:
```bash
Read docs/task-definition-of-done.md
```

#### 3.2: Identify Customization Points

Scan the task-definition-of-done template and identify all `[NEEDS CLARIFICATION: ...]` markers. Common markers:
- `build command` - Command to build project
- `test command` - Command to run tests
- `coverage command` - Command to check coverage
- `coverage target` - Minimum coverage % required
- `lint command` - Command to check code quality
- `format command` - Command to check formatting
- `task dod threshold` - Overall completion threshold

Present findings to user:
"📝 task-definition-of-done.md has [N] customization points:
1. Build command: (e.g., go build ./..., npm run build)
2. Test command: (e.g., go test ./..., npm test)
3. Coverage target: What % coverage required? (default: 80%)
4. Lint command: (e.g., golangci-lint run, eslint .)
5. [etc.]

These are critical - they determine if tasks can be marked complete.
Should I help you customize these now?"

#### 3.3: Customize Based on User Input

**If user wants to customize now:**
For each marker, ask the user for their value and use Edit to replace it.

**If user wants to skip:**
Note in TodoWrite that manual customization is needed.

**🤔 Ask yourself before proceeding**:
1. Is task-definition-of-done.md created?
2. Did I read the entire template?
3. Did I identify all [NEEDS CLARIFICATION] markers?
4. Did I explain what each command does?
5. If user chose to customize, did I update all requested markers?
6. Did I emphasize these commands are critical for /s:implement?

---

### 📋 Step 4: Finalization

**🎯 Goal**: Summarize what was created, note what still needs customization, explain next steps.

Review TodoWrite to see which templates were customized vs. skipped.

**📝 Present Summary**:
```
✅ Quality Gate Templates Initialized

Created:
  • docs/definition-of-ready.md - Definition of Ready
  • docs/definition-of-done.md - Definition of Done
  • docs/task-definition-of-done.md - Task Definition of Done

Customized:
  [List which templates were fully customized]

Needs Manual Customization:
  [List which templates still have [NEEDS CLARIFICATION] markers]

🔍 To find markers needing customization:
  grep -n "NEEDS CLARIFICATION" docs/definition-of-ready.md docs/definition-of-done.md docs/task-definition-of-done.md

📝 Next Steps:
  1. Commit templates to git (so team shares standards)
  2. [If any uncustomized] Manually customize remaining markers
  3. Start using with: /s:specify "your feature description"

🔗 How It Works:
  • /s:specify validates DOR before SDD/PLAN, DOD after PRD/SDD/PLAN
  • /s:implement validates TASK-DOD after each task
  • Quality gates prevent shortcuts and ensure quality
```

**🤔 Verify before finalizing**:
1. Is TodoWrite showing all template steps as completed?
2. Have I clearly listed which templates need manual customization?
3. Did I explain how to find remaining markers?
4. Did I explain next steps (commit, customize, use)?
5. Did I explain how validation integrates with /s:specify and /s:implement?

---

## ⚠️ Error Handling

If `the-startup init` fails:
- Check if `{{STARTUP_PATH}}/bin/the-startup` exists
- Suggest running `the-startup install` if binary missing
- Show error output to user

If templates exist and user didn't use `--force`:
- The command will prompt for confirmation (expected)
- Let user respond to the prompt

## 💡 Remember

- Use TodoWrite to track progress through templates
- Work sequentially: DOR → DOD → TASK-DOD
- Wait for user confirmation between templates
- Only replace [NEEDS CLARIFICATION] markers
- Let `the-startup init` handle file operations
- Parse command output to drive conversation
