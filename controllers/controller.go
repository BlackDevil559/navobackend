package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/BlackDevil559/novahack2/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const dbName = "novahack2"
const colName1 = "User"
const colName2 = "Food"
const colName3 = "Order"
var collection1_user *mongo.Collection
var collection2_food *mongo.Collection
var collection3_order *mongo.Collection


func init(){
	err:=godotenv.Load("./.env")
    if err!=nil{
        log.Println("Error loading .env file")
    }
    connectionstring := os.Getenv("MONGODB_URI")
	clientOption:=options.Client().ApplyURI(connectionstring)
	client,err:=mongo.Connect(context.TODO(),clientOption)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection Sucess")
	collection1_user=client.Database(dbName).Collection(colName1)
	collection2_food=client.Database(dbName).Collection(colName2)
	collection3_order=client.Database(dbName).Collection(colName3)
	fmt.Println("Collection Intance is Ready")
} 

func addNewUser(user model.User){
	inserted,err:=collection1_user.InsertOne(context.Background(),user)
	if err!=nil{
		panic(err)
	}
	fmt.Println("User is added successfully",inserted.InsertedID)
	go func() {
		subject := "Welcome to HungerPoint!"
		body := fmt.Sprintf("Hello %s,\n\nThank you for joining us. We are excited to have you on board!", user.Name)
		recipientEmail := user.Email
		GeneralMailScript(subject, body, recipientEmail)
	}()
}

func addNewFood(food model.Food){
	inserted,err:=collection2_food.InsertOne(context.Background(),food)
	if err!=nil{
		panic(err)
	}
	fmt.Println("Food is added successfully",inserted.InsertedID)
}

func deleteUser(userId string){
	id,_:=primitive.ObjectIDFromHex(userId)
	filter:=bson.M{"_id":id}
	result,err:=collection1_user.DeleteOne(context.Background(),filter)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Modified count",result.DeletedCount)
}

func deleteFoodItem(foodId string){
	id,_:=primitive.ObjectIDFromHex(foodId)
	filter:=bson.M{"_id":id}
	result,err:=collection2_food.DeleteOne(context.Background(),filter)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Modified count",result.DeletedCount)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius=6371
	dLat:=(lat2-lat1)*math.Pi/180.0
	dLon:=(lon2-lon1)*math.Pi/180.0
	lat1=lat1*math.Pi/180.0
	lat2=lat2*math.Pi/180.0
	a:=math.Sin(dLat/2)*math.Sin(dLat/2)+
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c:=2*math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	return earthRadius*c
}

func getLatLong(userId primitive.ObjectID) (float64, float64, error) {
	var user model.User
	err:=collection1_user.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&user)
	if err!=nil{
		return 0, 0, err
	}
	lat := user.Latitude
	lon := user.Longitude
	return lat, lon, nil
}

func showFoodNearBy(userId string) ([]model.FoodWithUserInfo, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	userLat, userLon, err := getLatLong(id)
	if err != nil {
		return nil, err
	}
	var foodItems []model.Food
	cursor, err := collection2_food.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &foodItems); err != nil {
		return nil, err
	}
	var nearbyFoods []model.FoodWithUserInfo
	for _, food := range foodItems {
		foodLat, foodLon, err := getLatLong(food.UserID)
		if err != nil {
			continue
		}
		distance := haversine(userLat, userLon, foodLat, foodLon)
		if distance <= 5 {
			var user model.User
			err := collection1_user.FindOne(context.TODO(), bson.M{"_id": food.UserID}).Decode(&user)
			if err != nil {
				continue
			}
			nearbyFoods = append(nearbyFoods, model.FoodWithUserInfo{
				Food:        food,
				UserName:    user.Name,
				UserAddress: user.Address,
				UserRating:  user.Rating,
				UserLat:     foodLat,
				UserLon:     foodLon,
			})
		}
	}
	return nearbyFoods, nil
}


