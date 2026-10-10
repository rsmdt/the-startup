export interface Mailer {
  send(to: string, subject: string, body: string): Promise<void>;
}

export class NotificationService {
  constructor(private readonly mailer: Mailer) {}

  async paymentReceipt(email: string, amountCents: number, currency: string): Promise<void> {
    const amount = (amountCents / 100).toFixed(2);
    await this.mailer.send(email, "Payment received", `We received ${amount} ${currency}.`);
  }

  async welcome(email: string, displayName: string): Promise<void> {
    await this.mailer.send(email, "Welcome", `Hi ${displayName}, your account is ready.`);
  }
}
