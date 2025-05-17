package translator

var GameTranslations = translations{
	"ru": {
		"game creation": {
			oneOf: []oneOf{
				{message: SimpleMessage("Мы играем чтобы пощекотать нервы, хочешь тоже?")},
				{message: SimpleMessage("Я вижу ты не из трусливых, что насчёт остальных?")},
				{message: SimpleMessage("Вы опять за старое? Так уж и быть")},
				{message: SimpleMessage("Ну что, рискнёте сыграть в игру с летальным исходом?")},
				{message: SimpleMessage("Смерть или слава? Давайте узнаем!")},
				{message: SimpleMessage("Кто-то сегодня не вернётся домой...")},
				{message: SimpleMessage("Пистолет заряжен, очередь за вами.")},
				{message: SimpleMessage("Русская рулетка — лучший способ проверить удачу.")},
				{message: SimpleMessage("Готовы испытать судьбу?")},
				{message: SimpleMessage("Один выстрел — и ты в истории.")},
				{message: SimpleMessage("Судьба любит смелых. Или нет?")},
				{message: SimpleMessage("На кону жизнь. Ну, почти.")},
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
				{message: SimpleMessage("Добро пожаловать в клуб самоубийц.")},
				{message: SimpleMessage("Шанс выжить? Примерно как у снега в аду.")},
				{message: SimpleMessage("Ты либо герой, либо покойник.")},
				{message: SimpleMessage("Один патрон, шесть камер. Удачи.")},
				{message: SimpleMessage("Ты уверен, что хочешь это сделать?")},
				{message: SimpleMessage("Сейчас узнаем, насколько ты везучий.")},
				{message: SimpleMessage("Ты либо выйдешь победителем, либо станешь уроком для остальных.")},
				{message: SimpleMessage("Ну что, готов к своему последнему клику?")},
				{message: SimpleMessage("Судьба улыбается... или нет?")},
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
				{
					allOf: []oneOf{
						{message: SimpleMessage("Игра началась. Пусть удача будет на вашей стороне.")},
						{message: SimpleMessage("Колесо фортуны крутится...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Ну что, паника или спокойствие?")},
						{message: SimpleMessage("Пистолет уже в игре...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Время узнать, кто сегодня проиграл.")},
						{message: SimpleMessage("Палец на спусковом курке...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Готовы к самому страшному клику в вашей жизни?")},
						{message: SimpleMessage("Курок нажат...")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Сейчас всё решит один выстрел.")},
						{message: SimpleMessage("Тишина перед бурей...")},
					},
				}},
		},
		"gunslinger killed": {message: SimpleMessage("%gunslinger прострелил себе голову")},
		"killed by atomic bullet": {
			oneOf: []oneOf{
				{message: SimpleMessage("Щёлк — и внезапно комната превратилась в эпицентр ядерного гриба. В револьере оказалась атомная пуля")},
				{message: SimpleMessage("...от группы смельчаков остался только радиоактивный след. Видимо, кто-то подсунул в барабан не ту пулю")},
				{message: SimpleMessage("Ккомната наполнилась ярким светом... Все игроки мгновенно испарились. Кто-то явно жульничал с боеприпасами")},
			},
		},
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
