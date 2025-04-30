package dynamodb

type UserDto struct {
	Username string
	Fullname string
	Email    string
}

type AddressDto struct {
	Username string
	Label    string
	Address  Address
}

type OrderDto struct {
	Username string
	OrderId  string
	Status   string
	Address  Address
	Items    []OrderItemDto
}

type OrderItemDto struct {
	IdItem      string
	ProductName string
	Price       float64
	Status      string
}
