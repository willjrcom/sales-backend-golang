package processentity

type StatusProcess string

const (
	ProcessStatusPending   StatusProcess = "Pending"
	ProcessStatusStarted   StatusProcess = "Started"
	ProcessStatusFinished  StatusProcess = "Finished"
	ProcessStatusPaused    StatusProcess = "Paused"
	ProcessStatusContinued StatusProcess = "Continued"
)

func GetAllDeliveryStatus() []StatusProcess {
	return []StatusProcess{
		ProcessStatusPending,
		ProcessStatusStarted,
		ProcessStatusFinished,
		ProcessStatusPaused,
		ProcessStatusContinued,
	}
}
