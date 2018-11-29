package libs

import (
    "github.com/go-ini/ini"
    "github.com/op/go-logging"
    "math/rand"
    "os"
    "time"
)

var log = logging.MustGetLogger("aen.manager")

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
func getBlockHeight()   ( getBlockHeight int64) {
    //var height Blocks
    //if conf.mongoDB[0] == conf.mongoDB[1] {
    //    defer func() { //Catch errors, and resume
    //        r := recover()
    //        if r != nil {
    //            log.Error("Database Access problem:", r)
    //        } } ()
    //    session, err := mgo.Dial(conf.Database)
    //    if err != nil {
    //        conf.mongoDB[1] = 0
    //        panic(err)
    //    }
    //    defer session.Close()
    //    err = session.DB(conf.DBName).C("chainInfo").Find(nil).One(&height)
    //    if err != nil {
    //        conf.mongoDB[1] = 0
    //        log.Warning("No DB Space when trying to get height")
    //        panic(err)
    //    }
    //} else {
    //    log.Warning("Skipping Database connection, will retry in [",conf.mongoDB,"] cycles!")
    //    conf.mongoDB[1]++
    //}
        getBlockHeight = 0
    return
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
    cfgFile := getRoot() + "/" + file
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
