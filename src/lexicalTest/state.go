package lexicalTest

import "fmt"

const eps = ' '

// TODO: 重构！！！！！
type State struct {
	endFlag     bool
	markFlag    byte
	toNextState map[byte][]*State
}

func NewState(endFlag bool) *State {
	return &State{endFlag, eps, make(map[byte][]*State)}
}

func (s *State) getFather(fatherSet map[*State]*State) *State {
	return fatherSet[s]
}
func (s *State) Merge(hasVisited map[*State]bool) {
	if hasVisited[s] {
		return
	}
	hasVisited[s] = true
	// 消除空白态
	for s.haveBlankStates() {
		mapOfReachableStateOfBlankStates := s.formMapOfReachableStateOfBlankStates()
		if s.isNextBlankStatesHaveEndState() {
			s.setEndFlag(true)
		}
		s.cleanBlankStates()
		s.AddNextStates(mapOfReachableStateOfBlankStates)
	}

	//对非空白态的子节点进行处理
	allNextStates := s.getAllNextStates()
	for _, nextState := range allNextStates {
		nextState.Merge(hasVisited)
	}
	return
}
func (s *State) setEndFlag(value bool) {
	s.endFlag = value
}
func (s *State) formMapOfReachableStateOfBlankStates() map[byte][]*State {
	blankStates := s.getNextBlankStates()
	return getStatesToNext(blankStates)
}
func (s *State) isNextBlankStatesHaveEndState() bool {
	nextBlankStates := s.getNextBlankStates()
	return haveEndState(nextBlankStates)
}
func (s *State) haveBlankStates() bool {
	return len(s.getNextBlankStates()) != 0
}

func (s *State) getNextBlankStates() []*State {
	return s.getNextStates(eps)
}
func (s *State) getAllNextStates() []*State {
	result := make([]*State, 0)
	for char := range s.toNextState {
		result = append(result, s.getNextStates(char)...)
	}
	return result
}

func (s *State) AddNextStates(addedMap map[byte][]*State) {
	for char, states := range addedMap {
		s.toNextState[char] = append(s.toNextState[char], states...)
	}
}

func (s *State) cleanBlankStates() {
	s.cleanNextStates(eps)
}
func (s *State) cleanNextStates(char byte) {
	delete(s.toNextState, char)
}

// TODO： N -> X | Z 有问题..  这个的DFA不正确...
// TODO: 有时候测试可以通过，有时候不可以...
//var i=0
func (s *State) DFA(hasVisited map[*State]bool) *State {
	//i++
	//if i==10{
	//	panic("toolllll")
	//}
	if hasVisited[s] {
		return s
	}
	hasVisited[s] = true
	for char := range s.toNextState {
		dfaState := s.getDFAState(char)
		s.cleanNextStates(char)
		s.LinkByChar(char,dfaState.DFA(hasVisited))
	}
	return s
}



// TODO: 这个函数还是有问题的     D+.D+|D+ 这种情况不能判断
func (s *State) getDFAState(char byte) *State{
	states := s.toNextState[char]
	if len(states) == 1 {
		return states[0]
	}
	dfaState := NewState(s.hasEndFlag(s.toNextState[char]))
	dfaState.toNextState = s.formMapOfReachableStateOfAllNextStates()

	if s.toNextIsSame(dfaState) {
		dfaState = s.toNextState[char][0]
		return dfaState
	}
	if s.hasSelf(char) {
		dfaState.LinkByChar(char, dfaState)
	}
	return dfaState
}

func (s *State) formMapOfReachableStateOfAllNextStates() map[byte][]*State {
	allNextStates := s.getAllNextStates()
	return getStatesToNext(allNextStates)
}



