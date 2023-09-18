package request

type MalformedRequest struct {
	Status  int    `json:"-"`
	Message string `json:"request"`
}

func (mr MalformedRequest) Error() string {
	return mr.Message
}

type MalformedRequesOTP struct {
	Status      int    `json:"-"`
	ReleaseTime string `json:"release_time,omitempty"`
}

func (mro MalformedRequesOTP) Error() string {
	return mro.ReleaseTime
}

type MalformedRequestBannedLogin struct {
	Status  int    `json:"-"`
	Message string `json:"release_time"`
}

func (mrbl MalformedRequestBannedLogin) Error() string {
	return mrbl.Message
}
