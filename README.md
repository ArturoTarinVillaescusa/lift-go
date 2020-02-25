## Project objective
Build a Lift control system capable to manage:

 - N lifts
 - M building floors
 
The system objectives are:

 - Pick-up users from the floor they're located
 - Drop-down users in the floor they request
 - Save energy, i.e. non stop in floors where any user wants to stop or where there's no user to be picked-up
 
 ## These are my assumptions for this project

My project is all in one piece on purpose, to facilitate the reviewing task. It could be done in modules, and there are
mark comments indicating which pieces could be extracted from main.go.

I've used this interface:

```bash
type ElevatorControlSystem interface {
	Status()
	Update(elevatorID int, floor int, direction string)
	PickUpButtonWasPushed(userID string, pickUpFloor int, dropOffFloor int)
	Step()
}
```
*NewElevatorControlSystem*

Initializer of an elevator controller system interface with a numberOfElevators and a numberOfFloors

Then initializes an array of elevators to give service to the users. I've added the user name field just for traceability
and debugging purposes. At the end, the execution of Status() makes clearer which actions did take the elevator controller
to satisfy the users requests.

*Status()*

Prints information about the state of the elevators: where is it, if it is going up or down, and what list of tasks have
been assigned to it.

*Update*

It allows to the elevator controler system to intentionally change the floor and the direction of an elevator,
without the intervention of a user pressing the pick-up button. It could be the equivalent to an engineer
using his master key when they are fixing an elevator in a building

*PickUpButtonWasPushed*

For the sake of simplicity I'm assuming that a PickUpButtonWasPushed action is composed of the following actions:

1) A user located in pickUpFloor pushes the button to call an elevator

2) When one of the elevators arrives to the pickUpFloor, the user enters and pushes the button corresponding to the
   floor he wants to go (dropOffFloor)

*Step()*

Builds a list with all the actions in the scenario:

1) Moves the elevators up and down, skipping the floors where there isn't a user step-out request, and also skipping
   the ones where nobody has pushed the button to go there.
      
2) Stores the information related to a user pushing the PickUp button, waiting an elevator to arrive, or enters to it

3) Prints the list of actions performed by the elevators in the whole elevator control system.

4) I'm not implementing the code to prevent an user entering in an elevator whose direction is opposite to his trip
   direction.

I've also added these auxiliary functions, not exposed in the Elevator Controller System Interface, but used internally
by the exposed methods:

*printStepListSimulation*

Used internally by Step()

*chooseTheMostOptimalElevator*

Is the scheduler optimizer. I've chosen this set of rules to assign the best elevator for every pick-up request:

- When a user pushes the PickUp button at any floor, the Elevator chosen to give him service will be chosen
 attending to those rules, in ascending order of importance:

1) The elevatorProximity: this is the difference between the floor where the nearest elevator is and the user's
   pickup floor.

2) The elevator matching more than 7 in the value to 'elevatorWithMoreDropOffTripsMatchingTheRequestedDropOff', assuming
   this is a rule I've made up: the most people going to the same plant the more savings in electricity (maybe
   it is not happening in the real world, but ... I didn't want to leave just a FIFO queue of pickups and drop-offs)

3) The chosen elevation direction must match the direction of the user's requested trip. I.e: if user's trip is going up,
   only an elevator going up can be assigned to this user
   
I also have chosen to offer an Elevator interface, so that the the Elevator Control System has the elevator funcionalities
centralized in an interface.

## System requirements

I've chosen Go to work on this project, so to run it you need to have [installed Golang](https://golang.org/dl/) in your environment too:

Your Go environment variables must be similar to mine:

```sh
arturotarin@QOSMIO-X70B:~/go/src/lift-go
12:16:30 $ echo $GOPATH 
/home/arturotarin/go

arturotarin@QOSMIO-X70B:~/go/src/lift-go
12:16:41 $ echo $GOROOT 
/home/arturotarin/go/go1.12.7
```
 
