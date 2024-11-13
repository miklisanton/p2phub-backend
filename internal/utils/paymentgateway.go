package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func VerifySignature (confirmReq map[string]interface{}) error {
    Logger.LogInfo().Fields(map[string]interface{}{
        "request": confirmReq,
    }).Msg("Confirm request")
	// Verify signature
    sign, ok := confirmReq["sign"].(string)
    if !ok {
        return fmt.Errorf("Sign not found")
    }
    delete(confirmReq, "sign")
    jsonData, err := json.Marshal(confirmReq)
    if err != nil {
        return err
    }
	// Manually escape forward slashes
	escapedData := strings.ReplaceAll(string(jsonData), "/", "\\/")
	// Create signature and add it to request headers
	base64Req := base64.StdEncoding.EncodeToString([]byte(escapedData))
    Logger.LogInfo().Str("key", os.Getenv("GATEWAY_API_KEY")).Msg("API key")
	hash := md5.Sum([]byte(base64Req + os.Getenv("GATEWAY_API_KEY")))
	// Compare hash
	if fmt.Sprintf("%x", hash) != sign {
        Logger.LogInfo().Str("hash", fmt.Sprintf("%x", hash)).Msg("Hash")
        Logger.LogInfo().Str("sign", sign).Msg("Sign")
		return fmt.Errorf("Invalid signature")
	}
    return nil
}
