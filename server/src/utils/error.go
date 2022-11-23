package utils

import "github.com/rs/zerolog/log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Error().Msgf("%s: %s", msg, err)
	}
}

func LogWithInfo(cmp, msg string) {

	log.Info().Msgf("%s: %s", cmp, msg)

}
