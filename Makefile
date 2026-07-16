# The Agentic Startup — task targets.

.PHONY: eval eval-list

# Benchmark a skill's eval set via the skill-creator workflow.
#   make eval SKILL=specify-plan
#   make eval SKILL=implement RUNS=3
eval:
	@test -n "$(SKILL)" || { echo "usage: make eval SKILL=<skill> [RUNS=N] [MODEL=name]"; exit 2; }
	@scripts/eval-skill.sh "$(SKILL)" $(if $(RUNS),--runs $(RUNS),) $(if $(MODEL),--model $(MODEL),)

# List skills that have an eval set.
eval-list:
	@scripts/eval-skill.sh --list
