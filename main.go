package main

import (
	"fmt"
	"math"
)

/***********************************************************************
 ***********************************************************************
 ***                                                                 ***
 ***  ELEVATOR CONTROL SYSTEM INTERFACE DEFINITION AND FUNCIONALITY  ***
 ***                                                                 ***
 ***********************************************************************
 ***********************************************************************/

// The Elevator Control System Interface
type ElevatorControlSystem interface {
	Status()
	Update(elevatorID int, floor int, direction string)
	PickUpButtonWasPushed(userID string, pickUpFloor int, dropOffFloor int)
	Step()
}

// Stores the information generated the Elevator Control System
type elevatorControlSystem struct {
	Elevators    []Elevator // List of the elevators in our system and their current status
	NUMELEVATORS int        // Number of elevators in our system
	TOPFLOOR     int        // Top floor building in our system
}

/**
 * Initializes the Elevator Control System Intefrace
	@ numberOfElevators int
	@ numberOfFloors int
*/
func NewElevatorControlSystem(numberOfElevators int, numberOfFloors int) ElevatorControlSystem {
	control := &elevatorControlSystem{
		Elevators:    []Elevator{},
		NUMELEVATORS: numberOfElevators,
		TOPFLOOR:     numberOfFloors,
	}

	for i := 0; i < numberOfElevators; i++ {
		control.Elevators = append(control.Elevators, NewElevator(i, control.TOPFLOOR))
	}

	return control
}

/**
 * 	It tells the status of every elevator of the system: which floor is in, if it is going up, going down
	or if it is stopped. It also will tell the list of tasks it has been assigned so far from the list of
	pick-up requests done by the users so far
*/
func (control *elevatorControlSystem) Status() {
	fmt.Printf("\nCURRENT STATUS OF THIS ELEVATOR CONTROL SYSTEM: IS MANAGING %d PLANTS AND %d ELEVATORS\n"+
		"======================================================================================\n",
		control.TOPFLOOR, control.NUMELEVATORS)
	for i := range control.Elevators {
		elev := control.Elevators[i]
		fmt.Printf("\n* Elevator %d is in floor %d", i, elev.getFloorNumber())
		trips := elev.getAssignedTrips()
		if len(trips) > 0 {
			fmt.Printf(", going %v, and it has been assigned %d tasks:\n\n", elev.getDirection(), len(trips))
			for j := range trips {
				trip := trips[j]
				fmt.Printf(" - %v is in floor %d %v. Wants to go to floor %d.\n", trip.userID, trip.fromFloor, trip.userAction, trip.toFloor)
			}
		} else {
			fmt.Printf(", stopped.")
		}
	}
}

/****
* It allows to the elevator controler system to intentionally change the floor and the direction of an elevator,
  without the intervention of a user pressing the pick-up button. It could be the equivalent to an engineer
  using his master key when they are fixing an elevator in a building
*/
func (control *elevatorControlSystem) Update(elevatorID int, floor int, direction string) {
	control.Elevators[elevatorID].setFloorNumber(floor)
	control.Elevators[elevatorID].setDirection(direction)
}

/**
* For the sake of simplicity, let's assume that a PickUpButtonWasPushed action is composed of the following actions:

  1) A user located in pickUpFloor pushes the button to call an elevator
  2) When one of the elevators arrives to the pickUpFloor, the user enters and pushes the button corresponding
     to his desired dropOffFloor
*/
func (control *elevatorControlSystem) PickUpButtonWasPushed(userID string, pickUpFloor int, dropOffFloor int) {
	// Choose the most optimal elevator and assign it this the current pick-up task, defined by the parameters
	chooseTheMostOptimalElevator(control, userID, pickUpFloor, dropOffFloor)
}

/**
 *  Builds the step list of actions performed by the elevators to complete the tasks they have been
	assigned by the control system and calls PrintStepListSimulation() once it is ready in order to to
    print it
*/
func (control *elevatorControlSystem) Step() {
	for i := 0; i < len(control.Elevators); i++ {
		control.Elevators[i].Step()
	}
	control.printStepListSimulation()
}

