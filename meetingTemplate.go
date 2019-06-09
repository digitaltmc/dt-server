package main

type MeetingRole struct {
	RoleName    string
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
		MeetingItem{RoleName: "SAA", Duration: "2"},
		MeetingItem{RoleName: "President", Duration: "1"},
		MeetingItem{RoleName: "TMD", Duration: "1"},
		MeetingItem{RoleName: "GE", Duration: "3"},
		MeetingItem{RoleName: "AhCounter", Duration: "1"},
		MeetingItem{RoleName: "Timer", Duration: "1"},
		MeetingItem{RoleName: "Grammarian", Duration: "2"},
		MeetingItem{RoleName: "TTM", Duration: "15"},
		MeetingItem{RoleName: "TTIE", Duration: "5"},
		MeetingItem{RoleName: "Speaker", Duration: "7"},
		MeetingItem{RoleName: "Speaker", Duration: "7"},
		MeetingItem{RoleName: "Speaker", Duration: "7"},
		MeetingItem{RoleName: "IE", Duration: "3"},
		MeetingItem{RoleName: "IE", Duration: "3"},
		MeetingItem{RoleName: "IE", Duration: "3"},
		MeetingItem{RoleName: "AhCounter", Duration: "1"},
		MeetingItem{RoleName: "Timer", Duration: "1"},
		MeetingItem{RoleName: "Grammarian", Duration: "2"},
		MeetingItem{RoleName: "GE", Duration: "5"},
		MeetingItem{RoleName: "President", Duration: "2"},
	}
}

func getMeetingTemplate() {

}
