package speech

func (s *SpeechBase) SetPosition(pos string) {
	s.position = pos
}

func (s *SpeechBase) GetPosition() string {
	return s.position
}

func (s *SpeechBase) SetSpeed(speed int) {
	s.speed = speed
}

func (s *SpeechBase) GetSpeed() int {
	return s.speed
}
