package injecter

type phpmyadmin struct {
}

func ConnectPHPMyAdmin(address, username, password string) (Handle, error) {
	return &phpmyadmin{}, nil
}

func (d *phpmyadmin) Query(query string) ([]map[string]string, error) {
	return nil, nil
}

func (d *phpmyadmin) Exec(query string) error {
	return nil
}
