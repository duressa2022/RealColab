package domain

type HomeInformation struct {
	ProfileUrl             string                 `json:"profileUrl" bson:"profileUrl"`
	FirstName              string                 `json:"firstName" bson:"firstName"`
	LastName               string                 `json:"lastName" bson:"lastName"`
	Rating                 UserRating             `json:"rating" bson:"rating"`
	TaskData               TaskInformation        `json:"taskData" bson:"taskData"`
	RecentlyCompletedTasks []*PrivateTask         `json:"recentlyCompletedTasks" bson:"recentlyCompletedTasks"`
	UpcomingTasks          []*PrivateTask         `json:"upcomingTasks" bson:"upcomingTasks"`
	UnReadNotifications    map[string]interface{} `json:"unReadNotifications" bson:"unReadNotifications"`
}
