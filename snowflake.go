package gophbot

import (
	"strconv"
	"time"
)

// Snowflake is a convenience typealias depicting the format used to store snowflakes
type Snowflake = string

// SnowflakeTime decodes a snowflake and fetches the timestamp element.
func SnowflakeTime(s Snowflake) time.Time {
	const discordEpoch int64 = 1420070400000

	// Ignoring this error, as we'll just return a timestamp of 0
	i, _ := strconv.ParseInt(s, 10, 64)

	return time.Unix((i>>22+discordEpoch)/1000, 0)
}
