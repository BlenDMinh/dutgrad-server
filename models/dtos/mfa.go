package dtos

type MFASetupRequest struct {
}

type MFASetupResponse struct {
	Secret          string   `json:"secret"`
	QRCodeDataURL   string   `json:"qr_code_data_url"`
	BackupCodes     []string `json:"backup_codes"`
	ProvisioningURI string   `json:"provisioning_uri"`
}

type MFAVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

type MFADisableRequest struct {
}

type MFALoginVerifyRequest struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	MFACode       string `json:"mfa_code"`
	UseBackupCode bool   `json:"use_backup_code"`
}

type MFALoginCompleteRequest struct {
	Code          string `json:"code" binding:"required"`
	UseBackupCode bool   `json:"use_backup_code"`
}

type MFATempTokenResponse struct {
	TempToken string `json:"temp_token"`
	ExpiresAt string `json:"expires_at"`
}

type MFAExternalVerifyRequest struct {
	Code          string `json:"code" binding:"required"`
	UseBackupCode bool   `json:"use_backup_code"`
}
