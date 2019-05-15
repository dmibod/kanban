package kernel

// ID type
type ID string

// EmptyID value
var EmptyID = ID("")

// IsValid id
func (id ID) IsValid() bool {
	return id != EmptyID
}

// String converts ID to string
func (id ID) String() string {
	return string(id)
}

// WithSet builds MemberID
func (id ID) WithSet(setID ID) MemberID {
	return MemberID{ID: id, SetID: setID}
}

// WithID builds MemberID
func (setID ID) WithID(id ID) MemberID {
	return MemberID{ID: id, SetID: setID}
}

// MemberID represents complex ID
type MemberID struct {
	SetID ID
	ID    ID
}
