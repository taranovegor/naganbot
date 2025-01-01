package translator

var GameTranslations = translations{
	"ru": {
		"game creation": {
			oneOf: []oneOf{
				{message: SimpleMessage("Кто готов испытать свою судьбу в блистательной русской рулетке?! Ня-хаха~!")},
				{message: SimpleMessage("Вызываю всех отважных и кавайных героев! Начнём? Угу~!")},
				{message: SimpleMessage("Ооо, так вы снова хотите пощекотать нервы? Хорошо-хорошо, я с вами, десу~!")},
			},
		},
		"joining the game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Уааа! Ты пришёл! Присоединяйся, и покажи всем, какой ты герой-чан~!")},
				{message: SimpleMessage("Ты что, бессмертный? Или просто кавайный? В любом случае, добро пожаловать~!")},
				{message: SimpleMessage("Кири-кири, усаживайся, маэстро судьбы скоро начнёт игру!")},
				{message: SimpleMessage("Хех, ты точно готов рискнуть всем ради эпичного момента? Ня~?")},
				{message: SimpleMessage("Один из шести, шанс на победу! Веришь в себя? Ох, как круто!")},
				{message: SimpleMessage("Это место — для настоящих героев! Ну что, будет сугой?")},
				{message: SimpleMessage("А ты раньше играл в русскую рулетку? Сейчас всё будет по-эпическому!")},
			},
		},
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: SimpleMessage("Окей, судари и сударушки, начинаем волшебную битву судьбы! Ура-ра~!")},
						{message: SimpleMessage("Барабан заряжен, патрон вставлен... щелчок... *ба-бах*!")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Все готовы? Тогда волшебная игра началась! Кавайный барабан вращается!")},
						{message: SimpleMessage("Пистолет передан следующему игроку, пусть удача улыбнётся тебе, ня~!")},
					},
				},
			},
		},
		"gunslinger killed":     {message: SimpleMessage("%gunslinger встретил свою судьбу... Бака! Это была их последняя пуля, ояшии~!")},
		"joined the game":       {message: SimpleMessage("Игроки собрались в этот судьбоносный день: %date. Кто выживет?")},
		"game join list item":   {message: SimpleMessage("%num. %gunslinger — настоящий воин света!")},
		"owner of the game":     {message: SimpleMessage("заслуженный сенпай этой рулетки")},
		"shot in game":          {message: SimpleMessage("Ахх, %gunslinger выбывает из игры! О, это было драматично~!")},
		"top is not determined": {message: SimpleMessage("Ещё рано для списка кавайных чемпионов! Потерпите, ня~!")},
		"top players by games":  {message: SimpleMessage("Топ-%number самых блистательных героев русской рулетки~:")},
		"top game player": {message: PluralMessage{
			one:  "%i. %user — испытал судьбу %times раз, каваиии~!",
			few:  "%i. %user — испытал судьбу %times раза, каваиии~!",
			many: "%i. %user — испытал судьбу %times раз, каваиии~!",
		}},
		"available only in chat": {message: SimpleMessage("Эта игра не для одиночек, ня~! Приводи друзей, угу~!")},
		"player already in game": {message: SimpleMessage("Ты уже в игре, бака! Подожди своей очереди, ня~!")},
		"player is not kicked":   {message: SimpleMessage("Этот игрок всё ещё дышит! Может, у него есть шанс, ня~?")},
		"wait for game timeout":  {message: SimpleMessage("Игра ещё не началась! Терпение, мой юный герой~!")},
		"active game not found":  {message: SimpleMessage("Эх, никто ещё не начал игру! Кто будет первым?! Сугоооой~!")},
		"user game statistics": {message: PluralMessage{
			one:  "Ты участвовал в %games игр и проиграл %shots раз. Но ты всё равно кавайный~!",
			few:  "Ты участвовал в %games игр и проиграл %shots раза. Но ты всё равно кавайный~!",
			many: "Ты участвовал в %games игр и проиграл %shots раз. Но ты всё равно кавайный~!",
		}},
	},
}
