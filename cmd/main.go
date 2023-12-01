package main

func main() {
	daemonMode := true

	switch {
	case daemonMode:
		RunDaemon()
	default:
		RunController()
	}
}