func bookFoodItem(foodId string, NumberServing string, consumerId string) {
    foodObjID, err := primitive.ObjectIDFromHex(foodId)
    if err != nil {
        panic(err)
    }
    consumerObjID, err := primitive.ObjectIDFromHex(consumerId)
    if err != nil {
        panic(err)
    }
    var food model.Food
    err = collection2_food.FindOne(context.TODO(), bson.M{"_id": foodObjID}).Decode(&food)
    if err != nil {
        panic(err)
    }
    var producer model.User
    err = collection1_user.FindOne(context.TODO(), bson.M{"_id": food.UserID}).Decode(&producer)
    if err != nil {
        panic(err)
    }
    var consumer model.User
    err = collection1_user.FindOne(context.TODO(), bson.M{"_id": consumerObjID}).Decode(&consumer)
    if err != nil {
        panic(err)
    }
    numServing, err := strconv.Atoi(NumberServing)
    if err != nil {
        panic(err)
    }
    newServings := food.NumberServing - numServing
    if newServings >= 0 {
        if newServings > 0 {
            _, err := collection2_food.UpdateOne(
                context.TODO(),
                bson.M{"_id": foodObjID},
                bson.M{"$set": bson.M{"numberserving": newServings}},
            )
            if err != nil {
                panic(err)
            }
        } else {
            deleteFoodItem(foodId)
        }
        newOrder := model.Order{
            ID:         primitive.NewObjectID(),
            ConsumerID: consumerObjID,
            ProducerID: food.UserID,
            IsRated:    false,
            Rating:     0.0,
            Timestamp:  time.Now(),
        }
        _, err := collection3_order.InsertOne(context.TODO(), newOrder)
        if err != nil {
            panic(err)
        }
        fmt.Println("Order created successfully!")
        go func() {
            subject := "Order Confirmation"
            body := fmt.Sprintf(
                "Hello %s,\n\nThank you for your order! Here are the details:\n\nFood Item: %s\nProvider: %s\nProvider Address: %s\nNumber of Servings: %d\nTotal Price: %d\n\nEnjoy your meal!",
                consumer.Name,
                food.Title,
                producer.Name,
                producer.Address,
                numServing,
                numServing*food.Price,
            )
            recipientEmail := consumer.Email
            err := GeneralMailScript(subject, body, recipientEmail)
            if err != nil {
                fmt.Printf("Failed to send confirmation email: %v\n", err)
            }
        }()
    } else {
        fmt.Println("Not enough servings available.")
    }
}

func loginUser(phone string, password string) (model.User, error) {
	var user model.User
	err:=collection1_user.FindOne(context.TODO(), bson.M{"phone": phone}).Decode(&user)
	if err!=nil{
		return user,err
	}
	if user.Password!=password{
		return user,errors.New("invalid credentials")
	}
	return user,nil
}

func showAllOrder(userId string) ([]model.Order, error) {
    userObjID, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %v", err)
    }
    filter := bson.M{"consumerid": userObjID}
    cursor, err := collection3_order.Find(context.TODO(), filter)
    if err != nil {
        return nil, fmt.Errorf("error fetching orders: %v", err)
    }
    defer cursor.Close(context.TODO())
    var orders []model.Order
    for cursor.Next(context.TODO()) {
        var order model.Order
        if err := cursor.Decode(&order); err != nil {
            return nil, fmt.Errorf("error decoding order: %v", err)
        }
        orders = append(orders, order)
    }
    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %v", err)
    }
    return orders, nil
}

