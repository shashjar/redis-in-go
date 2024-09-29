package protocol

// TODO: finish implementing Redis protocol
/**
 * Implementation of the Redis Serialization Protocol (RESP): https://redis.io/docs/latest/develop/reference/protocol-spec/
 */

const (
	SIMPLE_STRING   = "+"
	SIMPLE_ERROR    = "-"
	INTEGER         = ":"
	BULK_STRING     = "$"
	ARRAY           = "*"
	NULL            = "_"
	BOOLEAN         = "#"
	DOUBLE          = ","
	BIG_NUMBER      = "("
	BULK_ERROR      = "!"
	VERBATIM_STRING = "="
	MAP             = "%"
	SET             = "~"
	PUSH            = ">"
)
