package video_collection

type Collection interface {
	Run(cfg *Configuration, retry Retry) error
	Reload(cfg *Configuration, retry Retry) error
	Exit(cfg *Configuration)
}
