package main

import (
    "database/sql"
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
    "github.com/go-redis/redis"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
    "github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
    "github.com/iplay88keys/my-recipe-library/pkg/api/users"
    "github.com/iplay88keys/my-recipe-library/pkg/config"
    "github.com/iplay88keys/my-recipe-library/pkg/repositories"
    "github.com/iplay88keys/my-recipe-library/pkg/token"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    cfg.Port = "8080"
    cfg.Static = "ui/build"

    db, err := connectToMySQL(&cfg.MySQLCreds)
    if err != nil {
        panic(err)
    }

    redisClient, err := connectToRedis(cfg.RedisURL)
    if err != nil {
        panic(err)
    }

    recipesRepo := repositories.NewRecipesRepository(db)
    ingredientsRepo := repositories.NewIngredientsRepository(db)
    stepsRepo := repositories.NewStepsRepository(db)
    usersRepo := repositories.NewUsersRepository(db)

    redisRepo := repositories.NewRedisRepository(redisClient)
    tokenService := token.NewService(cfg.AccessSecret, cfg.RefreshSecret)

    a := api.New(&api.Config{
        Port:                  cfg.Port,
        StaticDir:             cfg.Static,
        Validate:              tokenService.ValidateToken,
        RetrieveAccessDetails: redisRepo.RetrieveTokenDetails,
        Endpoints: []*api.Endpoint{
            recipes.CreateRecipe(recipesRepo.Insert),
            recipes.ListRecipes(recipesRepo.List),
            recipes.GetRecipe(
                recipesRepo.Get,
                ingredientsRepo.GetForRecipe,
                stepsRepo.GetForRecipe,
            ),
            users.Register(
                usersRepo.ExistsByUsername,
                usersRepo.ExistsByEmail,
                usersRepo.Insert,
            ),
            users.Login(
                usersRepo.Verify,
                tokenService.CreateToken,
                redisRepo.StoreTokenDetails,
            ),
            users.Logout(
                tokenService.ValidateToken,
                redisRepo.DeleteTokenDetails,
            ),
        },
    })

    fmt.Printf("Serving at http://localhost:%s\n", cfg.Port)
    stopApi := a.Start()

    defer stopApi()
    defer disconnectFromMySQL(db)
    defer disconnectFromRedis(redisClient)

    blockUntilSigterm()
}

func connectToMySQL(config *config.MySQLCreds) (db *sql.DB, err error) {
    if config.URL != "" {
        var unquotedURL string
        url := config.URL

        unquotedURL, err = strconv.Unquote(url)
        if err == nil {
            url = unquotedURL
        }

        db, err = sql.Open("mysql", strings.TrimSpace(strings.TrimPrefix(url, "mysql://")))
        if err != nil {
            return nil, err
        }

    } else {
        cfg := mysql.Cfg(config.InstanceName, config.User, config.Password)
        cfg.DBName = config.DBName

        db, err = mysql.DialCfg(cfg)
        if err != nil {
            return nil, err
        }
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}

func connectToRedis(redisURL string) (redis.Cmdable, error) {
    options, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }

    redisClient := redis.NewClient(options)

    _, err = redisClient.Ping().Result()
    if err != nil {
        return nil, err
    }

    return redisClient, nil
}

func disconnectFromMySQL(db *sql.DB) {
    var stats sql.DBStats
    stats = db.Stats()

    var maxCount int
    for stats.InUse != 0 {
        if maxCount == 10 {
            break
        }

        stats = db.Stats()

        fmt.Printf("Waiting on open mySQL connections: %d in use\n", stats.InUse)

        maxCount += 1
        time.Sleep(100 * time.Millisecond)
    }

    err := db.Close()
    if err != nil {
        panic(err)
    }
}

func disconnectFromRedis(client redis.Cmdable) {
    err := client.(*redis.Client).Close()
    if err != nil {
        panic(err)
    }
}

func blockUntilSigterm() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

    <-sigs
}
