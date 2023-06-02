package routines

import (
	"watch-ms/DI"
	"watch-ms/service/episode"
	"watch-ms/service/watch"
)

func Init(services DI.ServiceContainer) {
	startEpisodeRoutines(services.EpisodeService)
	startWatchRoutines(services.WatchService)
}

func startEpisodeRoutines(service episode.Service) {
	go service.OnUploadNewEpisode()
}

func startWatchRoutines(service watch.Service) {
	go service.OnUpdateWatchThumbnail()
}
