package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Foxtrot-14/FitRang/notification-service/config"
	"github.com/Foxtrot-14/FitRang/notification-service/eventbus"
	"github.com/Foxtrot-14/FitRang/notification-service/graph"
	"github.com/Foxtrot-14/FitRang/notification-service/repository"
	"github.com/Foxtrot-14/FitRang/notification-service/services"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const defaultPort = "8080"

func main() {

	godotenv.Load(".env")
	kafkaCfg := config.LoadKafkaConfig()

	eventBus, err := eventbus.NewEventBus(eventbus.Config{
		Brokers: kafkaCfg.Brokers,
	})
	if err != nil {
		log.Fatalf("failed to init event bus: %v", err)
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo Connect Error: ", err)
	}

	db := client.Database("notification-service")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	messageRepo := repository.NewMessageRepository(db)
	messageService := services.NewMessageService(messageRepo)

	_, err = eventBus.NewConsumer(
		"notification-service",
		[]string{"notification"},
		func(key, value []byte) error {
			if err := messageService.CreateMessage(value); err != nil {
				log.Printf("error while creating message: %s", err)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("failed to init event bus: %v", err)
	}

	resolver := &graph.Resolver{
		MessageService: messageService,
	}

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolver,
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
