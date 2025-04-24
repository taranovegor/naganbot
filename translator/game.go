package translator

var GameTranslations = translations{
	"ru": {
		"game creation": {
			oneOf: []oneOf{
				{message: SimpleMessage("Мы играем чтобы пощекотать нервы, хочешь тоже?")},
				{message: SimpleMessage("Я вижу ты не из трусливых, что насчёт остальных?")},
				{message: SimpleMessage("Вы опять за старое? Так уж и быть")},
			},
		},
		"joining the game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Могу поспорить, что ты - труп")},
				{message: SimpleMessage("Ты считаешь себя бессмертным?")},
				{message: SimpleMessage("Ну хорошо, присаживайся, игра скоро начнётся")},
				{message: SimpleMessage("Решил напоследок рискнуть?")},
				{message: SimpleMessage("1 к 6 - вот твой шанс выживания в русской рулетке")},
				{message: SimpleMessage("Здесь произойдёт охуенное убийство")},
				{message: SimpleMessage("Ты играл когда-нибудь до этого в рулетку?")},
			},
		},
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: SimpleMessage("Решили рискнуть? Тогда начнём")},
						{message: SimpleMessage("Стопка, щелчок, выстрел...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Хорошо, раз все в сборе, начнём")},
						{message: SimpleMessage("Пистолет передаётся следующему...")},
					},
				},
			},
		},
		"gunslinger killed":     {message: SimpleMessage("%gunslinger прострелил себе голову")},
		"joined the game":       {message: SimpleMessage("Участники игры %date")},
		"game join list item":   {message: SimpleMessage("%num. %gunslinger")},
		"owner of the game":     {message: SimpleMessage("владелец")},
		"shot in game":          {message: SimpleMessage("выбыл")},
		"top is not determined": {message: SimpleMessage("Пока что мы не можем составить топ игроков")},
		"top players by games":  {message: SimpleMessage("Топ-%number выбывших игроков:")},
		"top game player": {message: PluralMessage{
			one:  "%i. %user - %times раз",
			few:  "%i. %user - %times раза",
			many: "%i. %user - %times раз",
		}},
		"available only in chat": {message: SimpleMessage("Русская рулетка игра не для одного")},
		"player already in game": {message: SimpleMessage("Вы уже в игре")},
		"player is not kicked":   {message: SimpleMessage("У него еще есть пульс. Может, откачают?")},
		"wait for game timeout":  {message: SimpleMessage("Ещё не время начинать игру, мы должны залечь на дно")},
		"active game not found":  {message: SimpleMessage("Ещё никто не начал игру")},
		"user game statistics":   {message: SimpleMessage("Вы участвовали в %games игре(-ах) и проиграли %shots раз(а)")},
	},
}
