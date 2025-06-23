package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kitessafikadu/kevent/internal/database"
)

func (app *application) createEvent (c *gin.Context){
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func (app *application) getAllEvents(c *gin.Context){
	events, err:=app.models.Events.GetAll()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (app *application) getEvent(c *gin.Context){
	id, err:=strconv.Atoi(c.Param("id"))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	}
	event,err:= app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
	}
	c.JSON(http.StatusOK, event)
}

	func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	updatedEvent := &database.Event{}
	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id
	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// Add attendees to an event
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.Get(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd,err:= app.models.Users.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if userToAdd == nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err:= app.models.Attendees.GetByEventAndUser(event.Id, user.Id)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
			return
		}

		if existingAttendee != nil{
			c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
			return
		}

		attendee := database.Attendee{
			EventId: event.Id,
			UserId: userToAdd.Id,
		}

		_,err= app.models.Attendees.Insert(&attendee)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		}

		c.JSON(http.StatusCreated, attendee)
}

// Get Attendees for an event
func (app *application) getAttendeesForEvent(c *gin.Context){
	id, err:=strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return 
	}
	users,err:=app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for events"})
		return
	}
	c.JSON(http.StatusOK, users)
}