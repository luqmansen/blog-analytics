package mongodb

import (
	"context"
	"github.com/luqmansen/web-analytics/analytics"
	"github.com/luqmansen/web-analytics/configs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(URL string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

func NewMongoRepository(m configs.Database) (analytics.AnalyticRepository, error) {

	client, err := newMongoClient(m.URI, m.Timeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.NewMongoRepository")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Second)
	defer cancel()

	coll := client.Database(m.Database).Collection("analytics")
	idx, err := coll.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bsonx.Doc{
				{"url", bsonx.Int32(1)},
				{"ip", bsonx.Int32(1)},
				{"created_at.date",  bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Infof("Index %s created\n", idx)

	return &mongoRepository{
		client:   client,
		database: m.Database,
		timeout:  time.Duration(m.Timeout) * time.Second,
	}, nil
}

func (m mongoRepository) GetAll() ([]*analytics.Analytic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.timeout)*time.Second)
	defer cancel()

	var data []*analytics.Analytic

	coll := m.client.Database(m.database).Collection("analytics")
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.GetAll")
	}
	if err = cur.All(ctx, &data); err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.GetAll")
	}
	return data, nil
}

func (m mongoRepository) Store(analytic *analytics.Analytic) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.timeout)*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection("analytics")
	_, err := coll.InsertOne(
		ctx,
		bson.M{
			"created_at": analytic.CreatedAt,
			"url":        analytic.URL,
			"ip":         analytic.IP,
			"info": bson.M{
				"ip":       analytic.Info.IP,
				"city":     analytic.Info.City,
				"country":  analytic.Info.Country,
				"location": analytic.Info.Location,
				"org":      analytic.Info.Org,
				"region":   analytic.Info.Region,
				"timezone": analytic.Info.Timezone,
			},
		},
	)
	if err != nil {
		var merr mongo.WriteException
		merr = err.(mongo.WriteException)
		errCode := merr.WriteErrors[0].Code

		if errCode == 11000 {
			return errors.Wrap(analytics.ErrorDuplicate, "repository.mongodb.Store")
		}
		return errors.Wrap(err, "repository.mongodb.Store")
	}
	return nil
}
