package repository

// LeadRepository is deprecated
// Use ContactRepository instead with ContactStatus = "Leads"
// This file exists only for backward compatibility

type LeadRepository = ContactRepository

// NewLeadRepository creates a new contact repository (alias for leads)
// Deprecated: Use NewContactRepository instead
func NewLeadRepository(db interface{}) *ContactRepository {
	// This is just an alias pointing to ContactRepository
	return nil
}
