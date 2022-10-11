package main

import Routers "restapi/routers"

func main() {

	var PORT = ":3000"

	Routers.StartServer().Run(PORT)
}
