package main

import ()

func UserProfileHandler(c *SocketContext) error {

	return c.JSON(BaseData)
}

func init() {
	AddHandler("data", UserProfileHandler)
}
