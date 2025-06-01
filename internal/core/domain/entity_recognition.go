package domain

// EntityType represents the type of a named entity
type EntityType string

// Common entity types
const (
	EntityTypePerson       EntityType = "person"
	EntityTypeOrganization EntityType = "organization"
	EntityTypeLocation     EntityType = "location"
	EntityTypeDate         EntityType = "date"
	EntityTypeTime         EntityType = "time"
	EntityTypeEmail        EntityType = "email"
	EntityTypeURL          EntityType = "url"
	EntityTypeProduct      EntityType = "product"
	EntityTypeCurrency     EntityType = "currency"
	EntityTypePercentage   EntityType = "percentage"
	EntityTypePhoneNumber  EntityType = "phone_number"
	EntityTypeIPAddress    EntityType = "ip_address"
	EntityTypeOther        EntityType = "other"
)
