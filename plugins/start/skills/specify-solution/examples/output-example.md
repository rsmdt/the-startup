# Example SDD Status Output

🏗️ SDD Status: 003-notification-system

Architecture:
- Pattern: Event-driven microservice
- Key Components: NotificationService, ChannelAdapter, TemplateEngine
- External Integrations: Stripe Webhooks, SendGrid, Twilio

Sections Completed:
- Approach: ✅ Complete
- Components: ✅ Complete
- Interfaces & data: ⚠️ Needs user decision on webhook retry strategy
- Behavior: 🔄 In progress

ADRs:
- ADR-1: Use event sourcing for notification state: ✅ Confirmed
- ADR-2: SendGrid vs Postmark for email: ⏳ Pending confirmation

Validation Status:
- 15 items passed
- 4 items pending

Next Steps:
- Resolve ADR-2 (email provider selection)
