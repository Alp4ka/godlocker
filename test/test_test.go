package test

import (
	"context"
	"flag"
	"fmt"
	"github.com/Alp4ka/godlocker/clientlocker"
	"github.com/Alp4ka/godlocker/redislocker"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

type test_testType string
type test_testCase func(t *testing.T)

const (
	test_redisTestLabel test_testType = "redis"
)

var (
	test_redisAddr     string
	test_redisPassword string
	test_redisDB       int
	test_redisUsername string
	test_testFunc      test_testCase
)

const (
	test_unimplemented = "unimplemented!"
)

var (
	test_mapping = map[test_testType]test_testCase{
		test_redisTestLabel: test_RedisLocker,
	}
)

func test_mapTestType(testType test_testType) test_testCase {
	if f, ok := test_mapping[testType]; !ok {
		panic(test_unimplemented)
	} else {
		return f
	}
}

func init() {
	test_redisAddr = os.Getenv("redis-addr")
	test_redisPassword = os.Getenv("redis-pass")
	test_redisUsername = os.Getenv("redis-username")
	val, err := strconv.Atoi(os.Getenv("redis-db"))
	if err != nil {
		panic(fmt.Sprintf("Error while parsing redis-db: %s", err.Error()))
	}
	test_redisDB = val
	test_testFunc = test_mapTestType(test_testType(*flag.String("type", string(test_redisTestLabel), "Setup test type!")))
}

func Test_Main(t *testing.T) {
	t.Log("Test started!")

	test_testFunc(t)

	t.Log("Test finished!")
}

func test_Lock(t *testing.T, clientID int) error {
	workID := rand.Int()

	t.Logf("[Lock-%d-%d] Init hard work\n", workID, clientID)
	mu, err := clientlocker.CL().Lock(context.TODO(), clientID)
	if err != nil {
		t.Logf("[Lock-%d-%d] Hard work lock error: %s\n", workID, clientID, err)
		return err
	}
	t.Logf("[Lock-%d-%d] Start hard work\n", workID, clientID)
	time.Sleep(time.Second * 15)
	t.Logf("[Lock-%d-%d] End hard work\n", workID, clientID)

	err = clientlocker.CL().Unlock(context.TODO(), mu)
	if err != nil {
		t.Logf("[Lock-%d-%d] Hard work unlock error1: %s\n", workID, clientID, err)
		return err
	}
	t.Logf("[Lock-%d-%d] Released hard work\n", workID, clientID)
	return nil
}

func test_TryLock(t *testing.T, clientID int) error {
	workID := rand.Int()

	_, err := clientlocker.CL().TryLock(nil, clientID)
	if err != nil {
		t.Logf("[TryLock-%d-%d] Hard work lock error: %s\n", workID, clientID, err)
		return err
	}
	return nil
}

func test_RedisLocker(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     test_redisAddr,
		Password: test_redisPassword,
		DB:       test_redisDB,
		Username: test_redisUsername,
	})
	defer client.Conn().Close()

	redisLocker := redislocker.NewRedisLocker(client)
	clientLocker := clientlocker.NewClientLocker(redisLocker)
	clientlocker.ReplaceGlobals(clientLocker)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		err := test_Lock(t, 1)
		if err != nil {
			t.Errorf("1-1 Lock shouldn't throw an error in this scenario")
		}

	}()
	go func() {
		defer wg.Done()
		err := test_Lock(t, 2)
		if err != nil {
			t.Errorf("2-1 Lock shouldn't throw an error in this scenario")
		}
	}()

	time.Sleep(time.Second * 3)
	go func() {
		defer wg.Done()
		err := test_Lock(t, 1)
		if err != nil {
			t.Errorf("1-2 Lock shouldnt throw an error in this scenario")
		}
	}()
	wg.Wait()

}
