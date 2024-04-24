package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FetchDataParams struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type ResponsePayload struct {
	Code    int           `json:"code"`
	Msg     string        `json:"msg"`
	Records []MongoRecord `json:"records"`
}

type MongoRecord struct {
	Key        string    `bson:"key" json:"key"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	TotalCount int       `bson:"totalCount" json:"totalCount"`
}

func (s *Server) fetchDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Implemented", http.StatusNotImplemented)
		return
	}

	jsonEncoder := json.NewEncoder(w)

	var filter FetchDataParams
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, filter.StartDate)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)

		jsonEncoder.Encode(ResponsePayload{
			Code: http.StatusBadRequest,
			Msg:  "Invalid startTime format",
		})
		return
	}

	endDate, err := time.Parse(layout, filter.EndDate)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)

		jsonEncoder.Encode(ResponsePayload{
			Code: http.StatusBadRequest,
			Msg:  "Invalid endDate format",
		})
		return
	}

	if filter.MinCount == 0 {
		http.Error(w, "", http.StatusBadRequest)

		jsonEncoder.Encode(ResponsePayload{
			Code: http.StatusBadRequest,
			Msg:  "Missing minCount",
		})
		return
	}

	if filter.MaxCount == 0 {
		http.Error(w, "", http.StatusBadRequest)
		jsonEncoder.Encode(ResponsePayload{
			Code: http.StatusBadRequest,
			Msg:  "Missing maxCount",
		})
		return
	}

	filterQuery := bson.M{
		"createdAt": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	matchStage := bson.D{
		{Key: "$match", Value: filterQuery},
	}

	unwindStage := bson.D{{Key: "$unwind", Value: "$counts"}}

	groupStage := bson.D{
		{Key: "$group", Value: bson.M{
			"_id":        "$_id",
			"key":        bson.M{"$first": "$key"},
			"createdAt":  bson.M{"$first": "$createdAt"},
			"totalCount": bson.M{"$sum": "$counts"},
		}},
	}

	matchTotalCountStage := bson.D{
		{Key: "$match", Value: bson.M{
			"totalCount": bson.M{"$gte": filter.MinCount, "$lte": filter.MaxCount},
		}},
	}

	coll := s.mongoClient.Database(databaseName).Collection(collectionName)

	collection, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, groupStage, matchTotalCountStage})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	var records []MongoRecord

	if err = collection.All(context.TODO(), &records); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		jsonEncoder.Encode(ResponsePayload{
			Code: http.StatusInternalServerError,
			Msg:  "Internal Server Error",
		})
	}

	response := ResponsePayload{
		Code:    0,
		Msg:     "success",
		Records: records,
	}

	jsonEncoder.Encode(response)
}
