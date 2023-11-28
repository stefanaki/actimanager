package main

func main() {
	daemonMode := false

	switch {
	case daemonMode:
		runDaemon()
	default:
		runControllers()
	}
}
