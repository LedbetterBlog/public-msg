package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoDBPoolManager 结构体
type MongoDBPoolManager struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDBPoolManager 创建 MongoDBPoolManager 实例
func NewMongoDBPoolManager(uri, dbName string) (*MongoDBPoolManager, error) {
	// 创建 MongoDB 客户端连接
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("连接 MongoDB 失败: %v", err)
	}

	// 测试连接
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("测试 MongoDB 连接失败: %v", err)
	}

	// 获取指定数据库
	db := client.Database(dbName)

	// 成功连接并测试通过
	log.Println("MongoDB 连接成功")

	return &MongoDBPoolManager{client: client, db: db}, nil
}

// InsertData 向 MongoDB 插入数据
func (m *MongoDBPoolManager) InsertData(collectionName string, document interface{}) (string, error) {
	collection := m.db.Collection(collectionName)
	result, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(string)
	if !ok {
		return "", err
	}
	return id, nil
}

// Close 关闭 MongoDB 连接
func (m *MongoDBPoolManager) Close() error {
	return m.client.Disconnect(context.TODO())
}
