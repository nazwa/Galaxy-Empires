package main

import ()

func UserProfileHandler(c *SocketContext) error {

	return SendJson(Store, c.Session)
}

func init() {
	AddHandler("data", UserProfileHandler)
}
