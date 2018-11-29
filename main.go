package main

import (
    "fmt"
    "github.com/AENCO-Global/Chain-Manager/libs"
    "github.com/op/go-logging"
)

// ---------------------------
//    __  __       _
//   |  \/  |     (_)
//   | \  / | __ _ _ _ __
//   | |\/| |/ _` | | '_ \
//   | |  | | (_| | | | | |
//   |_|  |_|\__,_|_|_| |_|
func main() {
//    log.Critical("----- Starting Agent ver:", conf.Version ," Up ----- \n")
    var conf = libs.getGlobals()
    var log = logging.MustGetLogger("aen.manager")
    var params = libs.getParams()

    for { // Keep the application alive, and reduce CPU usage
//        time.Sleep(time.Duration(conf.Heartbeat) * time.Minute )
        fmt.Println("Keep Alive:")
    }
//    log.Critical("Terminating Ver:", conf.Version  ," Agent ----- End ", " Execution -----")
}
//    ______           _
//   |  ____|         | |
//   | |__   _ __   __| |
//   |  __| | '_ \ / _` |
//   | |____| | | | (_| |
//   |______|_| |_|\__,_|
//   --------------------
