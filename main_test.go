package main

import (
	"testing"
)

type pickup struct {
	userID       string
	pickupFloor  int
	dropOffFloor int
}

var testcases = []struct {
	numberOfElevators int
	numberOfFloors    int
	pickUps           []pickup
	out               string
}{
	{
		numberOfElevators: 16,
		numberOfFloors:    10,
		pickUps: []pickup{
			{userID: "User1", pickupFloor: 0, dropOffFloor: 1},
			{userID: "User2", pickupFloor: 2, dropOffFloor: 10},
			{userID: "User3", pickupFloor: 4, dropOffFloor: 10},
			{userID: "User4", pickupFloor: 9, dropOffFloor: 10},
			{userID: "User5", pickupFloor: 10, dropOffFloor: 9},
			{userID: "User6", pickupFloor: 5, dropOffFloor: 10},
			{userID: "User7", pickupFloor: 9, dropOffFloor: 3},
			{userID: "User8", pickupFloor: 1, dropOffFloor: 10},
			{userID: "User9", pickupFloor: 3, dropOffFloor: 10},
			{userID: "User10", pickupFloor: 2, dropOffFloor: 10},
			{userID: "User11", pickupFloor: 5, dropOffFloor: 10},
			{userID: "User12", pickupFloor: 3, dropOffFloor: 1},
			{userID: "User13", pickupFloor: 4, dropOffFloor: 3},
			{userID: "User14", pickupFloor: 7, dropOffFloor: 7},
			{userID: "User15", pickupFloor: 10, dropOffFloor: 6},
			{userID: "User16", pickupFloor: 1, dropOffFloor: 4},
			{userID: "User17", pickupFloor: 3, dropOffFloor: 0},
			{userID: "User18", pickupFloor: 9, dropOffFloor: 1},
			{userID: "User19", pickupFloor: 2, dropOffFloor: 7},
			{userID: "User20", pickupFloor: 0, dropOffFloor: 5},
		},
		out: "",
	},
	/*
		{
			numberOfElevators: 2,
			numberOfFloors:    30,
			pickUps:
				[]pickup{
						{userID: "User1", pickupFloor:  0, dropOffFloor: 1,},
						{userID: "User2", pickupFloor:  12, dropOffFloor: 10,},
						{userID: "User3", pickupFloor:  24, dropOffFloor: 10,},
						{userID: "User4", pickupFloor:  19, dropOffFloor: 13,},
						{userID: "User5", pickupFloor:  10, dropOffFloor: 9,},
						{userID: "User6", pickupFloor:  5, dropOffFloor: 17,},
						{userID: "User7", pickupFloor:  9, dropOffFloor: 30,},
						{userID: "User8", pickupFloor:  1, dropOffFloor: 10,},
						{userID: "User9", pickupFloor:  3, dropOffFloor: 10,},
						{userID: "User10", pickupFloor:  2, dropOffFloor: 10,},
						{userID: "User11", pickupFloor:  5, dropOffFloor: 10,},
						{userID: "User12", pickupFloor:  3, dropOffFloor: 21,},
						{userID: "User13", pickupFloor:  4, dropOffFloor: 23,},
						{userID: "User14", pickupFloor:  7, dropOffFloor: 17,},
					},
			out: "",
		},

	*/
}

func TestElevatorControlSystem(t *testing.T) {
	t.Parallel()

	for _, c := range testcases {
		control := NewElevatorControlSystem(c.numberOfElevators, c.numberOfFloors)
		for i := range c.pickUps {
			pickup := c.pickUps[i]
			control.PickUpButtonWasPushed(pickup.userID, pickup.pickupFloor, pickup.dropOffFloor)
			control.Step()
		}
	}
}

func BenchmarkElevatorControlSystem(b *testing.B) {
	println(b.N)

	for i := 0; i < b.N; i++ {
		for _, c := range testcases {
			control := NewElevatorControlSystem(c.numberOfElevators, c.numberOfFloors)
			for i := range c.pickUps {
				pickup := c.pickUps[i]
				control.PickUpButtonWasPushed(pickup.userID, pickup.pickupFloor, pickup.dropOffFloor)
			}
		}

	}
}
