package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID       primitive.ObjectID  `json:"_id" bson:"_id,omitempty"` // Unique identifier for the user
	Name     string              `json:"name,omitempty"`          // Name of the user
	Phone    string              `json:"phone,omitempty"`         // Phone number of the user
	Email    string              `json:"email,omitempty"`         // Email address of the user
	Password string              `json:"password,omitempty"`      // Encrypted password of the user
	Latitude float64 			 `json:"latitude,omitempty"`     // Latitude for user's location
	Longitude float64 			 `json:"longitude,omitempty"`   // Longitude for user's location
	Address  string              `json:"address,omitempty"`       // Address of the user
	Rating   float64             `json:"rating,omitempty"`       // User rating
	NumberSelled int			 `json:"numberselled,omitempty"`       // Numberselled 
}

type Food struct{
	ID            primitive.ObjectID `json:"_id" bson:"_id,omitempty"` // Unique identifier for the food item
	ImageUrl      string             `json:"imageurl,omitempty"`       // URL of the food image
	Title         string             `json:"title,omitempty"`          // Title of the food item
	Description   string             `json:"description,omitempty"`    // Description of the food item
	NumberServing int      			 `json:"numberserving,omitempty"`  // Number of servings
	Price         int                `json:"price,omitempty"`          // Price of the food item
	UserID        primitive.ObjectID `json:"userid,omitempty"` // Reference to the user who added the food item
}

type Order struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`          // Unique identifier for the order
	ConsumerID  primitive.ObjectID `json:"consumer_id,omitempty"` // ID of the user who placed the order
	ProducerID  primitive.ObjectID `json:"producer_id,omitempty"` // ID of the user providing the food
	FoodID      primitive.ObjectID `json:"food_id,omitempty"`    // ID of the food which was ordered
	IsRated     bool               `json:"is_rated,omitempty"`       // Whether the order has been rated
	Rating      float64            `json:"rating"`          // Rating given for the order
	Timestamp   time.Time          `json:"timestamp,omitempty"`                        // Time when the order was placed
}

type FoodWithUserInfo struct {
	Food        Food       `json:"food"`          // Food item details
	UserName    string     `json:"user_name"`     // Name of the user providing the food
	UserAddress string     `json:"user_address"`  // Address of the user providing the food
	UserRating  float64    `json:"user_rating"`   // Rating of the user providing the food
	UserLat     float64    `json:"user_latitude"` // Latitude of the user
	UserLon     float64    `json:"user_longitude"`// Longitude of the user
}
