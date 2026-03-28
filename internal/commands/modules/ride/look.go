package ride

import (
	"math"
)

func calculateLookAt(bx, by, bz float64, tx, ty, tz float64) (yaw, pitch float32) {
	dx := tx - bx
	dy := ty - by
	dz := tz - bz

	yaw = float32(math.Atan2(-dx, dz) * 180 / math.Pi)
	horizDist := math.Sqrt(dx*dx + dz*dz)
	pitch = float32(math.Atan2(-dy, horizDist) * 180 / math.Pi)

	return yaw, pitch
}

func smoothLook(currentYaw, currentPitch, targetYaw, targetPitch, smoothing float32) (newYaw, newPitch float32) {
	if smoothing <= 0 {
		smoothing = 0.1
	}
	if smoothing > 1 {
		smoothing = 1
	}

	// 处理 yaw 跨越 ±180° 的情况
	// 选择最短路径：顺时针或逆时针
	yawDiff := normalizeAngle(targetYaw - currentYaw)
	pitchDiff := targetPitch - currentPitch

	newYaw = currentYaw + yawDiff*smoothing
	newPitch = currentPitch + pitchDiff*smoothing

	// 规范化结果到 [-180, 180]
	newYaw = normalizeAngle(newYaw)

	return newYaw, newPitch
}

func normalizeAngle(angle float32) float32 {
	for angle < -180 {
		angle += 360
	}
	for angle > 180 {
		angle -= 360
	}
	return angle
}
