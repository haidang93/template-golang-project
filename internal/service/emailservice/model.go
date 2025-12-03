package emailservice

type SendData struct {
	TemplateUUID      *string
	From              *Email
	To                *[]Email
	Cc                *[]Email
	Bcc               *[]Email
	TemplateVariables *map[string]interface{}
}

type Email struct {
	Email *string
	Name  *string
}

func (s *SendData) Data() map[string]interface{} {
	from := map[string]*string{}
	to := []map[string]*string{}
	cc := []map[string]*string{}
	bcc := []map[string]*string{}

	if s.From != nil {
		from["email"] = s.From.Email
		from["name"] = s.From.Name
	}

	if s.To != nil {
		for _, address := range *s.To {
			to = append(to, map[string]*string{
				"email": address.Email,
				"name":  address.Name,
			})
		}
	}

	if s.Cc != nil {
		for _, address := range *s.Cc {
			cc = append(cc, map[string]*string{
				"email": address.Email,
				"name":  address.Name,
			})
		}
	}

	if s.Bcc != nil {
		for _, address := range *s.Bcc {
			bcc = append(bcc, map[string]*string{
				"email": address.Email,
				"name":  address.Name,
			})
		}
	}

	return map[string]interface{}{
		"from":               from,
		"to":                 to,
		"cc":                 cc,
		"bcc":                bcc,
		"template_uuid":      s.TemplateUUID,
		"template_variables": s.TemplateVariables,
	}
}
