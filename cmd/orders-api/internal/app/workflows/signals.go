package workflows

// SignalChannels are all available channels used by our workflow
var SignalChannels = struct {
	CONFIRM_ORDER_CHANNEL     string
	CANCEL_ORDER_CHANNEL      string
	EXPIRE_ORDER_CHANNEL      string
	PAYMENT_PROCESSED_CHANNEL string
}{
	CONFIRM_ORDER_CHANNEL:     "CONFIRM_ORDER_CHANNEL",
	CANCEL_ORDER_CHANNEL:      "CANCEL_ORDER_CHANNEL",
	EXPIRE_ORDER_CHANNEL:      "EXPIRE_ORDER_CHANNEL",
	PAYMENT_PROCESSED_CHANNEL: "PAYMENT_PROCESSED_CHANNEL",
}

// ConfirmOrderSignal is the body that is sent through the  ConfirmOrder channel
type ConfirmOrderSignal struct {
	OrderID         string
	PaymentOptionID string
}
