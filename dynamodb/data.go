package dynamodb

var Users = []UserDto{{
	Username: "alexdebrie",
	Fullname: "Alex DeBrie",
	Email:    "alexdebrie1@gmail.com",
}, {
	Username: "a1nzdev",
	Fullname: "Brawer Nuñez",
	Email:    "a1nzdev28@gmail.com",
}, {
	Username: "johndoe",
	Fullname: "John Doe",
	Email:    "jhondoe1@hotmail.com",
}}

var Addresses = []AddressDto{{
	Username: "alexdebrie",
	Label:    "Work",
	Address: Address{
		Street:      "123 Main St",
		PostalCode:  90211,
		State:       "CA",
		CountryCode: "USA",
	},
}, {
	Username: "a1nzdev",
	Label:    "Home",
	Address: Address{
		Street:      "421 Main St",
		PostalCode:  70201,
		State:       "PA",
		CountryCode: "USA",
	},
}, {
	Username: "a1nzdev",
	Label:    "Work",
	Address: Address{
		Street:      "942 Main St",
		PostalCode:  80100,
		State:       "PA",
		CountryCode: "USA",
	},
}, {
	Username: "johndoe",
	Label:    "School",
	Address: Address{
		Street:      "555 Main St",
		PostalCode:  20201,
		State:       "AS",
		CountryCode: "USA",
	},
}}

var Orders = []OrderDto{
	{
		Username: "alexdebrie",
		OrderId:  "a3f2c1b7",
		Status:   "PLACED",
		Address: Address{
			Street:      "123 Main St",
			PostalCode:  90211,
			State:       "CA",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "b8d4e2a1",
				ProductName: "Macbook Pro",
				Price:       1399.99,
				Status:      "FILLED",
			},
			{
				IdItem:      "c7f3a8b2",
				ProductName: "Amazon Echo",
				Price:       69.99,
				Status:      "FILLED",
			},
		},
	},
	{
		Username: "a1nzdev",
		OrderId:  "e8b4d9a2",
		Status:   "SHIPPED",
		Address: Address{
			Street:      "421 Main St",
			PostalCode:  70201,
			State:       "PA",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "9e2d7c5a",
				ProductName: "Mechanical Keyboard",
				Price:       89.90,
				Status:      "FILLED",
			},
			{
				IdItem:      "f1a3b7c8",
				ProductName: "Gaming Mouse",
				Price:       59.90,
				Status:      "SHIPPED",
			},
		},
	},
	{
		Username: "johndoe",
		OrderId:  "4d7e2f9c",
		Status:   "DELIVERED",
		Address: Address{
			Street:      "555 Main St",
			PostalCode:  20201,
			State:       "AS",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "7a4b3d2c",
				ProductName: "Noise Cancelling Headphones",
				Price:       299.99,
				Status:      "DELIVERED",
			},
		},
	},
}
