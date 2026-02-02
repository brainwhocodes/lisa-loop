package tui

type OutputTab int

const (
	OutputTabTranscript OutputTab = iota
	OutputTabDiffs
	OutputTabReasoning
)

func (t OutputTab) String() string {
	switch t {
	case OutputTabTranscript:
		return "Transcript"
	case OutputTabDiffs:
		return "Diffs"
	case OutputTabReasoning:
		return "Reasoning"
	default:
		return "Transcript"
	}
}
