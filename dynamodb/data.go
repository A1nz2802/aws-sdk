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
	{
		Username: "a1nzdev",
		OrderId:  "b2c5e1f8",
		Status:   "PENDING",
		Address: Address{
			Street:      "123 Oak Ave",
			PostalCode:  10001,
			State:       "NY",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "a9d1e4f2",
				ProductName: "Wireless Headphones",
				Price:       129.99,
				Status:      "PROCESSING",
			},
		},
	},
	{
		Username: "a1nzdev",
		OrderId:  "c3d6f2a9",
		Status:   "DELIVERED",
		Address: Address{
			Street:      "789 Pine Ln",
			PostalCode:  90210,
			State:       "CA",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "b8e2f5a1",
				ProductName: "Smart Speaker",
				Price:       49.99,
				Status:      "DELIVERED",
			},
			{
				IdItem:      "c7f3a6b2",
				ProductName: "USB Hub",
				Price:       19.95,
				Status:      "DELIVERED",
			},
		},
	},
	{
		Username: "a1nzdev",
		OrderId:  "d4e7a3b6",
		Status:   "PROCESSING",
		Address: Address{
			Street:      "456 Elm St",
			PostalCode:  60601,
			State:       "IL",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "e6a4b7c9",
				ProductName: "Webcam",
				Price:       79.00,
				Status:      "PROCESSING",
			},
		},
	},
	{
		Username: "a1nzdev",
		OrderId:  "e5f8b4c7",
		Status:   "SHIPPED",
		Address: Address{
			Street:      "987 Birch Ct",
			PostalCode:  30303,
			State:       "GA",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "f4b8c1d2",
				ProductName: "Ergonomic Chair",
				Price:       299.50,
				Status:      "SHIPPED",
			},
		},
	},
	{
		Username: "a1nzdev",
		OrderId:  "f6a9c5d2",
		Status:   "PENDING",
		Address: Address{
			Street:      "234 Willow Dr",
			PostalCode:  77001,
			State:       "TX",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "a2c6d9e3",
				ProductName: "Laptop Stand",
				Price:       35.00,
				Status:      "PENDING",
			},
			{
				IdItem:      "b1d7e0f4",
				ProductName: "Monitor Arm",
				Price:       99.99,
				Status:      "PENDING",
			},
		},
	},
	{
		Username: "a1nzdev",
		OrderId:  "a7b0d3e5",
		Status:   "DELIVERED",
		Address: Address{
			Street:      "567 Maple Ave",
			PostalCode:  19104,
			State:       "PA",
			CountryCode: "USA",
		},
		Items: []OrderItemDto{
			{
				IdItem:      "c8e1f4a5",
				ProductName: "Desk Lamp",
				Price:       45.75,
				Status:      "DELIVERED",
			},
		},
	},
}
