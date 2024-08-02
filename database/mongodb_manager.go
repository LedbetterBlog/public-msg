package database

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// MongoDBPoolManager 结构体
type MongoDBPoolManager struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDBPoolManager 创建 MongoDBPoolManager 实例
func NewMongoDBPoolManager(uri, dbName string) (*MongoDBPoolManager, error) {
	// 创建 MongoDB 客户端连接, 并设置超时时间
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).         // 连接超时时间
		SetServerSelectionTimeout(10 * time.Second). // 服务器选择超时时间
		SetSocketTimeout(5 * time.Second)            // 套接字读写超时时间
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("连接 MongoDB 失败: %v", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("测试 MongoDB 连接失败: %v", err)
	}

	// 获取指定数据库
	db := client.Database(dbName)

	// 成功连接并测试通过
	log.Println("MongoDB 连接成功")

	return &MongoDBPoolManager{client: client, db: db}, nil
}

// InsertData 向 MongoDB 插入数据
func (m *MongoDBPoolManager) InsertData(ctx context.Context, collectionName string, document interface{}) (string, error) {
	collection := m.db.Collection(collectionName)
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(string)
	if !ok {
		return "", err
	}
	return id, nil
}

// ReplaceOrUpdateData 更新或替换 MongoDB 中的数据
func (m *MongoDBPoolManager) ReplaceOrUpdateData(ctx context.Context, collectionName string, filter interface{}, replacement interface{}) (string, error) {
	collection := m.db.Collection(collectionName)
	result, err := collection.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return "", err
	}
	if result.MatchedCount == 0 {
		return "No document matched the filter", nil
	}
	return "Document replaced successfully", nil
}

// UpdateData 部分更新 MongoDB 中的数据
func (m *MongoDBPoolManager) UpdateData(ctx context.Context, collectionName string, filter interface{}, update interface{}) (int64, error) {
	collection := m.db.Collection(collectionName)
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

// FindOne 查询单个文档
func (m *MongoDBPoolManager) FindOne(ctx context.Context, collectionName string, filter interface{}) (bson.M, error) {
	collection := m.db.Collection(collectionName)
	var result bson.M

	fmt.Printf("Query filter: %v\n", filter)

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		// 查找的文档不存在要这么处理报错
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No document found for filter: %v", filter)
			// Debug: print all documents in the collection
			cursor, err := collection.Find(ctx, bson.M{})
			if err != nil {
				log.Printf("Error finding all documents: %v", err)
				return nil, err
			}
			defer func(cursor *mongo.Cursor, ctx context.Context) {
				err := cursor.Close(ctx)
				if err != nil {

				}
			}(cursor, ctx)
			var docs []bson.M
			if err = cursor.All(ctx, &docs); err != nil {
				log.Printf("Error decoding all documents: %v", err)
				return nil, err
			}
			log.Printf("Documents in collection: %v", docs)
			return nil, err
		}
		log.Printf("mongo FindOne error: %v", err)
		return nil, err
	}
	return result, nil
}

// Close 关闭 MongoDB 连接
func (m *MongoDBPoolManager) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
