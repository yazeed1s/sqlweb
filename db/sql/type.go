package sql

// DbType represents an enumerated type for various database types.
type DbType int

// Constants representing different database types.
const (
	MySQL DbType = iota + 1
	PostgreSQL
	SQLite
	Unsupported
)

// String returns the string representation of the DbType.
// It converts the DbType constant to its corresponding string value.
// If the DbType is not recognized, it returns "Unsupported".
func (t DbType) String() string {
	if t >= MySQL && t <= SQLite {
		return [...]string{"MySQL", "PostgreSQL", "SQLite"}[t-1]
	}
	return "Unsupported"
}

// EnumIndex returns the integer index of the DbType.
// The index corresponds to the position of the DbType constant in the iota sequence.
// For example, MySQL has an index of 1, PostgreSQL has an index of 2, and so on.
func (t DbType) EnumIndex() int {
	return int(t)
}