func addRating(orderId string, rating float64) {
	orderObjID, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		panic(fmt.Errorf("invalid order ID: %v", err))
	}
	var order model.Order
	err = collection3_order.FindOne(context.TODO(), bson.M{"_id": orderObjID}).Decode(&order)
	if err != nil {
		panic(fmt.Errorf("order not found: %v", err))
	}
	if order.IsRated {
		panic(fmt.Errorf("order already rated"))
	}
	_, err = collection3_order.UpdateOne(
		context.TODO(),
		bson.M{"_id": orderObjID},
		bson.M{
			"$set": bson.M{
				"rating":   rating,
				"is_rated": true,
			},
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to update order rating: %v", err))
	}
	var producer model.User
	err = collection1_user.FindOne(context.TODO(), bson.M{"_id": order.ProducerID}).Decode(&producer)
	if err != nil {
		panic(fmt.Errorf("producer not found: %v", err))
	}
	newRating := ((producer.Rating * float64(producer.NumberSelled)) + rating) / float64(producer.NumberSelled+1)
	_, err = collection1_user.UpdateOne(
		context.TODO(),
		bson.M{"_id": order.ProducerID},
		bson.M{
			"$set": bson.M{"rating": newRating},
			"$inc": bson.M{"numberselled": 1},
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to update producer rating: %v", err))
	}
	go func() {
		subject := "You've Received a New Rating!"
		body := fmt.Sprintf(
			"Hello %s,\n\nYou have received a new rating for your provided food item.\n\nRating: %.2f\nYour updated average rating: %.2f\n\nThank you for your service!\n",
			producer.Name, rating, newRating,
		)
		err := GeneralMailScript(subject, body, producer.Email)
		if err != nil {
			fmt.Printf("Failed to send email to producer: %v\n", err)
		}
	}()
	fmt.Println("Rating added successfully and notification sent to producer")
}

func showAllFood() ([]model.FoodWithUserInfo, error) {
	var foodItems []model.Food
	cursor, err := collection2_food.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &foodItems); err != nil {
		return nil, err
	}
	var allFoods []model.FoodWithUserInfo
	for _, food := range foodItems {
		var user model.User
		err := collection1_user.FindOne(context.TODO(), bson.M{"_id": food.UserID}).Decode(&user)
		if err != nil {
			continue
		}
		allFoods = append(allFoods, model.FoodWithUserInfo{
			Food:        food,
			UserName:    user.Name,
			UserAddress: user.Address,
			UserRating:  user.Rating,
			UserLat:     user.Latitude,
			UserLon:     user.Longitude,
		})
	}
	return allFoods, nil
}


func AddNewUser(w http.ResponseWriter, r *http.Request){
	fmt.Println("Adding User")
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","POST")
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)
	addNewUser(user)
	json.NewEncoder(w).Encode(user)
}

func AddNewFood(w http.ResponseWriter, r *http.Request){
	fmt.Println("Adding Food")
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","POST")
	var food model.Food
	json.NewDecoder(r.Body).Decode(&food)
	addNewFood(food)
	json.NewEncoder(w).Encode(food)
}

func DeleteUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","DELETE")
	params:=mux.Vars(r)
	deleteUser(params["id"])
	json.NewEncoder(w).Encode("User deleted")
	fmt.Println("User deleted")
}

func DeleteFoodItem(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","DELETE")
	params:=mux.Vars(r)
	deleteFoodItem(params["id"])
	json.NewEncoder(w).Encode("Fooditem deleted")
	fmt.Println("Fooditem deleted")
}

func ShowFoodNearBy(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","GET")
	params:=mux.Vars(r)
	food,err:=showFoodNearBy(params["id"])
	if err!=nil{
		panic(err)
	}
	json.NewEncoder(w).Encode(food)
	fmt.Println("Fooditem Displayed")
}

func BookFooditem(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","GET")
	params:=mux.Vars(r)
	bookFoodItem(params["id"],params["ns"],params["consumerid"])
	json.NewEncoder(w).Encode("Order Confirmed")
	fmt.Println("Order Confirmed")
}

func LoginUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","GET")
	params:=mux.Vars(r)
	user,err:=loginUser(params["phone"],params["password"])
	if err!=nil{
		panic(err)
	}
	json.NewEncoder(w).Encode(user)
	fmt.Println("Logged In")
}

func ShowAllOrder(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","GET")
	params:=mux.Vars(r)
	orders,err:=showAllOrder(params["id"])
	if err!=nil{
		panic(err)
	}
	json.NewEncoder(w).Encode(orders)
	fmt.Println("All orders sent")
}

func AddRating(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","PUT")
	params:=mux.Vars(r)
	floatRating, err := strconv.ParseFloat(params["rating"], 64)
	addRating(params["id"],floatRating)
	if err!=nil{
		panic(err)
	}
	json.NewEncoder(w).Encode("Rating Done")
	fmt.Println("Rating Done")
}

func ShowAllFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	w.Header().Set("User-Agent","CustomUserAgent/1.0")
	if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
	foodWithUserInfo, err := showAllFood()
	if err != nil {
		http.Error(w, "Error fetching food items", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(foodWithUserInfo)
	fmt.Println("All food items with user info displayed")
}

func GeneralMailScript(subject, body, recipientEmail string) error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", email, password, smtpServer)
	msg := fmt.Sprintf("Subject: %s\n\n%s", subject, body)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpServer, smtpPort),
		auth,
		email,
		[]string{recipientEmail},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}