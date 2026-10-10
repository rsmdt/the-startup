# The Agentic Startup — task targets.

.PHONY: eval eval-all eval-list

# Run a skill's eval set via the skill-creator workflow (with_skill only by default).
#   make eval SKILL=specify-plan
#   make eval SKILL=debug RUNS=3 BASELINE=1
eval:
	@test -n "$(SKILL)" || { echo "usage: make eval SKILL=<skill> [RUNS=N] [MODEL=name] [BASELINE=1]"; exit 2; }
	@scripts/eval-skill.sh "$(SKILL)" $(if $(RUNS),--runs $(RUNS),) $(if $(MODEL),--model $(MODEL),) $(if $(BASELINE),--baseline,)

# Run every eligible skill's eval set, JOBS at a time (default 3); logs go to EVAL_LOGS.
#   make eval-all
#   make eval-all JOBS=4 BASELINE=1
JOBS ?= 3
EVAL_LOGS ?= .eval-logs
eval-all:
	@mkdir -p $(EVAL_LOGS)
	@scripts/eval-skill.sh --list | awk 'NR>1 && !/excluded/ {print $$1}' | \
		xargs -P $(JOBS) -I{} sh -c 'scripts/eval-skill.sh {} $(if $(RUNS),--runs $(RUNS),) $(if $(MODEL),--model $(MODEL),) $(if $(BASELINE),--baseline,) > $(EVAL_LOGS)/{}.log 2>&1; echo "{} exit=$$?"'

# List skills that have an eval set.
eval-list:
	@scripts/eval-skill.sh --list
