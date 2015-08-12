package loadbalancer

type VIP struct {
	Name             string
	FrontendIP       string
	FrontendPort     uint8
	FrontendProtocol string
	PoolMembers      []PoolMember
}

type PoolMember struct {
	IP       string
	Port     uint8
	Protocol string
	// HealthCheck HealthCheck  - Not Yet Implemented
}

type LoadBalancer interface {
	CreateVIP(*VIP) error
	UpdateVIP(*VIP) error
	DeleteVIP(string) error
	GetVIP(string) (*VIP, error)
	GetAllVIPs() ([]*VIP, error)
	AddPoolMember(*PoolMember) error
	DeletePoolMember(string) error
	DeleteAllPoolMembers() error
	GetPoolMembers() ([]*PoolMember, error)
}
