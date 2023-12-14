package handler

import "log"

// Handler package entrypoint
func PortForward(r string, h string, p string, lp string) {

	log.Println("Opening connection for:", h)
	createSession(getBastion(), h, p, lp)
}