func (s *State) toNextIsSame(reference *State) bool {
	if len(s.toNextState)!=len(reference.toNextState){
		return false
	}

	for char, nextStates := range reference.toNextState {
		HasVisitOfS := make(map[*State]bool)
		HasVisitOfRef := make(map[*State]bool)
		for _, nextState1 := range nextStates {
			HasVisitOfRef[nextState1] = true
		}
		for _, nextState1 := range s.toNextState[char] {
			HasVisitOfS[nextState1] = true
		}
		//if len(HasVisitOfS)!=len(HasVisitOfRef){
		//	return false
		//}
		for state := range HasVisitOfS{
			if HasVisitOfRef[state]==false{
				return false
			}
		}
		for state := range HasVisitOfRef{
			if HasVisitOfS[state]==false{
				return false
			}
		}

	}
	return true
}

func (s *State) stateIsLiving(char byte, x *State) bool {
	for _, state := range s.toNextState[char] {
		if state == x {
			return true
		}
	}
	return false
}

func (s *State) CanBeStartOfDFA(hasVisited map[*State]bool) bool {
	if hasVisited[s] {
		return true
	}
	hasVisited[s] = true

	charsOfLinkingToNextStates := s.getTheCharsOfLinkingToNextStates()
	for _, char := range charsOfLinkingToNextStates {
		nextStates := s.getNextStates(char)
		if len(nextStates) != 1 {
			return false
		}
		// 往后搜索
		for _, state := range nextStates {
			if state.CanBeStartOfDFA(hasVisited) == false {
				return false
			}
		}
	}
	return true
}
func (s *State) IsMatch(pattern string) bool {
	// 空匹配
	nextStates := s.toNextState[eps]
	for _, nextState := range nextStates {
		if nextState.IsMatch(pattern) {
			return true
		}
	}
	if pattern == "" {
		return s.endFlag
	}

	ch := pattern[0]

	// 实匹配
	nextStates = s.toNextState[ch]
	for _, nextState := range nextStates {
		if nextState.IsMatch(pattern[1:]) {
			return true
		}
	}
	return false
}

func (s *State) MarkDown(specialChar byte, stateIsVisit map[*State]bool) {
	currentState := s
	if stateIsVisit[currentState] {
		return
	}
	stateIsVisit[currentState] = true
	s.markFlag = specialChar
	for _, nextStates := range s.toNextState {
		for _, nextState := range nextStates {
			nextState.MarkDown(specialChar, stateIsVisit)
		}
	}
}

func (s *State) Show(startId int, stateToId map[*State]int, stateIsVisit map[*State]bool) {
	currentState := s
	if stateIsVisit[currentState] {
		return
	}
	stateIsVisit[currentState] = true
	stateToId[currentState] = startId
	for bytes, nextStates := range s.toNextState {
		for _, nextState := range nextStates {
			nextState.Show(len(stateToId), stateToId, stateIsVisit)
			option := string(bytes)
			fmt.Printf("id:%d%s|%s --%s--> id:%d%s|%s\n",
				stateToId[currentState],
				currentState.getEndMark(),
				string(currentState.markFlag),
				option,
				stateToId[nextState],
				nextState.getEndMark(),
				string(currentState.markFlag),
			)
		}
	}
}

func (s *State) LinkByChar(ch byte, nextState *State) {
	s.toNextState[ch] = append(s.toNextState[ch], nextState)
}
func (s *State) Link(nextState *State) {
	s.LinkByChar(eps, nextState)
}

func (s *State) getNextStates(char byte) []*State {
	return s.toNextState[char]
}
func (s *State) getTheCharsOfLinkingToNextStates() []byte {
	chars := make([]byte, 0)
	for char := range s.toNextState {
		chars = append(chars, char)
	}
	return chars
}

func getStatesToNext(states []*State) map[byte][]*State {
	result := make(map[byte][]*State)
	for _, state := range states {
		for char, nextStates := range state.toNextState {
			for _, nextState := range nextStates {
				result[char] = append(result[char], nextState)
			}
		}
	}
	return result
}
func (s *State) hasSelf(char byte) bool {
	for _, state := range s.toNextState[char] {
		if state == s {
			return true
		}
	}
	return false
}
func (s *State) hasEndFlag(states []*State) bool {
	for _, state := range states {
		if state.endFlag {
			return true
		}
	}
	return false
}

func (s *State) getEndMark() string {
	if s.endFlag == true {
		return "(OK)"
	}
	return "    "
}
