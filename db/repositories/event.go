package repositories

import "order/db/models"

type EventSourcingRepository interface {
	FindOne(eventId string) (*models.EventSource, error)
	FindByEventName(eventName string) (event *models.EventSource, err error)
	InsertEvent(event *models.EventSource) error
}

//FindOne find event by id
func (repo *defaultRepository) FindOne(eventId string) (event *models.EventSource, err error) {
	return
}

//FindByEventName find event by name
func (repo *defaultRepository) FindByEventName(eventName string) (event *models.EventSource, err error) {
	return
}

//InsertEvent insert a event
func (repo *defaultRepository) InsertEvent(event *models.EventSource) error {
	return nil
}
