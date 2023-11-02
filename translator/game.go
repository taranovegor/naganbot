package translator

var GameTranslations = translations{
	"ru": {
		"game creation": {
			oneOf: []oneOf{
				{message: "Мы играем чтобы пощекотать нервы, хочешь тоже?"},
				{message: "Я вижу ты не из трусливых, что насчёт остальных?"},
				{message: "Вы опять за старое? Так уж и быть"},
			},
		},
		"joining the game": {
			oneOf: []oneOf{
				{message: "Могу поспорить, что ты - труп"},
				{message: "Ты считаешь себя бессмертным?"},
				{message: "Ну хорошо, присаживайся, игра скоро начнётся"},
				{message: "Решил напоследок рискнуть?"},
				{message: "1 к 6 - вот твой шанс выживания в русской рулетке"},
				{message: "Здесь произойдёт охуенное убийство"},
				{message: "Ты играл когда-нибудь до этого в рулетку?"},
			},
		},
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: "Решили рискнуть? Тогда начнём"},
						{message: "Стопка, щелчок, выстрел..."},
					},
				},
				{
					allOf: []oneOf{
						{message: "Хорошо, раз все в сборе, начнём"},
						{message: "Пистолет передаётся следующему..."},
					},
				},
			},
		},
		"gunslinger killed":       {message: "%gunslinger прострелил себе голову"},
		"killed by atomic bullet": {message: "Комнату озарила яркая вспышка. В нагане была атомная пуля."},
		"joined the game":         {message: "Участники игры %date"},
		"game join list item":     {message: "%num. %gunslinger"},
		"owner of the game":       {message: "владелец"},
		"shot in game":            {message: "выбыл(а)"},
		"top is not determined":   {message: "Пока что мы не можем составить топ игроков"},
		"top players by games":    {message: "Топ-%number выбывших игроков:"},
		"top game player":         {message: "%i. %user - %times раз(а)"},
		"available only in chat":  {message: "Русская рулетка игра не для одного"},
		"player already in game":  {message: "Вы уже в игре"},
		"player is not kicked":    {message: "У него еще есть пульс. Может, откачают?"},
		"wait for game timeout":   {message: "Ещё не время начинать игру, мы должны залечь на дно"},
		"active game not found":   {message: "Ещё никто не начал игру"},
		"user game statistics":    {message: "Вы участвовали в %games игре(-ах) и проиграли %shots раз(а)"},
		"deprecated command":      {message: "Текущая команда в дальнейшем будет удалена, используйте вместо неё %new"},
	},
}
