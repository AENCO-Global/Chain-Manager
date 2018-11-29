package libs

import (
    "flag"
)

type Params struct {
    Lat     string
    Long    string
    ID      string
    Name    string
}

func GetParams() (getParams Params) {
    //var latPtr = flag.String("lat", "", "Latitude")
    //var lonPtr = flag.String("lon", "", "Longitude")
    //var idPtr =  flag.String("id", "", "Force the Node ID")
    //var namePtr =  flag.String("name", "", "Force the Node Name")
    flag.Parse()

    getParams = Params{
        //*latPtr,
        //*lonPtr,
        //*idPtr,
        //*namePtr
    }

    return
}