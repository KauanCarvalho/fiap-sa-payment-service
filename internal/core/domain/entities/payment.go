package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

const PaymentStatusPending = "pending"

type Payment struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Amount            float64            `bson:"amount" json:"amount"`
	Status            string             `bson:"status" json:"status"`
	ExternalReference string             `bson:"external_reference" json:"external_reference"`
	Provider          string             `bson:"provider" json:"provider"`
	PaymentMethod     string             `bson:"payment_method" json:"payment_method"`
	QRCode            string             `bson:"qr_code" json:"qr_code"`
}