/***** ELEVATOR CONTROL SYSTEM HELPER FUNCTIONS, NOT EXPOSED IN THE INTERFACE *************/
func (control *elevatorControlSystem) printStepListSimulation() {
	fmt.Printf("\n\nSTEP LIST OF OUR SYSTEM OF %d FLOORS AND %d ELEVATORS\n"+
		"=====================================================\n", control.TOPFLOOR, control.NUMELEVATORS)
	for i := range control.Elevators {
		fmt.Printf("\nElevator %d Performed This Step List\n----------------------------------------\n", i)
		for j := range control.Elevators[i].getStepList() {
			step := control.Elevators[i].getStep(j)
			// fmt.Printf("%v\n", step)
			if step.userAction == exitingFromElevator {
				fmt.Printf("Floor %d, going %v. "+
					"%v is %v in floor %d.\n",
					step.elevInFloor, step.elevDirection, step.userID, step.userAction, step.elevInFloor)
			} else {
				if step.userAction == gettingIntoAElevator {
					fmt.Printf("Floor %d, going %v. "+
						"%v is %v.\n",
						step.elevInFloor, step.elevDirection, step.userID, step.userAction)
				} else {
					fmt.Printf("Floor %d, going %v. "+
						"%v pressed the pick-up button in floor %d and wants to go to floor %d. %v.\n",
						step.elevInFloor, step.elevDirection, step.userID, step.fromFloor, step.toFloor, step.userAction)
				}
			}
		}
	}
}

/**
 *
   This Is the scheduler optimizer.

   I've chosen this set of rules to assign the best elevator for every pick-up request:

	- When a user pushes the PickUp button at any floor, the Elevator chosen to give him service will be elected attending to:

		1) The elevatorProximity: the difference between the floor where the nearest elevator is and the user's pickup floor
   		   This rule will be rejected in favour of the next one, in case it happens ...

		2) The elevator matching more than 7 in the value to 'elevatorWithMoreDropOffTripsMatchingTheRequestedDropOff', assuming
   		   this rule forged only in my imagination: the most people going to the same plant the more savings in electricity (maybe
   			it is not true in the real world, but ... I didn't want to leave just a FIFO queue of pickups and drop-offs)
*/
func chooseTheMostOptimalElevator(control *elevatorControlSystem, userID string, pickUpFloor int, dropOffFloor int) Elevator {
	var chosenElevator int
	var nearestElevator int = 9999 // Forces the calculation of the nearest elevator
	var elevatorWithMoreDropOffTripsMatchingTheRequestedDropOff int
	var maxDropOffLoad int

	// Direction of the trip requested by the user
	tripDirection := ""
	if pickUpFloor < dropOffFloor {
		tripDirection = UP
	} else {
		tripDirection = DOWN
	}

	// Get the nearest elevator going in the same direction than the user wants to go
	for i := range control.Elevators {
		elevatorProximity := int(math.Abs(float64(control.Elevators[i].getFloorNumber() - pickUpFloor)))
		// Our optimal elevator to pick up is the nearest one
		if elevatorProximity < nearestElevator && control.Elevators[i].getDirection() == tripDirection {
			chosenElevator = i
			nearestElevator = elevatorProximity
		}
	}

	// Get the elevator carrying more amount of people wanting to go to the drop-off floor
	for i := range control.Elevators {
		peopleGoingToTheSameDropOffFloor := 0

		// Look at all the assigned trips of this elevator and count how much  people wants
		// to stop in the same dropOffFloor than our user
		for j := range control.Elevators[i].getAssignedTrips() {
			if control.Elevators[i].getAssignedTrip(j).toFloor == dropOffFloor {
				peopleGoingToTheSameDropOffFloor++
			}
		}
		// If there's an elevator going to our same plant and going in the same tripDirection of our current user,
		// it will be chosen as a candidate to be choosen instead of the nearest one
		if peopleGoingToTheSameDropOffFloor > maxDropOffLoad && control.Elevators[i].getDirection() == tripDirection {
			elevatorWithMoreDropOffTripsMatchingTheRequestedDropOff = i
			maxDropOffLoad = peopleGoingToTheSameDropOffFloor
		}
	}

	// If there's an elevator carrying more than 7 people in it, going to our same plant and going in the same direction of our current user,
	// it will choose this one to assign the pick-up task instead of the nearest one
	if maxDropOffLoad > 7 {
		chosenElevator = elevatorWithMoreDropOffTripsMatchingTheRequestedDropOff
	}

	newTrip := TripDetails{
		userID:        userID,
		userAction:    waitingInAFloor,
		elevInFloor:   control.Elevators[chosenElevator].getFloorNumber(),
		fromFloor:     pickUpFloor,
		toFloor:       dropOffFloor,
		tripDirection: tripDirection,
	}

	control.Elevators[chosenElevator].setAssignedTrips(newTrip)

	return control.Elevators[chosenElevator]
}

