package main

type MeetingRole struct {
	Role        string
	DisplayName string
}

var AgendaTemplate []MeetingItem
var RolesTemplate []MeetingRole

const secondsPerMinute int = 60

func init() {

	rolesTemplateConfig := []string{
		"SAA", "SAA",
		"President", "President",
		"TMD", "TMD",
		"GE", "GE",
		"AhCounter", "AhCounter",
		"Timer", "Timer",
		"Grammarian", "Grammarian",
		"ShareMaster", "ShareMaster",
		"TTM", "TTM",
		"TTIE", "TTIE",
		"Speaker", "Speaker 1",
		"Speaker", "Speaker 2",
		"Speaker", "Speaker 3",
		"IE", "IE 1",
		"IE", "IE 2",
		"IE", "IE 3",
		"VPE", "VPE",
	}

	RolesTemplate := make([]MeetingRole, len(rolesTemplateConfig)/2)
	for i, _ := range rolesTemplateConfig {
		if i%2 != 0 {
			continue
		}
		RolesTemplate[i/2] = MeetingRole{rolesTemplateConfig[i], rolesTemplateConfig[i+1]}
	}

	AgendaTemplate = []MeetingItem{
		MeetingItem{Role: "SAA", Duration: 2 * secondsPerMinute},
		MeetingItem{Role: "President", Duration: secondsPerMinute},
		MeetingItem{Role: "TMD", Duration: secondsPerMinute},
		MeetingItem{Role: "GE", Duration: 3 * secondsPerMinute},
		MeetingItem{Role: "AhCounter", Duration: secondsPerMinute},
		MeetingItem{Role: "Timer", Duration: secondsPerMinute},
		MeetingItem{Role: "Grammarian", Duration: 2 * secondsPerMinute},
		MeetingItem{Role: "TTM", Duration: 15 * secondsPerMinute},
		MeetingItem{Role: "TTIE", Duration: 5 * secondsPerMinute},
		MeetingItem{Role: "Speaker", Duration: 7 * secondsPerMinute},
		MeetingItem{Role: "Speaker", Duration: 7 * secondsPerMinute},
		MeetingItem{Role: "Speaker", Duration: 7 * secondsPerMinute},
		MeetingItem{Role: "IE", Duration: 3 * secondsPerMinute},
		MeetingItem{Role: "IE", Duration: 3 * secondsPerMinute},
		MeetingItem{Role: "IE", Duration: 3 * secondsPerMinute},
		MeetingItem{Role: "AhCounter", Duration: secondsPerMinute},
		MeetingItem{Role: "Timer", Duration: secondsPerMinute},
		MeetingItem{Role: "Grammarian", Duration: 2 * secondsPerMinute},
		MeetingItem{Role: "GE", Duration: 5 * secondsPerMinute},
		MeetingItem{Role: "President", Duration: 2 * secondsPerMinute},
	}
}

func getMeetingTemplate() {

}
