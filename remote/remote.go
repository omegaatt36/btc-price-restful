package remote

import (
	"encoding/json"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// APIConfig refer to APIconfig.josn
type apiConfig struct {
	Auth  string `json:"auth"`
	Delay int    `json:"delay"`
}

var apiConfigs map[string]apiConfig

type apiSet []API

// ParseConfig will call a file like APIconf.json and the format is
/* {
    	"sourceName": {
        	"auth":"auth token",
        	"delay" : second
		},
		...
	}
*/
func ParseConfig(file *os.File) error {
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&apiConfigs)
	if err != nil {
		return err
	}
	return nil
}

// InitAPIs init all remote API and get first
// If fist call failed, will get latest id from redis
func InitAPIs() {
	apis := make(apiSet, 0)

	cmcFactory := CoinMarketCapFactory{}
	cmcConf, ok := apiConfigs[sourceNameCoinMarketCap]
	if !ok {
		logrus.Warning("can't find key:", sourceNameCoinMarketCap, ", please check APIconfig.json")
	} else {
		cmc, _ := cmcFactory.Create(cmcConf.Auth)
		apis = append(apis, cmc)
	}

	cgFactory := CoinGeckoFactory{}
	cgConf, ok := apiConfigs[sourceNameCoinGecko]
	if !ok {
		logrus.Warning("can't find key:", sourceNameCoinGecko, ", please check APIconfig.json")
	} else {
		cg, _ := cgFactory.Create(cgConf.Auth)
		apis = append(apis, cg)
	}
	cdFactory := CoinDeskFactory{}
	cdConf, ok := apiConfigs[sourceNameCoinDesk]
	if !ok {
		logrus.Warning("can't find key:", sourceNameCoinDesk, ", please check APIconfig.json")
	} else {
		cd, _ := cdFactory.Create(cdConf.Auth)
		apis = append(apis, cd)
	}

	apis.firstCall()

	for _, v := range apis {
		go runTicker(v)
	}
}

func (apis apiSet) firstCall() {
	for _, v := range apis {
		go firstCallAPI(v)
	}
}

func firstCallAPI(api API) {
	err := api.CallRemote()
	if err != nil {
		logrus.Info(err.Error())
		err = api.InitFormRedis()
		if err != nil {
			logrus.Info(err.Error())
		}
		return
	}
	err = api.InsertDB()
	if err != nil {
		logrus.Info(err.Error())
		return
	}
}

func runTicker(api API) {
	name := api.GetSourceName()
	d := time.Duration(time.Second * time.Duration(apiConfigs[name].Delay))
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		<-ticker.C
		err := api.CallRemote()
		if err != nil {
			logrus.Info(err.Error())
			continue
		}
		err = api.InsertDB()
		if err != nil {
			logrus.Info(err.Error())
			continue
		}
		logrus.Debugf("%s get new data", name)
	}
}
