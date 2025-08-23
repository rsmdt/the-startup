Before automating infrastructure, ask yourself:
- Do I do this manually more than twice? (If yes, automate it)
- What breaks when I'm not watching? (Automate monitoring first)
- How do I roll back if this deployment fails? (Build rollback before deploy)
- Can this wait until tomorrow morning? (If no, automate the fix)

Before deploying, ask yourself:
- Have I tested this deployment process in a staging environment?
- Do I have metrics showing the current system is healthy?
- Is the rollback plan faster than debugging forward?
- Will users notice if something breaks? (If yes, use blue-green deployment)

Before scaling infrastructure, ask yourself:
- What's the actual bottleneck? (Measure, don't guess)
- Am I solving for current load or hypothetical future load?
- What's the cheapest fix that gets me 6 months of runway?
- Is this complexity worth the performance gain?

Before adding monitoring, ask yourself:
- What decision will this metric help me make?
- Will this alert wake me up for something I can fix?
- Am I monitoring symptoms (user impact) or causes (system metrics)?
- How will I know when the system is healthy vs. struggling?

Security checkpoints:
- Are secrets in environment variables, not code?
