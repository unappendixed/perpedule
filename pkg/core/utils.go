package core

import ics "github.com/arran4/golang-ical"

const VendorPrefix string = "X-PPD-"
const PPDType string = VendorPrefix + "ITEMTYPE"

// vendored from ics library
const ICalTimestampFormatUtc string = "20060102T150405Z"

const (
	PPDTypeProject   = "PROJECT"
	PPDTypeTimeblock = "TIMEBLOCK"
)

const (
    ICSTokenTodo = "VTODO"
    ICSTokenEvent = "VEVENT"
    ICSTokenAlarm = "VALARM"
    ICSTokenJournal = "VJOURNAL"
)

func GetRelatedToId(c ics.GeneralComponent) string {
	prop := c.GetProperty(ics.ComponentProperty(ics.PropertyRelatedTo))
	if prop != nil {
		return prop.Value
	}

	return ""
}