/***** END OF THE ELEVATOR CONTROL SYSTEM HELPER FUNCTIONS, NOT EXPOSED IN THE INTERFACE *************/

/****************************************************************
 ****************************************************************
 ***                                                          ***
 ***                 ELEVATOR FUNCIONALITY                    ***
 ***                                                          ***
 ****************************************************************
 *****************************************************************/
const UP = "UP"
const DOWN = "DOWN"
const waitingInAFloor = "This elevator will take him there"
const gettingIntoAElevator = "getting into the elevator"
const exitingFromElevator = "exiting from the elevator"

type Elevator interface {
	goToNextFloorInElevatorsTaskList() int
	Step()
	// ELEVATOR GETTERS AND SETTERS
	getFloorNumber() int
	setFloorNumber(floorNumber int)
	getDirection() string
	setDirection(direction string)
	getStepList() StepList
	getStep(step int) TripDetails
	getAssignedTrips() TripQueue
	getAssignedTrip(trip int) TripDetails
	setAssignedTrips(details TripDetails)
	// END OF ELEVATOR GETTERS AND SETTERS
}

// Stores the information about the current status of an elevator
type elevator struct {
	elevID        int       // 0..NUMELEVATORS
	floorNumber   int       // 0..TOPFLOOR: which floor is the elevator in
	topFloor      int       // Top floor of the building
	direction     string    // Up, Down, Stopped
	assignedTrips TripQueue // Queue of assignedTrips assigned to an elevator
	stepList      StepList  // List of steps performed by this Elevator to complete his assignedTrips
}

type TripQueue []TripDetails
type StepList []TripDetails

// Stores the state of an elevator trip
type TripDetails struct {
	userID        string // Not necessary, but added for debugging and tracing purposes
	userAction    string // "waiting for an elevator", "in the elevator", "dropping-off the elevator"
	elevInFloor   int    // Information about the elevator chosen by the system to perform this task
	elevDirection string // Up, Down, Stopped
	fromFloor     int    // Floor where the user presses the pick-up button (0..TOPFLOOR)
	toFloor       int    // Floor where the user wants to go (0..TOPFLOOR)
	tripDirection string // Up, Down, Stopped
}

func NewElevator(i int, topFloor int) Elevator {
	return &elevator{
		elevID:        i,
		floorNumber:   0,
		topFloor:      topFloor,
		direction:     UP,
		assignedTrips: make(TripQueue, 0),
		stepList:      make(StepList, 0),
	}
}

