export type Currency = "EUR" | "USD";

export interface PaymentResult {
  paymentId: string;
  status: "succeeded" | "requires_action" | "failed";
}

export interface PaymentGateway {
  createIntent(amountCents: number, currency: Currency, idempotencyKey: string): Promise<PaymentResult>;
  refund(paymentId: string): Promise<void>;
}

export class PaymentService {
  constructor(private readonly gateway: PaymentGateway) {}

  async charge(orderId: string, amountCents: number, currency: Currency): Promise<PaymentResult> {
    if (!Number.isInteger(amountCents) || amountCents <= 0) {
      throw new RangeError("amountCents must be a positive integer");
    }
    return this.gateway.createIntent(amountCents, currency, `order-${orderId}`);
  }

  async refund(paymentId: string): Promise<void> {
    await this.gateway.refund(paymentId);
  }
}
