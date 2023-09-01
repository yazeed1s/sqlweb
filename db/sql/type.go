package sql

type DbType int

const (
	MySQL DbType = iota + 1
	PostgreSQL
	SQLite
	Unsupported
)

func (t DbType) String() string {
	if t >= MySQL && t <= SQLite {
		return [...]string{"MySQL", "PostgreSQL", "SQLite"}[t-1]
	}
	return "Unsupported"
}

func (t DbType) EnumIndex() int {
	return int(t)
}
