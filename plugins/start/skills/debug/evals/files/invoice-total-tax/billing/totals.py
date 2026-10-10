from dataclasses import dataclass
from decimal import Decimal, ROUND_HALF_UP


@dataclass
class LineItem:
    description: str
    unit_price: Decimal
    quantity: int = 1


@dataclass
class Invoice:
    items: list
    tax_rate: Decimal = Decimal("0.10")
    discount: Decimal = Decimal("0")


def round_currency(amount):
    return amount.quantize(Decimal("0.01"), rounding=ROUND_HALF_UP)


def line_total(item):
    return item.unit_price * item.quantity


def subtotal(invoice):
    return sum((line_total(item) for item in invoice.items), Decimal("0"))


def tax_for(amount, rate):
    return round_currency(amount * rate)


def apply_discount(amount, discount):
    return max(amount - discount, Decimal("0"))


def invoice_total(invoice):
    net = subtotal(invoice)
    total = net + tax_for(net, invoice.tax_rate)
    total = apply_discount(net, invoice.discount)
    return float(round_currency(total))
