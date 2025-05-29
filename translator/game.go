package translator

var GameTranslations = translations{
	"ru": {
		"something went wrong": {message: SimpleMessage("Наган заклинило. Мы уже вызвали оружейника.")},
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
		"gunslinger killed": {
			oneOf: []oneOf{
				{message: SimpleMessage("Раздался выстрел — и %gunslinger больше с нами")},
				{message: SimpleMessage("Барабан провернулся, курок щёлкнул... %gunslinger выбывает")},
				{message: SimpleMessage("Пуля нашла свою жертву — это был %gunslinger")},
				{message: SimpleMessage("Кровь, порох и тишина... %gunslinger проиграл")},
				{message: SimpleMessage("%gunslinger решил, что жил слишком долго")},
				{message: SimpleMessage("В этот раз не повезло %gunslinger")},
				{message: SimpleMessage("На этот раз смерть выбрала %gunslinger")},
				{message: SimpleMessage("Роковой выстрел для %gunslinger")},
				{message: SimpleMessage("%gunslinger больше не ответит")},
				{message: SimpleMessage("Прощай, %gunslinger, мы тебя запомним...")},
				{message: SimpleMessage("Игра окончена для %gunslinger")},
				{message: SimpleMessage("Момент истины для %gunslinger оказался последним")},
				{message: SimpleMessage("%gunslinger получил билет в один конец.")},
			},
		},
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
		"player already in game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Ты уже в игре. Смерть тебя уже запомнила")},
				{message: SimpleMessage("Жаждешь смерти? Подожди своей очереди")},
				{message: SimpleMessage("Ты уже в списке на расстрел. Терпение.")},
				{message: SimpleMessage("Ты уже в деле. Расслабься и жди выстрела.")},
				{message: SimpleMessage("Ты уже участвуешь. Или передумал?")},
				{message: SimpleMessage("Смерть любит нетерпеливых")},
				{message: SimpleMessage("Ты уже в игре, не торопи смерть")},
				{message: SimpleMessage("Один раз в игре - достаточно для проверки удачи")},
				{message: SimpleMessage("Ты уже в списке потенциальных покойников")},
			},
		},
		"player is not kicked": {
			oneOf: []oneOf{
				{message: SimpleMessage("Револьвер заклинило. Видимо, у кого-то иммунитет")},
				{message: SimpleMessage("Порох отсырел — приговор не был приведён в исполнение")},
				{message: SimpleMessage("Смерть сегодня капризничает. Кое-кто временно неуязвим")},
				{message: SimpleMessage("Казнь откладывается — технические неполадки.")},
				{message: SimpleMessage("Кто-то подложил холостой патрон")},
			},
		},
		"player already played today": {
			oneOf: []oneOf{
				{message: SimpleMessage("Сегодня ты уже играл. Смерть не торопится - завтра будет новый шанс.")},
				{message: SimpleMessage("Ты сегодня уже ходил по краю.")},
				{message: SimpleMessage("Смерть сегодня уже видела тебя.")},
				{message: SimpleMessage("Ты уже сегодня играл в рулетку. Не хочешь дожить до завтра?")},
				{message: SimpleMessage("Один шанс в день - таковы правила этой игры.")},
				{message: SimpleMessage("Сегодня ты уже крутил барабан.")},
				{message: SimpleMessage("Смерть сегодня уже брала свою дань с тебя.")},
			},
		},
		"active game not found":                  {message: SimpleMessage("Ещё никто не начал игру")},
		"user game statistics":                   {message: SimpleMessage("Вы участвовали в %games игре(-ах) и проиграли %shots раз(а)")},
		"available settings below":               {message: SimpleMessage("Персонализируйте игру используя перечисленные ниже опции")},
		"settings can be changed only by admins": {message: SimpleMessage("Изменять настройки игры могут только администраторы чата")},
		"4 shot revolver":                        {message: SimpleMessage("Colt Cloverleaf - 4 игрока")},
		"6 shot revolver":                        {message: SimpleMessage("Colt Python - 6 игроков")},
		"7 shot revolver":                        {message: SimpleMessage("Наган - 7 игроков")},
		"revolver has been replaced": {message: PluralMessage{
			one:  "Теперь в игре может участвовать %i игрок",
			few:  "Теперь в игре могут участвовать %i игрока",
			many: "Теперь в игре могут участвовать %i игроков",
		}},
		"settings will be applied for next games": {message: SimpleMessage("Изменения вступят в силу со следующего раунда")},
	},
}
