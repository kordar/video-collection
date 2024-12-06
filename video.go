package video_collection

type Video interface {
	Run(config Configuration) error
	Stop()
}
