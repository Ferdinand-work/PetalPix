package services

import (
	"context"
	"errors"

	"github.com/Ferdinand-work/PetalPix/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) *UserServiceImpl {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) (*mongo.InsertOneResult, error) {
	res, err := u.userCollection.InsertOne(u.ctx, user)
	return res, err
}

func (u *UserServiceImpl) GetUser(id *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "user_id", Value: id}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.userCollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	defer cursor.Close(u.ctx)

	if len(users) == 0 {
		return nil, errors.New("documents not found")
	}

	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "user_id", Value: user.UserId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "user_name", Value: user.Name},
		bson.E{Key: "user_contact_no", Value: user.ContactNo},
		bson.E{Key: "user_email", Value: user.Email}}}}
	result, _ := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(id *string) error {
	filter := bson.D{bson.E{Key: "user_id", Value: id}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) Follow(usernames interface{}, userId string) (*[]string, error) {

	filter := bson.D{bson.E{Key: "user_id", Value: userId}}
	var update primitive.D
	switch v := usernames.(type) {
	case string:
		update = bson.D{
			bson.E{Key: "$push",
				Value: bson.D{bson.E{Key: "user_following",
					Value: v}}},
			bson.E{Key: "$inc",
				Value: bson.D{bson.E{Key: "user_following_count",
					Value: 1}}},
		}
	case []string:
		update = bson.D{bson.E{Key: "$push",
			Value: bson.D{bson.E{Key: "user_following",
				Value: bson.D{bson.E{Key: "$each",
					Value: v}}}}},
			bson.E{Key: "$inc",
				Value: bson.D{bson.E{Key: "user_following_count",
					Value: len(v)}}},
		}
	default:
		return nil, errors.New("invalid input type")
	}

	result, err := u.userCollection.UpdateOne(u.ctx, filter, update)
	if err != nil {
		return nil, errors.New("cannot update")
	}
	if result.MatchedCount < 1 {
		return nil, errors.New("no matched document found for update")
	}
	var following []string
	switch v := usernames.(type) {
	case string:
		following = append(following, v)
	case []string:
		following = v
	default:
		return nil, errors.New("invalid input type")
	}
	filter1 := bson.M{"user_id": bson.M{"$in": &following}}
	update = bson.D{
		bson.E{Key: "$push",
			Value: bson.D{bson.E{Key: "user_followers",
				Value: userId}}},
		bson.E{Key: "$inc",
			Value: bson.D{bson.E{Key: "user_followers_count",
				Value: 1}}},
	}
	_, err = u.userCollection.UpdateMany(context.TODO(), filter1, update)
	if err != nil {
		return nil, err
	}
	return &following, nil
}

func (u *UserServiceImpl) GetFollowing(id string) (*[]models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "user_id", Value: id}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	if err != nil {
		return nil, err
	}
	follwings := user.Following
	filter := bson.M{"user_id": bson.M{"$in": &follwings}}
	cursor, err := u.userCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	var pointer = &users
	if err = cursor.All(context.Background(), pointer); err != nil {
		return nil, err
	}
	return pointer, nil
}
func (u *UserServiceImpl) Unfollow(usernames interface{}, userId string) (*[]string, error) {
	filter := bson.D{bson.E{Key: "user_id", Value: userId}}
	var update primitive.D
	switch v := usernames.(type) {
	case string:
		update = bson.D{
			bson.E{Key: "$pull",
				Value: bson.D{bson.E{Key: "user_following",
					Value: v}}},
			bson.E{Key: "$inc",
				Value: bson.D{bson.E{Key: "user_following_count",
					Value: -1}}},
		}
	case []string:
		update = bson.D{bson.E{Key: "$pull",
			Value: bson.D{bson.E{Key: "user_following",
				Value: bson.D{bson.E{Key: "$each",
					Value: v}}}}},
			bson.E{Key: "$inc",
				Value: bson.D{bson.E{Key: "user_following_count",
					Value: -len(v)}}},
		}
	default:
		return nil, errors.New("invalid input type")
	}

	result, err := u.userCollection.UpdateOne(u.ctx, filter, update)
	if err != nil {
		return nil, errors.New("cannot update")
	}
	if result.MatchedCount < 1 {
		return nil, errors.New("no matched document found for update")
	}
	var unfollowing []string
	switch v := usernames.(type) {
	case string:
		unfollowing = append(unfollowing, v)
	case []string:
		unfollowing = v
	default:
		return nil, errors.New("invalid input type")
	}
	filter1 := bson.M{"user_id": bson.M{"$in": &unfollowing}}
	update = bson.D{
		bson.E{Key: "$pull",
			Value: bson.D{bson.E{Key: "user_followers",
				Value: userId}}},
		bson.E{Key: "$inc",
			Value: bson.D{bson.E{Key: "user_followers_count",
				Value: -1}}},
	}
	_, err = u.userCollection.UpdateMany(context.TODO(), filter1, update)
	if err != nil {
		return nil, err
	}
	return &unfollowing, nil
}