/*
  The elevator will go to the nearest drop-down floor found in the list of his assigned tasks.
  This task is required to be in the same direction as the elevator is currently going.
*/
func (elev *elevator) goToNextFloorInElevatorsTaskList() int {
	minDistance := 9999
	chosenFloor := elev.floorNumber
	aUserAskedToStopInThisFloor := 9999
	nearestFloorWhereAnUserIsWaitingIn := 9999

	// Default behavior: move the elevator to the next floor in the same direction that it is currently moving
	if elev.direction == UP {
		chosenFloor++
	} else {
		chosenFloor--
	}

	// Improvement: search all the assigned trips in the same direction and find ...
	for i := range elev.assignedTrips {
		trip := elev.assignedTrips[i]
		distance := trip.fromFloor - elev.floorNumber
		// ... the nearest assigned trip
		// where a user is waiting for an elevator, and whose trip direction matches the elevator direction
		if distance > 0 &&
			distance < minDistance &&
			trip.userAction == waitingInAFloor {
			nearestFloorWhereAnUserIsWaitingIn = trip.fromFloor
			minDistance = distance
		}
		// ... if there's someone the elevator who wants to stop before, do it
		for i := elev.floorNumber; i < nearestFloorWhereAnUserIsWaitingIn; i++ {
			if i == trip.toFloor && trip.userAction == gettingIntoAElevator {
				aUserAskedToStopInThisFloor = trip.toFloor
			}
		}
	}

	// If some of the improved searches have found a better option than going one by one, jump to the best of them instead of going floor by floor
	if nearestFloorWhereAnUserIsWaitingIn < 9999 || aUserAskedToStopInThisFloor < 9999 {
		// Jump to the nearest floor where there's either an user wants a pick-up and goes in the same direction of the elevator, or
		// an user in the elevator who directly wants to stop in
		if aUserAskedToStopInThisFloor < nearestFloorWhereAnUserIsWaitingIn {
			chosenFloor = aUserAskedToStopInThisFloor
			// If there's not an intermediate stop, jump to the nearest floor where someone is waitint for a pick-up in the same direction
		} else {
			chosenFloor = nearestFloorWhereAnUserIsWaitingIn
		}
	}

	return chosenFloor
}

/**
 * 	Builds the step list of actions performed by an elevator
 */
func (elev *elevator) Step() {
	// Keep moving steps until all the users picked-up by this Elevator
	// have been dropped-off to their destination floor
	for {
		// If the elevator has not assigned trips, stop it.
		// Later we could think in adding a state to the elev
		if len(elev.assignedTrips) == 0 {
			break
		}

		// Complete the current Elevator step
		// For every trip assigned to this Elevator
		for i := 0; i < len(elev.assignedTrips); i++ {
			// Updates the trip info, in order to properly build the list of steps later
			elev.assignedTrips[i].elevInFloor = elev.floorNumber
			elev.assignedTrips[i].elevDirection = elev.direction

			// If the elevator is in the user's floor and goes in the same direction of the requested user's trip, then pick-it up
			if elev.assignedTrips[i].fromFloor == elev.floorNumber {
				elev.assignedTrips[i].userAction = gettingIntoAElevator
				if NotInStepList(elev.stepList, elev.assignedTrips[i]) {
					elev.stepList = append(elev.stepList, elev.assignedTrips[i])
				}
			} else {
				if NotInStepList(elev.stepList, elev.assignedTrips[i]) {
					elev.stepList = append(elev.stepList, elev.assignedTrips[i])
				}
			}

			// If the trip destination floor matches the floor where the Elevator is,
			// drop-off from the Elevator the user who requested the trip,
			// remove this trip from the assignedTrips because it is completed
			// and take note of the step to trace the job that is doing this elevator
			if elev.assignedTrips[i].toFloor == elev.floorNumber {
				if elev.assignedTrips[i].userAction == gettingIntoAElevator {
					elev.assignedTrips[i].userAction = exitingFromElevator
					if NotInStepList(elev.stepList, elev.assignedTrips[i]) {
						elev.stepList = append(elev.stepList, elev.assignedTrips[i])
					}
					elev.assignedTrips = RemoveAssignedTrip(elev.assignedTrips, i)
				}
			}
		}

		// If all the users that wanted to step-out in this floor are out, then move to the next floor
		if noMoreUsersToStepOutInThisFloor(elev.assignedTrips, elev.floorNumber) {
			// Move the elevator to the next step
			if elev.direction == UP {
				if elev.floorNumber == elev.topFloor {
					// If the elevator reached the TOPFLOOR of the building then
					// change downwards and move to next floor
					elev.direction = DOWN
					elev.floorNumber--
				} else {
					// Move to the next floor
					// elev.floorNumber++
					elev.floorNumber = elev.goToNextFloorInElevatorsTaskList()
				}
			} else { // If the elevator is moving down
				if elev.floorNumber == 0 {
					// If the elevator reached the ground floor of the building then
					// change upwards and move to next floor
					elev.direction = UP
					elev.floorNumber++
				} else {
					// Move to the next floor
					// elev.floorNumber--
					elev.floorNumber = elev.goToNextFloorInElevatorsTaskList()
				}
			}
		}
	}
}

