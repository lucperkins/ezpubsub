package main

type DB struct {
	data []string
}

func NewDB() *DB {
	d := make([]string, 0)
	return &DB{
		data: d,
	}
}

func (d *DB) save(s string) {
	d.data = append(d.data, s)
}
