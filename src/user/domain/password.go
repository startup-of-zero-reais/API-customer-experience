package domain

type Password string

func (p Password) IsValid() bool {
	return len(p) > 0
}

func (p Password) Hash() Password {
	return p
}
