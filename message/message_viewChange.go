package message

var (
	ViewChangePropose MessageType = "ViewChangePropose"
	NewChange         MessageType = "NewChange"
)

type ViewChangeMsg struct {
	CurView  int
	NextView int
	SeqID    int
	FromNode uint64
}

type NewViewMsg struct {
	CurView  int
	NextView int
	NewSeqID int
	FromNode uint64
}
