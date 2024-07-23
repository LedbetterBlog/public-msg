package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

// FindOne 查询单个文档
func (m *MongoDBPoolManager) FindOne(collectionName string, filter interface{}) (bson.M, error) {
	collection := m.db.Collection(collectionName)
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		//log.Printf("mongo FindOne is error: %v", err)
		return nil, err
	}
	return result, nil
}

// FindAll 查询多个文档
func (m *MongoDBPoolManager) FindAll(collectionName string, filter interface{}) ([]bson.M, error) {
	collection := m.db.Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Close 关闭 MongoDB 连接
func (m *MongoDBPoolManager) Close() error {
	return m.client.Disconnect(context.TODO())
}
