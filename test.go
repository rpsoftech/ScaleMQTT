package main

import (
	"fmt"

	MT_SICS "github.com/bas-dehaan/MT-SICS"
)

func main() {
	// Connect to the scale via COM1
	connection, err := MT_SICS.Connect("/dev/cu.usbserial-0001")
	defer connection.Close() // Close the connection when the program ends
	if err != nil {
		panic(err)
	}

	// Get the scale out of standby
	// err = MT_SICS.PowerOn(connection)
	// if err != nil {
	// 	panic(err)
	// }

	// Weigh a sample
	_ = MT_SICS.SetMessage(connection, "HHHHH")
	measurement, err := MT_SICS.Weight(connection)
	if err != nil {
		panic(err)
	}

	// Print the measurement
	fmt.Println(measurement)
}
