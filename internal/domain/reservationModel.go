package domain

type ReservationStatus string

const (
	RESERVATION_PENDING                 ReservationStatus = "PENDING"
	RESERVATION_PENDING_ASSIGNMENT      ReservationStatus = "PENDING_ASSIGNMENT"
	RESERVATION_GUIDE_PREASSIGNED       ReservationStatus = "GUIDE_PREASSIGNED"
	RESERVATION_GUIDE_CONFIRMED         ReservationStatus = "GUIDE_CONFIRMED"
	RESERVATION_PAYMENT_PENDING         ReservationStatus = "PAYMENT_PENDING"
	RESERVATION_PARTIALLY_CONFIRMED     ReservationStatus = "PARTIALLY_CONFIRMED"
	RESERVATION_CONFIRMED               ReservationStatus = "CONFIRMED"
	RESERVATION_COMPLETED               ReservationStatus = "COMPLETED"
	RESERVATION_CANCELLED               ReservationStatus = "CANCELLED"
	RESERVATION_FORCE_MAJEURE_CANCELLED ReservationStatus = "FORCE_MAJEURE_CANCELLED"
)

type SignatureStatus string

const (
	SIGNATURE_NOT_SENT SignatureStatus = "NOT_SENT"
	SIGNATURE_SENT     SignatureStatus = "SENT"
	SIGNATURE_SIGNED   SignatureStatus = "SIGNED"
	SIGNATURE_REJECTED SignatureStatus = "REJECTED"
)

type VoucherStatus string

const (
	VOUCHER_NOT_GENERATED  VoucherStatus = "NOT_GENERATED"
	VOUCHER_GENERATED      VoucherStatus = "GENERATED"
	VOUCHER_PARTIALLY_SENT VoucherStatus = "PARTIALLY_SENT"
	VOUCHER_SENT           VoucherStatus = "SENT"
)
