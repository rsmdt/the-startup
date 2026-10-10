# PaymentService

Charges customers through the payment gateway.

## API

### `charge(customerId: string, amount: number): Promise<string>`

Charges `amount` (in euros, decimals allowed) to the customer's saved card and returns the charge ID. Retries up to three times on gateway errors.

### `refund(chargeId: string, amount?: number): Promise<void>`

Refunds the charge. Pass `amount` for a partial refund.

## Notes

- Only EUR is supported.
