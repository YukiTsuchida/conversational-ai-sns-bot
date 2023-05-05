package cmd

type Type string

const (
	PostMessage          = Type("PostMessage")
	GetMyMessages        = Type("GetMyMessages")
	GetOtherMessages     = Type("GetOtherMessages")
	SearchMessage        = Type("SearchMessage")
	GetMyProfile         = Type("GetMyProfile")
	GetOthersProfile     = Type("GetOthersProfile")
	UpdateMyProfile      = Type("UpdateMyProfile")
	PostActionPurposeLog = Type("PostActionPurposeLog")
)
