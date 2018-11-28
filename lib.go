package main



import (
    "encoding/json"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/go-ini/ini"
    "io/ioutil"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
)

/* --------------------------------------------------------------------------------
#     _____      _     ____  _            _      _    _      _       _     _
#    / ____|    | |   |  _ \| |          | |    | |  | |    (_)     | |   | |
#   | |  __  ___| |_  | |_) | | ___   ___| | __ | |__| | ___ _  __ _| |__ | |_
#   | | |_ |/ _ \ __| |  _ <| |/ _ \ / __| |/ / |  __  |/ _ \ |/ _` | '_ \| __|
#   | |__| |  __/ |_  | |_) | | (_) | (__|   <  | |  | |  __/ | (_| | | | | |_
#    \_____|\___|\__| |____/|_|\___/ \___|_|\_\ |_|  |_|\___|_|\__, |_| |_|\__|
#                                                               __/ |
#                                                              |___/
    This function needs a better way to get the block height, for now nem2-cli
   --------------------------------------------------------------------------------*/
type Blocks struct {
    Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Height      int64 `json:"height"`
}

func getBlockHeight()  int64 {
    var height Blocks
    if conf.mongoDB[0] == conf.mongoDB[1] {
        defer func() { //Catch errors, and resume
            r := recover()
            if r != nil {
                log.Error("Database Access problem:", r)
            } } ()
        session, err := mgo.Dial(conf.Database)
        if err != nil {
            conf.mongoDB[1] = 0
            panic(err)
        }
        defer session.Close()
        err = session.DB(conf.DBName).C("chainInfo").Find(nil).One(&height)
        if err != nil {
            conf.mongoDB[1] = 0
            log.Warning("No DB Space when trying to get height")
            panic(err)
        }
    } else {
        log.Warning("Skipping Database connection, will retry in [",conf.mongoDB,"] cycles!")
        conf.mongoDB[1]++
    }
    return height.Height
}

// -------------------------------------------------------------------
//     _____      _                          _   _     _
//    / ____|    | |       /\               | | | |   (_)
//   | |  __  ___| |_     /  \   _ __  _   _| |_| |__  _ _ __   __ _
//   | | |_ |/ _ \ __|   / /\ \ | '_ \| | | | __| '_ \| | '_ \ / _` |
//   | |__| |  __/ |_   / ____ \| | | | |_| | |_| | | | | | | | (_| |
//    \_____|\___|\__| /_/    \_\_| |_|\__, |\__|_| |_|_|_| |_|\__, |
//                                      __/ |                   __/ |
//                                     |___/                   |___/
func getAnything(file string ,section string, key string) string{
    cfgFile := getRoot() + "/resources/" + file
    defer func() { //Catch errors, and resume
        r := recover()
        if r != nil {
            log.Error("Unable get values (Possible unable to find:", cfgFile, ") err", r)
        } } ()
    cfg, err := ini.Load(cfgFile )
    retVal := cfg.Section(section).Key(key).String()
    if err != nil {
        panic(err)
    }
    return retVal
}

// ------------------------------------------------------------------------
//     _____                                      _   _     _
//    / ____|                   /\               | | | |   (_)
//   | (___   __ ___   _____   /  \   _ __  _   _| |_| |__  _ _ __   __ _
//    \___ \ / _` \ \ / / _ \ / /\ \ | '_ \| | | | __| '_ \| | '_ \ / _` |
//    ____) | (_| |\ V /  __// ____ \| | | | |_| | |_| | | | | | | | (_| |
//   |_____/ \__,_| \_/ \___/_/    \_\_| |_|\__, |\__|_| |_|_|_| |_|\__, |
//                                           __/ |                   __/ |
//                                          |___/                   |___/
func saveAnything(file string ,section string, key string, value string) bool{
    cfgFile := conf.Root+"/resources/"+file
    defer func() { //Catch errors, and resume
        r := recover()
        if r != nil {
            log.Error("Unable to save values (Possibly unable to find or write file:", r)
        } } ()
    cfg, err := ini.Load(cfgFile )
    cfg.Section(section).Key(key).SetValue(value)
    if err != nil {
        panic(err)
    }
    err = cfg.SaveTo(cfgFile)
    if err != nil {
        panic(err)
    }
    return true
}



// ------------------------------------
//     _____ ______ ____ _____ _____
//    / ____|  ____/ __ \_   _|  __ \
//   | |  __| |__ | |  | || | | |__) |
//   | | |_ |  __|| |  | || | |  ___/
//   | |__| | |___| |__| || |_| |
//    \_____|______\____/_____|_|
func getGeo() {
    defer func() {
        r := recover()
        if r != nil { log.Error("GEO Details Denied", r) }
    }()
    type Geo struct {
        IP      string `json:"ip"`
        City    string `json:"city"`
        Region  string `json:"region"`
        Country string `json:"country"`
        Loc     string `json:"loc"`
        Org     string `json:"org"`
    }
    transport := &http.Transport{DisableKeepAlives: true}
    client := http.Client{Transport: transport}
    resp, err := client.Get("http://ipinfo.io")
    if err != nil {
        panic(err)
    } else {
        responseData, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Error("Response contains Errors:")
            panic(err)
        } else {
            bytes := []byte(responseData)
            var inGeo Geo
            err := json.Unmarshal(bytes, &inGeo)
            if err != nil {
                panic(err)
            }
            latLon := strings.Split(inGeo.Loc, ",")

            saveAnything("config-agent.ini", "report", "city", inGeo.City)
            saveAnything("config-agent.ini", "report", "country", inGeo.Country)
            saveAnything("config-agent.ini", "report", "region", inGeo.Region)
            saveAnything("config-agent.ini", "report", "lat",  latLon[0] )
            saveAnything("config-agent.ini", "report", "lon", latLon[1])
            saveAnything("config-agent.ini", "report", "publicIp", inGeo.IP)

        }
        defer resp.Body.Close()
    }
}

// --------------------------------------------------
//    ___  _  _            ___       _      _
//   | __|(_)| | ___  ___ | __|__ __(_) ___| |_  ___
//   | _| | || |/ -_)|___|| _| \ \ /| |(_-<|  _|(_-<
//   |_|  |_||_|\___|     |___|/_\_\|_|/__/ \__|/__/
func exists(filePath string) (exists bool) {
    _,err := os.Stat(filePath)
    if err != nil {
        exists = false
    } else {
        exists = true
    }
    return
}


// ----------------------------------------------------------------------------
//      ____                  __                   ____
//     / __ \____ _____  ____/ /___  ____ ___     / __ \____ _____  ____ ____
//    / /_/ / __ `/ __ \/ __  / __ \/ __ `__ \   / /_/ / __ `/ __ \/ __ `/ _ \
//   / _, _/ /_/ / / / / /_/ / /_/ / / / / / /  / _, _/ /_/ / / / / /_/ /  __/
//  /_/ |_|\__,_/_/ /_/\__,_/\____/_/ /_/ /_/  /_/ |_|\__,_/_/ /_/\__, /\___/
//                                                               /____/
func randomRange(min int, max int) (randomRange int){
    rand.Seed(time.Now().UnixNano())
    randomRange = rand.Intn(max-min)+min
    return
}
