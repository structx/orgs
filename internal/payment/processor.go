package payment

// Processor is a payment processor
type Processor interface {
	CreateAccountHolder() (string, error)
}
