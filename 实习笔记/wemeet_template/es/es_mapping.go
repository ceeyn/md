package es

import "encoding/json"

const (
	WebinarItineraryIndex = "webinar_itinerary_"
	WebinarIntroductionIndex   = "webinar_introduction_"
)

const esItineraryMapping = `
{
    "mappings": {
        "doc" : {
			"properties": {
				"meeting_id": {
					"type": "keyword"
				},
				"app_id": {
					"type": "keyword"
				},
				"app_uid": {
					"type": "keyword"
				},
				"itinerary_name": {
					"type": "text"
				},
				"itinerary_introduction": {
					"type": "text"
				},
				"create_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
				"update_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
			    "meeting_order_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_order_start_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_order_end_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_subject": {
					"type": "text"
				}
			}
        }
    }
}`


const esIntroductionMapping = `
{
    "mappings": {
        "doc" : {
			"properties": {
				"meeting_id": {
					"type": "keyword"
				},
				"app_id": {
					"type": "keyword"
				},
				"app_uid": {
					"type": "keyword"
				},
				"meeting_introduction": {
					"type": "text"
				},
				"create_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
				"update_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_order_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_order_start_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_order_end_time": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
				},
                "meeting_subject": {
					"type": "text"
				}
			}
        }
    }
}`


// Itinerary ...
type Itinerary struct {
	MeetingId            				  string `json:"meeting_id"`                          // 会议id
	AppId                				  string `json:"app_id"`
	AppUid               				  string `json:"app_uid"`
	ItineraryName                         string `json:"itinerary_name"`                       //日程名称
	ItineraryIntroduction                 string `json:"itinerary_introduction"`               //日程介绍
	CreateTime                            string `json:"create_time,omitempty"`
	UpdateTime                            string `json:"update_time,omitempty"`
	MeetingOrderTime                      string `json:"meeting_order_time,omitempty"`                   //会议预定时间
	MeetingOrderStartTime                 string `json:"meeting_order_start_time,omitempty"`             //会议预定开始时间
	MeetingOrderEndTime                   string `json:"meeting_order_end_time,omitempty"`               //会议预定开始时间
	MeetingSubject                        string `json:"meeting_subject,omitempty"`                      //会议主题
}

// Introduction ...
type Introduction struct {
	MeetingId            				  string `json:"meeting_id"`                          // 会议id
	AppId                				  string `json:"app_id"`
	AppUid               				  string `json:"app_uid"`
	MeetingIntroduction                   string `json:"meeting_introduction"`                 //会议介绍
	CreateTime                            string `json:"create_time,omitempty"`
	UpdateTime                            string `json:"update_time,omitempty"`
	MeetingOrderTime                      string `json:"meeting_order_time,omitempty"`                   //会议预定时间
	MeetingOrderStartTime                 string `json:"meeting_order_start_time,omitempty"`             //会议预定开始时间
	MeetingOrderEndTime                   string `json:"meeting_order_end_time,omitempty"`               //会议预定开始时间
	MeetingSubject                        string `json:"meeting_subject,omitempty"`                      //会议主题
}

// String ...
func (h Itinerary) String() string {
	str, _ := json.Marshal(h)
	return string(str)
}

// String ...
func (h Introduction) String() string {
	str, _ := json.Marshal(h)
	return string(str)
}