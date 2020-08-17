package client

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {
	var (
		a = `L 08/10/2020 - 12:11:04: "Dio Brando<142><[U:1:927963222]><Red>" say_team "so cloooose"`
		b = `L 08/10/2020 - 12:11:49: "råBïÐh¥þêrf0¢u§<141><[U:1:141229441]><Blue>" say "i need it"`
	)
	am := reSay.FindStringSubmatch(a)
	require.Equal(t, len(am), 5)
	bm := reSay.FindStringSubmatch(b)
	require.Equal(t, len(bm), 5)
}
