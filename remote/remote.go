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
	cmcFactory := CoinMarketCapFactory{}
	apis := make(apiSet, 0)
	value, ok := apiConfigs[sourceNameCoinMarketCap]
	if !ok {
		logrus.Warning("can't find key:", sourceNameCoinMarketCap, ", please check auth.json")
	} else {
		cmc, _ := cmcFactory.Create(value.Auth)
		apis = append(apis, cmc)
	}

	apis.firstCall()

	for _, v := range apis {
		go runTicker(v)
	}
}

func (apis apiSet) firstCall() {
	for _, v := range apis {
		err := v.CallRemote()
		if err != nil {
			logrus.Info(err.Error())
			err = v.InitFormRedis()
			if err != nil {
				logrus.Info(err.Error())
			}
			continue
		}
		err = v.InsertDB()
		if err != nil {
			logrus.Info(err.Error())
			continue
		}
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
