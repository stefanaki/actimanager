package main

func main() {
	daemonMode := true

	switch {
	case daemonMode:
		runDaemon()
	default:
		runController()
	}
}
