package auth
import (
  "gmcc/internal/crypto"
  "gmcc/internal/logger"
  "time"
  "fmt"
)

func TimeStampToBeijingTime(ts int64) string {
 	// 加载上海时区（北京时间）
 	loc, err := time.LoadLocation("Asia/Shanghai")
 	if err != nil {
 		// 加载失败时使用 UTC 作为 fallback
 		loc = time.UTC
 	}
 	// 时间戳转 time.Time
 	t := time.Unix(ts, 0).In(loc)
 	// 格式化
 	return t.Format("2006-01-02 15:04:05")
 }

var mcToken string

func LoginAccount(playeriid string) (string, error) {
	playerID := playeriid
	// 1. 尝试直接加载 MC Token
	var mcCache StoredToken
	// 修复：使用 := 声明并初始化 err
	err := crypto.LoadToken(playerID+"_mc.token", &mcCache)

	if err == nil && time.Now().Unix() < mcCache.ExpiresAt {
		logger.Info("使用缓存的 MC token")
		logger.Info("MC token过期时间:",TimeStampToBeijingTime(int64(mcCache.ExpiresAt)))
		mcToken = mcCache.AccessToken
	} else {
		logger.Info("MC token 不存在或已过期，执行 MS 登录...")

		// 2. 检查并获取 MS Token
		var msStored StoredToken
		msErr := crypto.LoadToken(playerID+"_ms.token", &msStored)
		
		var msAccessToken string
		if msErr == nil {
			// 尝试刷新 MS Token
			logger.Info("正在刷新 MS token...")
			if refreshErr := RefreshToken(playerID); refreshErr == nil {
				// 刷新成功，重新加载获取最新的 AccessToken
				_ = crypto.LoadToken(playerID+"_ms.token", &msStored)
				msAccessToken = msStored.AccessToken
			} else {
				logger.Info("刷新失败，需重新进行设备流登录")
				// 修复：使用 = 赋值，因为 err 在上方已由简短声明定义（或者这里重新定义）
				msAccessToken, err = MSToken(playerID)
			}
		} else {
			msAccessToken, err = MSToken(playerID)
		}

		if err != nil {
			return "", err
		}

		// 3. 走完整流程获取 MC Token
		xblToken, uhs, err := GetXBLAccess(msAccessToken)
		if err != nil { return "", err }

		xstsToken, uhs, err := GetXSTSToken(xblToken)
		if err != nil { return "", err }

		mcToken, err = GetMinecraftToken(playerID, xstsToken, uhs)
		if err != nil { return "", err }
		
		// 4. 最终检查所有权
    	owns, err := HasMinecraftOwnership(mcToken)
    	if err != nil { return "", err }
    	if !owns { return "",fmt.Errorf("该账户未拥有 Minecraft") }
	}

	logger.Info("登录成功")
	return mcToken, nil
}