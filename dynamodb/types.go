package dynamodb

// USER = USER#<username> + PROFILE#<username>
// USER_ADDRESS = NA + NA
// ORDER = USER#<username> + ORDER#<orderid>
// ORDER_ITEM = ITEM#<itemid> + ORDER#<orderid>

type Address struct {
	Street      string
	PostalCode  uint
	State       string
	CountryCode string
}

type User struct {
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"SK"`

	Username  string             `dynamodbav:"Username"`
	Fullname  string             `dynamodbav:"Fullname"`
	Email     string             `dynamodbav:"Email"`
	CreatedAt string             `dynamodbav:"Created_at"`
	Addresses map[string]Address `dynamodbav:"Addresses"`
}

type Order struct {
	PK string
	SK string

	IdOrder   string
	Username  string
	Address   Address
	Status    string
	CreatedAt string
}

type Item struct {
	PK string
	SK string

	IdItem      string
	IdOrder     string
	ProductName string
	Price       float64
	Status      string
}

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

type GenericItem struct {
	PK         string
	SK         string
	Attributes map[string]any
}
