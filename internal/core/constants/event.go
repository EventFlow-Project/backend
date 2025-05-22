package constants

type EventStatus string

const (
	EventStatusComingUp EventStatus = "Предстоит"
	EventStatusUnderway EventStatus = "Идёт"
	EventStatusHeld     EventStatus = "Прошло"
)

type EventModerationStatus string

const (
	EventModerationStatusPending  EventModerationStatus = "На модерации"
	EventModerationStatusApproved EventModerationStatus = "Одобрено"
	EventModerationStatusRejected EventModerationStatus = "Отклонено"
)

type EventTag string

const (
	EventTagConference EventTag = "Конференция"
	EventTagWorkshop   EventTag = "Воркшоп"
	EventTagMeetup     EventTag = "Митап"
	EventTagHackathon  EventTag = "Хакатон"
	EventTagExhibition EventTag = "Выставка"
	EventTagFestival   EventTag = "Фестиваль"
	EventTagOnline     EventTag = "Онлайн"
	EventTagOffline    EventTag = "Оффлайн"
)
