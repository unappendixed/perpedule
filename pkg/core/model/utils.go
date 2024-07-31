package core

import ics "github.com/arran4/golang-ical"

const VendorPrefix string = "X-PPD-"
const PPDType string = VendorPrefix + "ITEMTYPE"

const (
    PPDTypeProject = "PROJECT"
    PPDTypeTimeblock = "TIMEBLOCK"
)

func GetRelatedToId(c ics.ComponentBase) string {
    prop := c.GetProperty(ics.ComponentProperty(ics.PropertyRelatedTo))
    if prop != nil {
        return prop.Value
    }

    return ""
}
