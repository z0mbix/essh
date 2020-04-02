package config

// func TestTagDoubleStringWithAdditionalFlags(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p server1 ss -- -A -4 uptime
// 	args := []string{
// 		"--debug",
// 		"-r",
// 		"ap-southeast-1",
// 		"-p",
// 		"server1",
// 		"ss",
// 		"--",
// 		"-A",
// 		"-4",
// 		"uptime",
// 	}

// 	os.Args = args
// 	_, err := getESSHConfig()

// 	if err.Error() != "only specifiy an instance id or a tag, if a tag has a space, wrap in double quotes" {
// 		t.Error("failed to detect DoubleString error")
// 	}
// }

// func TestTagDoubleString(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p server1 ss
// 	args := []string{
// 		"--debug",
// 		"-r",
// 		"ap-southeast-1",
// 		"-p",
// 		"server1",
// 		"ss",
// 	}

// 	os.Args = args

// 	_, err := getESSHConfig()
// 	fmt.Println(err)
// 	if err.Error() != "only specifiy an instance id or a tag, if a tag has a space, wrap in double quotes" {
// 		t.Error("failed to detect DoubleString error")
// 	}
// }

// func TestTagDoubleStringWithAdditionalArgsButNoDoubleDash(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p server1 ddd -A -4 uptime
// }

// func TestTagSingleTagWithAdditionalFlagsButNoDoubleDash(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p server1 -A -4 uptime
// }

// func TestInstIDDoubleStringWithAdditionalFlags(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p i-xxxx ss -- -A -4 uptime
// }

// func TestInstIDDoubleString(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p i-xxxx ss
// }

// func TestInstIDWithDoubleStringWithAdditionalFlagsNoDoubleDash(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p i-xxxx dd -A -4 uptime
// }

// func TestNoInstIDOrTagWithAdditionalFlagsNoDoubleDash(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p -A -4 uptime
// }

// func TestTagCorrectTagAndAdditionalFlags(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p server1 -- -A -4 uptime
// }

// func TestCTagCorrectTag(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p server1
// }

// func TestTagCorrectQuotedTag(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p "server 1"
// }

// func TestTagCorrectQuotedTagWithAdditionalFlags(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p "server 1" -- -A -4 uptime
// }

// func TestInstIDCorrectWithAdditionalFlags(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p i-xxxxx -- -A -4 uptime
// }

// func TestInstIDCorrectJustID(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p i-xxxxx
// }

// func TestCorrectNoTagNoInstID(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p
// }

// func TestCorrectNoTagNoInstIDWithAdditionalFlags(t *testing.T) {
// 	// --debug -r ap-southeast-1 -p -- -A -4 uptime
// }
