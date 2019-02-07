package main

import (
	"echo-rest-api/api"
	"echo-rest-api/config"
	"echo-rest-api/service"
	"echo-rest-api/store"
	"flag"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	log.SetFormatter(&log.JSONFormatter{})
	configFile := flag.String("c", "config.yaml", "Path to config file")
	flag.Parse()
	var conf *config.Config
	if conf, err = config.NewConfig(*configFile); err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.Level(conf.LogLevel))
	log.Info("Starting service with configuration: ", conf.ConfigFile)
	store, err := store.NewStore(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	log.Info("Store created successfully")
	cs := service.NewCategoryService(store)
	ps := service.NewProductService(store)
	log.Info("Services created successfully")
	api := api.NewApi(conf, cs, ps)
	log.WithField("address", api.GetApiInfo().Address).
		WithField("mw", api.GetApiInfo().MW).
		WithField("routs", api.GetApiInfo().Routs).
		Info("Starting api")
	log.Fatal(api.Start())
}
