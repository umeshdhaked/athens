package dtos

type UploadKycRequest struct {
	DocType string
	//KycDocument multipart.File
	//Photo       multipart.File
}

type UploadKycResponse struct {
	Status string
}

type PendingKycsRequest struct {
	Pagination
}

type PendingKycsResponse struct {
	UserName    string
	Mobile      string
	DocType     string
	KycDocument []byte
	Photo       []byte
	KycId       int64
}

type KycStatusUpdateRequest struct {
	KycId   int64
	Status  string
	Comment string
}

type KycStatusUpdateResponse struct {
	KycId  int64
	Status string
}

type KycStatusRequest struct {
	UserId int64
}

type KycStatusResponse struct {
	KycId   int64
	Status  string
	Comment string
}
