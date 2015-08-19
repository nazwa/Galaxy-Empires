package main

import ()

func UserProfileHandler(c *SocketContext) error {

	return c.JSON(Store)
}

func init() {
	AddHandler("data", UserProfileHandler)
}
