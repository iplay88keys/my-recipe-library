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

	"cloud.google.com/go/cloudsqlconn/mysql/mysql"
	"code.cloudfoundry.org/go-envstruct"
	"github.com/go-redis/redis"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	"github.com/iplay88keys/my-recipe-library/pkg/api/users"
	"github.com/iplay88keys/my-recipe-library/pkg/config"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
	"github.com/iplay88keys/my-recipe-library/pkg/services"
	"github.com/iplay88keys/my-recipe-library/pkg/token"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Config{
		Port:   "8080",
		Static: "ui/build",
	}

	err := envstruct.Load(&cfg)
	if err != nil {
		panic(err)
	}

	db, disconnectMySQL, err := connectToMySQL(cfg.MySQLCreds)
	if err != nil {
		panic(err)
	}
	defer disconnectMySQL()

	redisClient, err := connectToRedis(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	defer disconnectFromRedis(redisClient)

	// Create repositories
	recipesRepo := repositories.NewRecipesRepository(db)
	ingredientsRepo := repositories.NewIngredientsRepository(db)
	stepsRepo := repositories.NewStepsRepository(db)
	usersRepo := repositories.NewUsersRepository(db)
	redisRepo := repositories.NewRedisRepository(redisClient)
	tokenService := token.NewService(cfg.AccessSecret, cfg.RefreshSecret)

	// Create services
	recipeService := services.NewRecipeService(recipesRepo, ingredientsRepo, stepsRepo, db)
	userService := services.NewUserService(usersRepo, redisRepo, tokenService)

	a := api.New(tokenService, redisRepo, &api.Config{
		Port:      cfg.Port,
		StaticDir: cfg.Static,
		Endpoints: []*api.Endpoint{
			recipes.CreateRecipe(recipeService),
			recipes.ListRecipes(recipeService),
			recipes.GetRecipe(recipeService),
			users.Register(userService),
			users.Login(userService),
			users.Logout(userService),
		},
	})

	fmt.Printf("Serving at http://localhost:%s\n", cfg.Port)
	stopApi := a.Start()
	defer stopApi()

	blockUntilSigterm()
}

func connectToMySQL(config config.MySQLCreds) (db *sql.DB, disconnect func(), err error) {
	if config.URL != "" {
		var unquotedURL string
		url := config.URL

		unquotedURL, err = strconv.Unquote(url)
		if err == nil {
			url = unquotedURL
		}

		db, err = sql.Open("mysql", strings.TrimSpace(strings.TrimPrefix(url, "mysql://")))
		if err != nil {
			return nil, nil, err
		}
		disconnect = func() {
			disconnectFromMySQL(db)
		}
	} else {
		cleanup, err := mysql.RegisterDriver(
			"cloudsql-mysql",
		)
		if err != nil {
			return nil, nil, err
		}

		db, err = sql.Open(
			"cloudsql-mysql",
			fmt.Sprintf("%s:%s@cloudsql-mysql(%s)/%s", config.User, config.Password, config.InstanceName, config.DBName),
		)
		if err != nil {
			return nil, nil, err
		}

		disconnect = func() {
			cleanup()
		}
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	return db, disconnect, nil
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
