package trinity

import (
	"github.com/bwmarrin/snowflake"
)

// GenerateSnowFlakeID generate snowflake id
func GenerateSnowFlakeID(nodenumber int64) int64 {

	// Create a new Node with a Node number of 1
	node, _ := snowflake.NewNode(nodenumber)

	// Generate a snowflake ID.
	id := node.Generate().Int64()
	return id

}
