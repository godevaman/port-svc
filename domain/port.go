package domain

// Port represents the domain model for a port.
type Port struct {
	ID          string json:"id"
	Name        string json:"name"
	City        string json:"city"
	Country     string json:"country"
	Coordinates []float64 json:"coordinates"
}

// Validate ensures the port data is valid.
func (p Port) Validate() error {
	if p.ID == "" {
		return ErrInvalidID
	}
	return nil
}

// ErrInvalidID is returned when a port ID is invalid.
var ErrInvalidID = errors.New("port ID cannot be empty")