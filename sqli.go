package injecter

type sqli struct {
}

func ConnectSQLI(address string) (Handle, error) {
	return &sqli{}, nil
}

func (d *sqli) Query(query string) ([]map[string]string, error) {
	return nil, nil
}

func (d *sqli) Exec(query string) error {
	return nil
}