/************** ELEVATOR INTERFACE GETTERS AND SETTERS *************/
func (elev *elevator) setAssignedTrips(details TripDetails) {
	elev.assignedTrips = append(elev.assignedTrips, details)
}

func (elev *elevator) getAssignedTrip(tripNumber int) TripDetails {
	trips := elev.getAssignedTrips()
	tripDetails := trips[tripNumber]
	return tripDetails
}

func (elev *elevator) getDirection() string {
	return elev.direction
}

func (elev *elevator) getStepList() StepList {
	return elev.stepList
}

func (elev *elevator) getStep(stepNumber int) TripDetails {
	steps := elev.getStepList()
	stepDetails := steps[stepNumber]
	return stepDetails
}

func (elev *elevator) setFloorNumber(floorNumber int) {
	elev.floorNumber = floorNumber
}

func (elev *elevator) setDirection(direction string) {
	elev.direction = direction
}

func (elev *elevator) getFloorNumber() int {
	return elev.floorNumber
}

func (elev *elevator) getAssignedTrips() TripQueue {
	return elev.assignedTrips
}

func (elev *elevator) changeAssignedTrips() *TripQueue {
	return &elev.assignedTrips
}

/************** END OF ELEVATOR INTERFACE GETTERS AND SETTERS *************/

/********* ELEVATOR INTERNAL HELPER FUNCTIONS, NOT OFFERED IN THE INTERFACE ***************/
// Helps the controller to remove a trip
func RemoveAssignedTrip(assignedTrips []TripDetails, i int) []TripDetails {
	return append(assignedTrips[:i], assignedTrips[i+1:]...)
}

// Avoids duplicates in the stepList
func NotInStepList(stepList []TripDetails, trip TripDetails) bool {
	for i := range stepList {
		if stepList[i].userID == trip.userID && stepList[i].userAction == trip.userAction &&
			stepList[i].fromFloor == trip.fromFloor && stepList[i].toFloor == trip.toFloor {
			return false
		}
	}
	return true
}

/**
 * 	It makes sure all the users that wanted to step-out in this floor are out, then move to the next floor
 */
func noMoreUsersToStepOutInThisFloor(assignedTrips TripQueue, floorNumber int) bool {
	for i := range assignedTrips {
		if assignedTrips[i].toFloor == floorNumber && assignedTrips[i].userAction == gettingIntoAElevator {
			return false
		}
	}
	return true
}

/********* END OF ELEVATOR INTERNAL HELPER FUNCTIONS, NOT OFFERED IN THE INTERFACE ***************/

func main() {
	// Initializes the Elevator Control System, in this case our control will manage
	// 16 elevators in a building of 10 plants
	control := NewElevatorControlSystem(16, 10)

	// List of users pressing the pick up button in the floor they are located, telling which floor they want to go
	control.PickUpButtonWasPushed("User1", 0, 10)
	control.PickUpButtonWasPushed("User2", 2, 5)
	control.PickUpButtonWasPushed("User3", 2, 5)
	control.PickUpButtonWasPushed("User4", 7, 5)
	control.PickUpButtonWasPushed("User5", 6, 4)

	control.PickUpButtonWasPushed("User6", 3, 1)
	control.PickUpButtonWasPushed("User7", 1, 0)
	control.PickUpButtonWasPushed("User8", 9, 7)
	control.PickUpButtonWasPushed("User9", 6, 8)
	control.PickUpButtonWasPushed("User10", 5, 9)
	control.PickUpButtonWasPushed("User11", 10, 6)
	control.PickUpButtonWasPushed("User12", 1, 3)
	control.PickUpButtonWasPushed("User13", 2, 5)
	control.PickUpButtonWasPushed("User14", 3, 2)

	// Will tell the status of every elevator of the system: which floor is in, if it is going up, going down
	// or if it is stopped. It also will tell the list of tasks it has been assigned so far from the list of
	// pick-up requests done by the users so far
	control.Status()

	// Will print the step list of actions performed by the elevators to complete the tasks they have been
	// assigned by the control
	control.Step()
}
