package splunk

var userCases = []struct {
	description string
	username    string
	password    string
}{
	{
		"aplhanum test",
		"sdajf1243EF",
		"sdgaf12342SVD",
	},
	{
		"special characters test",
		"@!#$asf12",
		"!@FAF@$",
	},
	{
		"check for white spaces",
		"  f  a424rqw3r$   ",
		" s  dafaf214@RQWE  ",
	},
}

var basicAuthCases = []struct {
	description string
	user        string
	password    string
	host        string
	port        string
	want        string
}{
	{
		"test the default settings with nill user",
		"",
		"",
		"",
		"",
		"401",
	},
	{
		"test the defaults with not null user",
		"admin",
		"changeme",
		"",
		"",
		"200",
	},
	{
		"test auth for custom host",
		"admin",
		"changeme",
		"localhost",
		"8089",
		"200",
	},
	{
		"test auth for custom host",
		"admin",
		"changeme",
		"12@#$svnsdn",
		"8089",
		"Invalid hostname :12@#$svnsdn provided",
	},
	{
		"test auth for custom host",
		"admin",
		"changeme",
		"localhost",
		"asfsv",
		"Invalid port : asfsv provided",
	},
	{
		"test auth for custom host",
		"admin",
		"changeme",
		"localhost",
		"65536",
		"Invalid port : 65536 provided",
	},
}

var searchCases = []struct {
	description string
	user        string
	password    string
	search      string
	want        string
}{
	{
		"perform simple search to list all the splunk_servers",
		"admin",
		"changeme",
		"index=_internal | stats count by splunk_server",
		"",
	},
	{
		"perform simple search w/ search keyword",
		"admin",
		"changeme",
		"search index=_internal | stats count by splunk_server",
		"",
	},
	{
		"perform a pipe prefixed search",
		"admin",
		"changeme",
		"| rest /services/apps/local | search disabled=0 | table label version",
		"",
	},
}
