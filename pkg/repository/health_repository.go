package repository

func (r RepositoryImpl) HealthCheck() error {
	db, err := r.DB.DB()
	if err != nil || db.Ping() != nil {
		return err
	}

	return nil
}
