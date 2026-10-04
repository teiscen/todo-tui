package notes

func TestNotes() Notes {
	return NewNotes("")
}

// func NotesTestPrint() {
// 	n := TestNotes()
// 	println(n.View())
// }
//
// func NotesTestToggle() {
// 	// fmt.Print("\034[H\033[2J") // Clear console ANSI escape sequence
// 	n := TestNotes()
//
// 	for range 10 {
// 		// In Focus
// 		fmt.Print("\034[H\033[2J")
// 		n.TextArea.Focus()
// 		println(n.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Not in Focus
// 		fmt.Print("\034[H\033[2J")
// 		n.TextArea.Blur()
// 		println(n.View())
// 		time.Sleep(1 * time.Second)
// 	}
// }
