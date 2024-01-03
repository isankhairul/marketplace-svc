package message

// Message wrapper.
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Common message
var (
	SuccessMsg       = Message{Code: 200, Message: "Success"}
	SuccessNoDataMsg = Message{Code: 204, Message: "No Data Found"}
	ValidationError  = Message{Code: 400, Message: "validation error"}
	ErrReq           = Message{Code: 400, Message: "required field"}
	FailedMsg        = Message{Code: 400, Message: "Failed"}
	ErrNoAuth        = Message{Code: 401, Message: "No Authorization"}
	ErrDB            = Message{Code: 500, Message: "Error has been occurred while processing database request"}
	ErrES            = Message{Code: 500, Message: "Error has been occurred while processing elasticsearch request"}
)

// Specific message
var (
	InvalidPhoneFormat                = Message{Code: 400, Message: "Format nomor HP tidak sesuai"}
	PhoneRequired                     = Message{Code: 400, Message: "No HP harus diisi"}
	PhoneRegistered                   = Message{Code: 400, Message: "No HP sudah terdaftar"}
	OTPPhoneNotRegistered             = Message{Code: 400, Message: "Nomor ponsel Anda tidak terdaftar"}
	CodeRequired                      = Message{Code: 400, Message: "Kode harus diisi"}
	CodeInvalid                       = Message{Code: 400, Message: "Kode yang anda masukkan salah"}
	CodeExpired                       = Message{Code: 400, Message: "Kode yang anda masukkan sudah kadaluarsa"}
	StrNoRegistered                   = Message{Code: 400, Message: "Nomor STR sudah terdaftar"}
	SIPTypeRequired                   = Message{Code: 400, Message: "Tipe SIP harus diisi"}
	SIPTypeInvalid                    = Message{Code: 400, Message: "Tipe SIP yang anda masukkan salah"}
	UnableGenerateCode                = Message{Code: 400, Message: "Gagal menghasilkan kode"}
	SpecialtyRequired                 = Message{Code: 400, Message: "Spesialisasi harus diisi"}
	ProfessionRequired                = Message{Code: 400, Message: "Profesi harus diisi"}
	UserApprovalRejected              = Message{Code: 400, Message: "Status pendaftaran ditolak"}
	PhoneNotActivated                 = Message{Code: 400, Message: "No HP anda belum di aktivasi"}
	PINRequired                       = Message{Code: 400, Message: "PIN harus diisi"}
	OldPINRequired                    = Message{Code: 400, Message: "PIN lama harus diisi"}
	TypeRequired                      = Message{Code: 400, Message: "Type harus diisi"}
	PhoneNotRegistered                = Message{Code: 400, Message: "Akun anda belum terdaftar, Silahkan Registrasi"}
	WrongPin                          = Message{Code: 400, Message: "Pin anda salah"}
	WrongOldPin                       = Message{Code: 400, Message: "Pin lama anda salah"}
	BannedLogin                       = Message{Code: 400, Message: "Anda baru bisa melakukan aktivitas login dalam {{time_release}}"}
	MedicalRecords                    = Message{Code: 400, Message: "Anda masih mempunyai Catatan Medis yang belum terisi. Harap melengkapi Catatan Medis yang tersisa sebelum menonaktifkan toggle"}
	UserAlreadyDeleted                = Message{Code: 400001, Message: "Akun Anda sudah tidak aktif, Hubungi Customer Service Kami untuk mengaktifkan kembali."}
	ImageRequired                     = Message{Code: 400, Message: "Image harus diisi"}
	FailedExractToken                 = Message{Code: 400, Message: "Fail Extract Token"}
	IdTransRequired                   = Message{Code: 400, Message: "ID Transaksi harus diisi"}
	ReportProblemUid                  = Message{Code: 400, Message: "Report Problem Uid harus diisi"}
	ReportProblemValue                = Message{Code: 400, Message: "Report Problem Value harus diisi"}
	DoctorNotFound                    = Message{Code: 400, Message: "Data Doctor Not Found"}
	MediaPathRequired                 = Message{Code: 400, Message: "Media Path harus diisi"}
	PlatformRequired                  = Message{Code: 400, Message: "Platform harus diisi"}
	FcmTokenRequired                  = Message{Code: 400, Message: "FCM Token harus diisi"}
	AccessTokenRequired               = Message{Code: 400, Message: "Access Token harus diisi"}
	RefreshTokenRequired              = Message{Code: 400, Message: "Refresh Token harus diisi"}
	WaitingOTPRelease                 = Message{Code: 400002, Message: "Silahkan tunggu selama {{time_release}}"}
	FailedGenerateToken               = Message{Code: 400, Message: "Fail generate token"}
	UserPending                       = Message{Code: 400, Message: "Akun anda masih pending"}
	AddressNotFound                   = Message{Code: 400, Message: "Address Not Found"}
	MerchantProductSearchSortNotFound = Message{Code: 400, Message: "The sort only distance or fulfill"}
)

var TelErrUserNotFound = Message{Code: 34000, Message: "Not found"}
var ErrDataExists = Message{Code: 34001, Message: "Data already exists"}
var ErrBadRouting = Message{Code: 34002, Message: "Inconsistent mapping between route and handler"}
var ErrInternalError = Message{Code: 34003, Message: "Error has been occurred while processing request"}
var ErrInvalidHeader = Message{Code: 34005, Message: "Invalid header"}
var ErrNoData = Message{Code: 34005, Message: "Data is not found"}
var ErrSaveData = Message{Code: 34005, Message: "Data cannot be saved, please check your request"}
var ErrPageNotFound = Message{Code: 39404, Message: "Page not found"}
var ErrReqParam = Message{Code: 4000, Message: "Invalid Request Parameter(s)"}
var ErrNoIndexName = Message{Code: 34005, Message: "Data index-name not in config"}
var ErrInvalidReqFilter = Message{Code: 34005, Message: "Invalid json format on 'filter' parameter"}

var AuthenticationFailed = Message{Code: 34006, Message: "JWT token is invalid"}
var UnauthorizedError = Message{Code: 34007, Message: "No authorization token was found"}
var NotAllowed = Message{Code: 405213, Message: "Request Not Allowed"}
var UnauthorizedTokenDevice = Message{Code: 401320, Message: "Token Device Unauthorized"}
var SessionLoginExpired = Message{Code: 401231, Message: "Sesi Anda telah habis. Silakan login kembali"}

// media-svc
// error code upload media
var ErrUploadMedia = Message{Code: 109400, Message: "Upload Failed"}
