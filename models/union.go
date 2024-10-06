package models

// 联合型，用来解决Admin和Student的联合查询(还没用到呢)

type User interface {
	UID() string
	UserName() string
}

func (s *Student) UID() string {
	return s.UserID
}

func (s *Student) UserName() string {
	return s.Username
}

func (a *Admin) UID() string {
	return a.UserID
}

func (a *Admin) UserName() string {
	return a.Username
}
