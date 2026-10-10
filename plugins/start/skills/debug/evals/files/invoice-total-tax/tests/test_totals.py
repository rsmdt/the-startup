from decimal import Decimal

from billing.totals import Invoice, LineItem, invoice_total


def test_invoice_total():
    invoice = Invoice(items=[LineItem("Consulting hour", Decimal("100.00"))], tax_rate=Decimal("0.10"))
    assert invoice_total(invoice) == 110.00
