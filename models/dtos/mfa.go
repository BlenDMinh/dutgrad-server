package dtos

// MFA setup request and response
type MFASetupRequest struct {
	// No password required
}

type MFASetupResponse struct {
	Secret          string   `json:"secret"`
	QRCodeDataURL   string   `json:"qr_code_data_url"`
	BackupCodes     []string `json:"backup_codes"`
	ProvisioningURI string   `json:"provisioning_uri"`
}

// MFA verify request
type MFAVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

// MFA disable request
type MFADisableRequest struct {
	// No password required
}

// MFA login verification request and response
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

// Session token for partial authentication
type MFATempTokenResponse struct {
	TempToken string `json:"temp_token"`
	ExpiresAt string `json:"expires_at"`
}

// MFAExternalVerifyRequest is the request body for verifying MFA for external auth
type MFAExternalVerifyRequest struct {
	Code          string `json:"code" binding:"required"`
	UseBackupCode bool   `json:"use_backup_code"`
}
