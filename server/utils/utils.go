package utils

import (
	"os"
	"strconv"
)

// set instance_id in env if you have mutil instance
// instance id should start from 1 and increase by step 1
// if not, instance id return 0
func GetInstanceID() int {
	instanceIDStr := os.Getenv("instance_id")
	instanceID, _ := strconv.Atoi(instanceIDStr)
	return instanceID
}

func GetAllInstanceCount() int {
	instanceCountStr := os.Getenv("instance_count")
	instanceCount, _ := strconv.Atoi(instanceCountStr)
	return instanceCount
}
