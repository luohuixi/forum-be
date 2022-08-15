package dao

func (d *Dao) Enforce(rvals ...interface{}) (bool, error) {
	return true, nil
	// return d.CB.Enforce(rvals) //TODO
}
