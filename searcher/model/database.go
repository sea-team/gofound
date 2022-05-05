package model

type Database struct {
	Name string
}

func (d *Database) GetName() string {
	if d.Name == "" {
		return "default"
	}
	return d.Name
}
