package id

type AccountID string

func (a AccountID) String() string {
	return string(a)
}

type TripID string

func (t TripID) String() string {
	return string(t)
}

type IdentityID string

func (i IdentityID) String() string {
	return string(i)
}
