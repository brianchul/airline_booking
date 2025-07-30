package queue

// BookingRequest represents a flight booking request
type BookingRequest struct {
	BookingUUID string             `json:"booking_uuid"`
	UserEmail   string             `json:"user_email"`
	UserTier    string             `json:"user_tier"`
	FlightID    string             `json:"flight_id"`
	Passengers  []PassengerDetails `json:"passengers"`
	PaymentInfo PaymentDetails     `json:"payment_info"`
	ContactInfo ContactDetails     `json:"contact_info"`
	SpecialReqs []string           `json:"special_requirements,omitempty"`
	ClientReqID string             `json:"client_request_id"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
}

// PassengerDetails contains passenger information
type PassengerDetails struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	PassportNo  string `json:"passport_no,omitempty"`
	SeatPref    string `json:"seat_preference,omitempty"`
}

// PaymentDetails contains payment information
type PaymentDetails struct {
	Method    string  `json:"method"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	CardLast4 string  `json:"card_last_4,omitempty"`
	PaymentID string  `json:"payment_id,omitempty"`
}

// ContactDetails contains contact information
type ContactDetails struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}