## How to build the project

If you want to get a runnable to execute the project you must issue this command:

```sh
arturotarin@QOSMIO-X70B:~/go/src/lift-go
15:36:23 $ go build -o main
```

## Run the main application

You can either run it:

* Calling the source:

```sh
arturotarin@QOSMIO-X70B:~/go/src/lift-go
15:42:23 $ go run main.go 
```

*Or running the executable application:

```sh
arturotarin@QOSMIO-X70B:~/go/src/lift-go
15:42:23 $ ./main
```

The result of both will be this output:

```sh
CURRENT STATUS OF THIS ELEVATOR CONTROL SYSTEM: IS MANAGING 10 PLANTS AND 16 ELEVATORS
======================================================================================

* Elevator 0 is in floor 0, going UP, and it has been assigned 14 tasks:

 - User1 is in floor 0 This elevator will take him there. Wants to go to floor 10.
 - User2 is in floor 2 This elevator will take him there. Wants to go to floor 5.
 - User3 is in floor 2 This elevator will take him there. Wants to go to floor 5.
 - User4 is in floor 7 This elevator will take him there. Wants to go to floor 5.
 - User5 is in floor 6 This elevator will take him there. Wants to go to floor 4.
 - User6 is in floor 3 This elevator will take him there. Wants to go to floor 1.
 - User7 is in floor 1 This elevator will take him there. Wants to go to floor 0.
 - User8 is in floor 9 This elevator will take him there. Wants to go to floor 7.
 - User9 is in floor 6 This elevator will take him there. Wants to go to floor 8.
 - User10 is in floor 5 This elevator will take him there. Wants to go to floor 9.
 - User11 is in floor 10 This elevator will take him there. Wants to go to floor 6.
 - User12 is in floor 1 This elevator will take him there. Wants to go to floor 3.
 - User13 is in floor 2 This elevator will take him there. Wants to go to floor 5.
 - User14 is in floor 3 This elevator will take him there. Wants to go to floor 2.

* Elevator 1 is in floor 0, stopped.
* Elevator 2 is in floor 0, stopped.
* Elevator 3 is in floor 0, stopped.
* Elevator 4 is in floor 0, stopped.
* Elevator 5 is in floor 0, stopped.
* Elevator 6 is in floor 0, stopped.
* Elevator 7 is in floor 0, stopped.
* Elevator 8 is in floor 0, stopped.
* Elevator 9 is in floor 0, stopped.
* Elevator 10 is in floor 0, stopped.
* Elevator 11 is in floor 0, stopped.
* Elevator 12 is in floor 0, stopped.
* Elevator 13 is in floor 0, stopped.
* Elevator 14 is in floor 0, stopped.
* Elevator 15 is in floor 0, stopped.

STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 0, going UP. User2 pressed the pick-up button in floor 2 and wants to go to floor 5. This elevator will take him there.
Floor 0, going UP. User3 pressed the pick-up button in floor 2 and wants to go to floor 5. This elevator will take him there.
Floor 0, going UP. User4 pressed the pick-up button in floor 7 and wants to go to floor 5. This elevator will take him there.
Floor 0, going UP. User5 pressed the pick-up button in floor 6 and wants to go to floor 4. This elevator will take him there.
Floor 0, going UP. User6 pressed the pick-up button in floor 3 and wants to go to floor 1. This elevator will take him there.
Floor 0, going UP. User7 pressed the pick-up button in floor 1 and wants to go to floor 0. This elevator will take him there.
Floor 0, going UP. User8 pressed the pick-up button in floor 9 and wants to go to floor 7. This elevator will take him there.
Floor 0, going UP. User9 pressed the pick-up button in floor 6 and wants to go to floor 8. This elevator will take him there.
Floor 0, going UP. User10 pressed the pick-up button in floor 5 and wants to go to floor 9. This elevator will take him there.
Floor 0, going UP. User11 pressed the pick-up button in floor 10 and wants to go to floor 6. This elevator will take him there.
Floor 0, going UP. User12 pressed the pick-up button in floor 1 and wants to go to floor 3. This elevator will take him there.
Floor 0, going UP. User13 pressed the pick-up button in floor 2 and wants to go to floor 5. This elevator will take him there.
Floor 0, going UP. User14 pressed the pick-up button in floor 3 and wants to go to floor 2. This elevator will take him there.
Floor 1, going UP. User7 is getting into the elevator.
Floor 1, going UP. User12 is getting into the elevator.
Floor 2, going UP. User2 is getting into the elevator.
Floor 2, going UP. User3 is getting into the elevator.
Floor 2, going UP. User13 is getting into the elevator.
Floor 3, going UP. User6 is getting into the elevator.
Floor 3, going UP. User12 is exiting from the elevator in floor 3.
Floor 3, going UP. User14 is getting into the elevator.
Floor 5, going UP. User2 is exiting from the elevator in floor 5.
Floor 5, going UP. User10 is getting into the elevator.
Floor 5, going UP. User13 is exiting from the elevator in floor 5.
Floor 5, going UP. User3 is exiting from the elevator in floor 5.
Floor 6, going UP. User5 is getting into the elevator.
Floor 6, going UP. User9 is getting into the elevator.
Floor 7, going UP. User4 is getting into the elevator.
Floor 8, going UP. User9 is exiting from the elevator in floor 8.
Floor 9, going UP. User8 is getting into the elevator.
Floor 9, going UP. User10 is exiting from the elevator in floor 9.
Floor 10, going UP. User1 is exiting from the elevator in floor 10.
Floor 10, going UP. User11 is getting into the elevator.
Floor 7, going DOWN. User8 is exiting from the elevator in floor 7.
Floor 6, going DOWN. User11 is exiting from the elevator in floor 6.
Floor 5, going DOWN. User4 is exiting from the elevator in floor 5.
Floor 4, going DOWN. User5 is exiting from the elevator in floor 4.
Floor 2, going DOWN. User14 is exiting from the elevator in floor 2.
Floor 1, going DOWN. User6 is exiting from the elevator in floor 1.
Floor 0, going DOWN. User7 is exiting from the elevator in floor 0.

Elevator 1 Performed This Step List
----------------------------------------

Elevator 2 Performed This Step List
----------------------------------------

Elevator 3 Performed This Step List
----------------------------------------

Elevator 4 Performed This Step List
----------------------------------------

Elevator 5 Performed This Step List
----------------------------------------

Elevator 6 Performed This Step List
----------------------------------------

Elevator 7 Performed This Step List
----------------------------------------

Elevator 8 Performed This Step List
----------------------------------------

Elevator 9 Performed This Step List
----------------------------------------

Elevator 10 Performed This Step List
----------------------------------------

Elevator 11 Performed This Step List
----------------------------------------

Elevator 12 Performed This Step List
----------------------------------------

Elevator 13 Performed This Step List
----------------------------------------

Elevator 14 Performed This Step List
----------------------------------------

Elevator 15 Performed This Step List
----------------------------------------
```

