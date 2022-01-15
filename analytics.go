package register

/*
type AnalyticsFunc func(context.Context, AnalyticsEvent) error

type AnalyticsFunction struct {
	reg *FunctionRegistrar
	fn  AnalyticsFunc
}

// AnalyticsEvent is the payload of an Analytics log event.
type AnalyticsEvent struct {
        EventDimensions []EventDimensions `json:"eventDim"`
        UserDimensions  interface{}       `json:"userDim"`
}

// EventDimensions holds Analytics event dimensions.
type EventDimensions struct {
        Name                    string      `json:"name"`
        Date                    string      `json:"date"`
        TimestampMicros         string      `json:"timestampMicros"`
        PreviousTimestampMicros string      `json:"previousTimestampMicros"`
        Params                  interface{} `json:"params"`
}

*/

/*

{
    "eventDim": [ // Contains a single event
        {
            "date": "20090213",
            "name": "screen_view",
            "params": {
                "firebase_conversion": {
                    "intValue": "1"
                },
                "firebase_event_origin": {
                    "stringValue": "auto"
                },
                "firebase_previous_class": {
                    "stringValue": "MainActivity"
                },
                "firebase_previous_id": {
                    "intValue": "1928209043426257906"
                },
                "firebase_previous_screen": {
                    "stringValue": "id-D-D"
                },
                "firebase_screen": {
                    "stringValue": "id-C-C"
                },
                "firebase_screen_class": {
                    "stringValue": "MainActivity"
                },
                "firebase_screen_id": {
                    "intValue": "1234567890000"
                }
            },
            "previousTimestampMicros": "1234567890000",
            "timestampMicros": "1234567890000"
        }
    ],
    "userDim": {
        // A UserDimensions object
    }
}
*/
