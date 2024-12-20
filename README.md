# Food Sharing Backend API
This repository contains the backend code for a food-sharing platform built using Go. The API allows users to add, view, and manage food items, place orders, and handle ratings. MongoDB is used for data storage.

## Project Structure
```
Backend/
├── controllers/          # API handler functions
│   └── controller.go
├── images/               # Directory for images
├── models/               # Data models
│   └── model.go
├── routers/              # API routes
│   └── router.go
├── .env                  # Environment variables
├── .gitignore            # Ignored files
├── dummydata.txt         # Test data
├── go.mod                # Go modules
├── go.sum                # Dependency checksums
├── main.exe              # Binary executable
├── main.go               # Entry point
└── README.md             # Documentation
```

## Getting Started
### Prerequisites
Before running the project, ensure you have the following installed:

- Go (v1.18+)
- MongoDB (Database)
- Gorilla Mux (Routing library)

### Installation
Clone the repository:
```
git clone https://github.com/yourusername/food-sharing-backend.git
cd food-sharing-backend
````
Install dependencies:
```
go mod tidy
```
Set up the .env file:

```
# MongoDB connection string
MONGODB_URI=mongodb+srv://<username>:<password>@cluster0.mvhrh.mongodb.net/

# Email configuration for SMTP (used for notifications)
EMAIL=<mailid>
EMAIL_PASSWORD=<app password>

# SMTP Server settings
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
```
Run the application:
```
go run main.go
```
The server will start on http://localhost:4000.

## API Endpoints
![image](https://github.com/user-attachments/assets/7dcad09c-20c9-4b43-a7f9-8edb2da55e42)
![image](https://github.com/user-attachments/assets/9171657a-f213-4ae0-a14b-05c321c86615)

## Models
The API uses the following data models:
### User Model
```
type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Phone        string             `json:"phone,omitempty"`
	Email        string             `json:"email,omitempty"`
	Password     string             `json:"password,omitempty"`
	Latitude     float64            `json:"latitude,omitempty"`
	Longitude    float64            `json:"longitude,omitempty"`
	Address      string             `json:"address,omitempty"`
	Rating       float64            `json:"rating,omitempty"`
	NumberSelled int                `json:"numberselled,omitempty"`
}
```
### Food Model
```
type Food struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImageUrl      string             `json:"imageurl,omitempty"`
	Title         string             `json:"title,omitempty"`
	Description   string             `json:"description,omitempty"`
	NumberServing int                `json:"numberserving,omitempty"`
	Price         int                `json:"price,omitempty"`
	UserID        primitive.ObjectID `json:"userid,omitempty"`
}
```
### Order Model
```
type Order struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ConsumerID  primitive.ObjectID `json:"consumer_id,omitempty"`
	ProducerID  primitive.ObjectID `json:"producer_id,omitempty"`
	FoodID      primitive.ObjectID `json:"food_id,omitempty"`
	IsRated     bool               `json:"is_rated,omitempty"`
	Rating      float64            `json:"rating,omitempty"`
	Timestamp   time.Time          `json:"timestamp,omitempty"`
}
```
## Contact
For questions or issues, contact:
Bhumesh Gaur

Email: gaurbhumesh559@gmail.com
