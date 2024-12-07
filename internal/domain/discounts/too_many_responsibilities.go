package discounts

type CheckoutUseCase struct {
	cartValidator    CartValidator
	paymentProcessor PaymentProcessor
	receiptSender    ReceiptSender
}

func (c *CheckoutUseCase) Checkout(user User, cart Cart) {
	c.cartValidator.Validate(cart)
	c.paymentProcessor.Process(user, cart)
	c.receiptSender.Send(user)
}

type CartValidator struct{}
type ReceiptSender struct{}
type PaymentProcessor struct{}

func (c *CartValidator) Validate(cart Cart) {
	// Logic about which carts are valid
}
func (p *PaymentProcessor) Process(user User, cart Cart) {
	// Logic about processing payments
}
func (r *ReceiptSender) Send(user User) {
	// Logic about sending receipts
}

type User struct{}
type Cart struct{}