You can try modifying the main() function in the main.go package to get different results.

## Testing the application

Run Benchmark test:

```
arturotarin@QOSMIO-X70B:~/go/src/lift-go
15:44:23 $ go test -bench=.

goos: linux
goarch: amd64
pkg: ArturoGo
BenchmarkElevatorControlSystem-8   	100
10000
200000
  200000	      6117 ns/op
PASS
ok  	lift-go	1.385s
```

Run parallel test executions:

```
arturotarin@QOSMIO-X70B:~/go/src/lift-go
15:47:58 $ go test

STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 1, going UP. User1 is exiting from the elevator in floor 1.


STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 1, going UP. User1 is exiting from the elevator in floor 1.
Floor 2, going UP. User2 is getting into the elevator.
Floor 10, going UP. User2 is exiting from the elevator in floor 10.


STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 1, going UP. User1 is exiting from the elevator in floor 1.
Floor 2, going UP. User2 is getting into the elevator.
Floor 10, going UP. User2 is exiting from the elevator in floor 10.

Elevator 1 Performed This Step List
----------------------------------------
Floor 0, going UP. User3 pressed the pick-up button in floor 4 and wants to go to floor 10. This elevator will take him there.
Floor 4, going UP. User3 is getting into the elevator.
Floor 10, going UP. User3 is exiting from the elevator in floor 10.


STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 1, going UP. User1 is exiting from the elevator in floor 1.
Floor 2, going UP. User2 is getting into the elevator.
Floor 10, going UP. User2 is exiting from the elevator in floor 10.

Elevator 1 Performed This Step List
----------------------------------------
Floor 0, going UP. User3 pressed the pick-up button in floor 4 and wants to go to floor 10. This elevator will take him there.
Floor 4, going UP. User3 is getting into the elevator.
Floor 10, going UP. User3 is exiting from the elevator in floor 10.

Elevator 2 Performed This Step List
----------------------------------------
Floor 0, going UP. User4 pressed the pick-up button in floor 9 and wants to go to floor 10. This elevator will take him there.
Floor 9, going UP. User4 is getting into the elevator.
Floor 10, going UP. User4 is exiting from the elevator in floor 10.


STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 1, going UP. User1 is exiting from the elevator in floor 1.
Floor 2, going UP. User2 is getting into the elevator.
Floor 10, going UP. User2 is exiting from the elevator in floor 10.
Floor 9, going DOWN. User5 pressed the pick-up button in floor 10 and wants to go to floor 9. This elevator will take him there.
Floor 10, going DOWN. User5 is getting into the elevator.
Floor 9, going DOWN. User5 is exiting from the elevator in floor 9.

Elevator 1 Performed This Step List
----------------------------------------
Floor 0, going UP. User3 pressed the pick-up button in floor 4 and wants to go to floor 10. This elevator will take him there.
Floor 4, going UP. User3 is getting into the elevator.
Floor 10, going UP. User3 is exiting from the elevator in floor 10.

Elevator 2 Performed This Step List
----------------------------------------
Floor 0, going UP. User4 pressed the pick-up button in floor 9 and wants to go to floor 10. This elevator will take him there.
Floor 9, going UP. User4 is getting into the elevator.
Floor 10, going UP. User4 is exiting from the elevator in floor 10.

...
... 
... (The testing runs several different combinations, check them in your output screen)
...
...

STEP LIST OF OUR SYSTEM OF 10 FLOORS AND 16 ELEVATORS
=====================================================

Elevator 0 Performed This Step List
----------------------------------------
Floor 0, going UP. User1 is getting into the elevator.
Floor 1, going UP. User1 is exiting from the elevator in floor 1.
Floor 2, going UP. User2 is getting into the elevator.
Floor 10, going UP. User2 is exiting from the elevator in floor 10.
Floor 9, going DOWN. User5 pressed the pick-up button in floor 10 and wants to go to floor 9. This elevator will take him there.
Floor 10, going DOWN. User5 is getting into the elevator.
Floor 9, going DOWN. User5 is exiting from the elevator in floor 9.
Floor 8, going DOWN. User13 pressed the pick-up button in floor 4 and wants to go to floor 3. This elevator will take him there.
Floor 4, going DOWN. User13 is getting into the elevator.
Floor 3, going DOWN. User13 is exiting from the elevator in floor 3.
Floor 2, going DOWN. User17 pressed the pick-up button in floor 3 and wants to go to floor 0. This elevator will take him there.
Floor 3, going DOWN. User17 is getting into the elevator.
Floor 0, going DOWN. User17 is exiting from the elevator in floor 0.
Floor 1, going UP. User19 pressed the pick-up button in floor 2 and wants to go to floor 7. This elevator will take him there.
Floor 2, going UP. User19 is getting into the elevator.
Floor 7, going UP. User19 is exiting from the elevator in floor 7.

Elevator 1 Performed This Step List
----------------------------------------
Floor 0, going UP. User3 pressed the pick-up button in floor 4 and wants to go to floor 10. This elevator will take him there.
Floor 4, going UP. User3 is getting into the elevator.
Floor 10, going UP. User3 is exiting from the elevator in floor 10.
Floor 9, going DOWN. User7 is getting into the elevator.
Floor 3, going DOWN. User7 is exiting from the elevator in floor 3.
Floor 2, going DOWN. User12 pressed the pick-up button in floor 3 and wants to go to floor 1. This elevator will take him there.
Floor 3, going DOWN. User12 is getting into the elevator.
Floor 1, going DOWN. User12 is exiting from the elevator in floor 1.

Elevator 2 Performed This Step List
----------------------------------------
Floor 0, going UP. User4 pressed the pick-up button in floor 9 and wants to go to floor 10. This elevator will take him there.
Floor 9, going UP. User4 is getting into the elevator.
Floor 10, going UP. User4 is exiting from the elevator in floor 10.
Floor 9, going DOWN. User14 pressed the pick-up button in floor 7 and wants to go to floor 7. This elevator will take him there.
Floor 7, going DOWN. User14 is getting into the elevator.
Floor 7, going DOWN. User14 is exiting from the elevator in floor 7.

Elevator 3 Performed This Step List
----------------------------------------
Floor 0, going UP. User6 pressed the pick-up button in floor 5 and wants to go to floor 10. This elevator will take him there.
Floor 5, going UP. User6 is getting into the elevator.
Floor 10, going UP. User6 is exiting from the elevator in floor 10.
Floor 9, going DOWN. User15 pressed the pick-up button in floor 10 and wants to go to floor 6. This elevator will take him there.
Floor 10, going DOWN. User15 is getting into the elevator.
Floor 6, going DOWN. User15 is exiting from the elevator in floor 6.

Elevator 4 Performed This Step List
----------------------------------------
Floor 0, going UP. User8 pressed the pick-up button in floor 1 and wants to go to floor 10. This elevator will take him there.
Floor 1, going UP. User8 is getting into the elevator.
Floor 10, going UP. User8 is exiting from the elevator in floor 10.
Floor 9, going DOWN. User18 is getting into the elevator.
Floor 1, going DOWN. User18 is exiting from the elevator in floor 1.

Elevator 5 Performed This Step List
----------------------------------------
Floor 0, going UP. User9 pressed the pick-up button in floor 3 and wants to go to floor 10. This elevator will take him there.
Floor 3, going UP. User9 is getting into the elevator.
Floor 10, going UP. User9 is exiting from the elevator in floor 10.

Elevator 6 Performed This Step List
----------------------------------------
Floor 0, going UP. User10 pressed the pick-up button in floor 2 and wants to go to floor 10. This elevator will take him there.
Floor 2, going UP. User10 is getting into the elevator.
Floor 10, going UP. User10 is exiting from the elevator in floor 10.

Elevator 7 Performed This Step List
----------------------------------------
Floor 0, going UP. User11 pressed the pick-up button in floor 5 and wants to go to floor 10. This elevator will take him there.
Floor 5, going UP. User11 is getting into the elevator.
Floor 10, going UP. User11 is exiting from the elevator in floor 10.

Elevator 8 Performed This Step List
----------------------------------------
Floor 0, going UP. User16 pressed the pick-up button in floor 1 and wants to go to floor 4. This elevator will take him there.
Floor 1, going UP. User16 is getting into the elevator.
Floor 4, going UP. User16 is exiting from the elevator in floor 4.

Elevator 9 Performed This Step List
----------------------------------------
Floor 0, going UP. User20 is getting into the elevator.
Floor 5, going UP. User20 is exiting from the elevator in floor 5.

PASS
ok  	lift-go	0.092s
```
