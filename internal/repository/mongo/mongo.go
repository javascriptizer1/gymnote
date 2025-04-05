package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"gymnote/internal/config"
)

var (
	colExercises = "exercises"
	colSessions  = "training_sessions"
	colLogs      = "training_logs"
)

type mongodb struct {
	db           *mongo.Database
	exerciseColl *mongo.Collection
	sessionColl  *mongo.Collection
	logColl      *mongo.Collection
	cfg          *config.DBConfig
}

func New(ctx context.Context, cfg *config.DBConfig) (*mongodb, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.ConnectionString()).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongo ping error: %w", err)
	}

	db := client.Database(cfg.Name)

	m := &mongodb{
		db:           db,
		cfg:          cfg,
		exerciseColl: db.Collection(colExercises),
		sessionColl:  db.Collection(colSessions),
		logColl:      db.Collection(colLogs),
	}

	if err := m.ensureIndexes(ctx); err != nil {
		return nil, fmt.Errorf("mongo ensure indexes err: %w", err)
	}

	return m, nil
}

func (m *mongodb) Close(ctx context.Context) error {
	return m.db.Client().Disconnect(ctx)
}

func (m *mongodb) ensureIndexes(ctx context.Context) error {
	if _, err := m.logColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "exercise_id", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.M{"session_date": 1},
			Options: options.Index().SetUnique(false),
		},
	}); err != nil {
		return err
	}

	if _, err := m.sessionColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "date", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}); err != nil {
		return err
	}

	if _, err := m.exerciseColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
	}); err != nil {
		return err
	}

	return nil
}
