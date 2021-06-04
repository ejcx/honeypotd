package notification

type Notification interface{
  Notify() error
}
