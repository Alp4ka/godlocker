# Go Distributed Locker (Go D Locker)
Sugared distributed locker implementation for golang 1.20.

# Usage 
```golang
import (
  "github.com/Alp4ka/godlocker/clientlocker"
  redislocker "github.com/Alp4ka/godlocker/redis"
  "context"
  "fmt"
  "math/rand"
  "time"
)

func clientLock(context context.Context, clientID int) {
  workID := rand.Int()

  mu, err := clientlocker.CL().Lock(ctx, clientID)
  if err != nil {
    fmt.Printf("[%d-%d] Hard work lock error: %s\n", workID, clientID, err)
    return
  }

  fmt.Printf("[%d-%d] Start hard work\n", workID, clientID)
  time.Sleep(time.Second * 10)
  fmt.Printf("[%d-%d] End hard work\n", workID, clientID)

  err = clientlocker.CL().Unlock(ctx, mu)
  if err != nil {
    fmt.Printf("[%d-%d] Hard work unlock error: %s\n", workID, clientID, err)
  }
}

func main() {
  // Here we create redis connection(redisConn)
  
  // Use redis locker. (It's possible to create custom distributed locker storage)
  locker := redislocker.NewRedisLocker(redisConn)
  
  // Custom locker implementation that uses clientID as locker label.
  clientLocker := clientlocker.NewClientLocker(redisLocker)
  
  // Since this moment we can call clientlocker.CL() to call global instance.
  clientlocker.ReplaceGlobals(clientLocker)
  
  ctx := context.Background()
  
  // Has to acquire clientID 1 locker.
  go clientLock(ctx, 1)
  
  // Has to acquire clientID 2 locker.
  go clientLock(ctx, 2)
  
  // Now we have to wait until clientID 1 lock releases it's resources.
  time.Sleep(time.Second)
  go clientLock(ctx, 1)
  
  return
}
```
