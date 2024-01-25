package companydto

type CompanyBySchemaName struct {
	SchemaName string `json:"schema_name"`
}

func (c *CompanyBySchemaName) validate() error {
	if c.SchemaName == "" {
		return ErrSchemaNameIsEmpty
	}

	return nil
}

func (c *CompanyBySchemaName) ToModel() (string, error) {
	if err := c.validate(); err != nil {
		return "", err
	}

	return c.SchemaName, nil
}
