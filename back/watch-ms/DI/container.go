package DI

import (
	"gorm.io/gorm"
	categoryController "watch-ms/controller/category"
	seasonController "watch-ms/controller/season"
	titleController "watch-ms/controller/title"
	watchController "watch-ms/controller/watch"
	categoryProvider "watch-ms/provider/category"
	episodeProvider "watch-ms/provider/episode"
	seasonProvider "watch-ms/provider/season"
	titleProvider "watch-ms/provider/title"
	watchProvider "watch-ms/provider/watch"
	categoryService "watch-ms/service/category"
	episodeService "watch-ms/service/episode"
	seasonService "watch-ms/service/season"
	titleService "watch-ms/service/title"
	watchService "watch-ms/service/watch"
)

type ServiceContainer struct {
	WatchService    watchService.Service
	SeasonService   seasonService.Service
	TitleService    titleService.Service
	CategoryService categoryService.Service
	EpisodeService  episodeService.Service
}

type ProviderContainer struct {
	WatchProvider    watchProvider.Provider
	SeasonProvider   seasonProvider.Provider
	TitleProvider    titleProvider.Provider
	CategoryProvider categoryProvider.Provider
	EpisodeProvider  episodeProvider.Provider
}

type ControllerContainer struct {
	WatchController    watchController.Controller
	SeasonController   seasonController.Controller
	TitleController    titleController.Controller
	CategoryController categoryController.Controller
}

type Container struct {
	Services    ServiceContainer
	Providers   ProviderContainer
	Controllers ControllerContainer
}

func NewContainer(dbInstance *gorm.DB) Container {
	providers := initProviders(dbInstance)
	services := initServices(providers)
	controllers := initControllers(services)
	return Container{
		Services:    services,
		Providers:   providers,
		Controllers: controllers,
	}
}

func initProviders(dbInstance *gorm.DB) ProviderContainer {
	episode := episodeProvider.NewEpisodeProvider(dbInstance)
	watch := watchProvider.NewWatchProvider(dbInstance, episode)
	season := seasonProvider.NewSeasonProvider(dbInstance)
	title := titleProvider.NewTitleProvider(dbInstance)
	category := categoryProvider.NewCategoryProvider(dbInstance)
	return ProviderContainer{
		WatchProvider:    watch,
		SeasonProvider:   season,
		TitleProvider:    title,
		CategoryProvider: category,
		EpisodeProvider:  episode,
	}
}

func initServices(providers ProviderContainer) ServiceContainer {
	episode := episodeService.NewEpisodeService(providers.EpisodeProvider)
	watch := watchService.NewWatchService(providers.WatchProvider)
	season := seasonService.NewSeasonService(providers.SeasonProvider)
	title := titleService.NewTitleService(providers.TitleProvider)
	category := categoryService.NewCategoryService(providers.CategoryProvider)
	return ServiceContainer{
		WatchService:    watch,
		SeasonService:   season,
		TitleService:    title,
		CategoryService: category,
		EpisodeService:  episode,
	}
}

func initControllers(services ServiceContainer) ControllerContainer {
	watch := watchController.NewWatchController(services.WatchService)
	season := seasonController.NewSeasonController(services.SeasonService)
	title := titleController.NewTitleController(services.TitleService)
	category := categoryController.NewCategoryController(services.CategoryService)
	return ControllerContainer{
		WatchController:    watch,
		SeasonController:   season,
		TitleController:    title,
		CategoryController: category,
	}
}
